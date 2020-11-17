package logy

import (
	"context"
	"fmt"
	"runtime"

	"github.com/google/uuid"
)

type levelMapping struct {
	Debug   nullFunc
	Info    nullFunc
	Warning nullFunc
	Error   nullFunc
	Fatal   nullFunc
	Unified nullFunc
}

type levelMappingFormat struct {
	Debug   nullFuncF
	Info    nullFuncF
	Warning nullFuncF
	Error   nullFuncF
	Fatal   nullFuncF
	Unified nullFuncF
}

type traceRoot struct {
	UUID     string
	Title    string
	MaxLevel logLevel

	DoList    [][]interface{}
	callStack []*traceInfo

	OwnLevel    bool
	OwnOccurLev logLevel
	OwnOutLev   logLevel
}

type traceInfo struct {
	Title    string
	FuncName string
	FilePath string
	Line     int
	DoList   [][]interface{}
}

type nullFunc func(...interface{})
type nullFuncF func(string, ...interface{})

type FormatLog func(fields LogFields, args ...interface{}) []interface{}
type FormatLogF func(fields LogFields, format string, args ...interface{}) (string, []interface{})

type ctxKey string
type logLevel int

type LogFields struct {
	UUID     string
	Level    string
	Title    string
	FuncName string
	FilePath string
	Line     int
}

const (
	offset = 2

	UUIDKeyName   ctxKey = "uuid"
	ctxKeyName    ctxKey = "Love"
	ctxSubKeyName ctxKey = "LoveMe"

	Debug   logLevel = 0
	Info    logLevel = 1
	Warning logLevel = 2
	Error   logLevel = 3
	Fatal   logLevel = 4

	logTypeSprint  int = 1
	logTypeSprintF int = 2
)

var (
	lm = levelMapping{
		Debug:   func(i ...interface{}) {},
		Info:    func(i ...interface{}) {},
		Warning: func(i ...interface{}) {},
		Error:   func(i ...interface{}) {},
		Fatal:   func(i ...interface{}) {},
		Unified: func(i ...interface{}) {},
	}

	lmf = levelMappingFormat{
		Debug:   func(f string, i ...interface{}) {},
		Info:    func(f string, i ...interface{}) {},
		Warning: func(f string, i ...interface{}) {},
		Error:   func(f string, i ...interface{}) {},
		Fatal:   func(f string, i ...interface{}) {},
		Unified: func(f string, i ...interface{}) {},
	}

	levelString = map[logLevel]string{
		0: "Debug",
		1: "Info",
		2: "Warning",
		3: "Error",
		4: "Fatal",
	}

	occurLev = Error
	outLev   = Info

	formatFunc FormatLog = func(fields LogFields, args ...interface{}) []interface{} {
		if 0 != len(args) {
			args[0] = fmt.Sprintf("%s\t%s:%d\t[%s] %v", fields.UUID, fields.FilePath, fields.Line, fields.Title, args[0])
			return append([]interface{}{fields.Level}, args...)
		} else {
			return append([]interface{}{fields.Level}, fmt.Sprintf("%s\t%s:%d\t[%s]", fields.UUID, fields.FilePath, fields.Line, fields.Title))
		}
	}

	formatFuncF FormatLogF = func(fields LogFields, format string, args ...interface{}) (string, []interface{}) {
		return fmt.Sprintf("%v\t%s\t%s:%d\t[%s] %v", fields.Level, fields.UUID, fields.FilePath, fields.Line, fields.Title, format), args
	}
)

func RegisterLogFormat(f FormatLog) {
	formatFunc = f
}

func RegisterLogFormatF(f FormatLogF) {
	formatFuncF = f
}

func RegisterDebug(debug nullFunc) {
	lm.Debug = debug
}

func RegisterDebugF(debug nullFuncF) {
	lmf.Debug = debug
}

func RegisterWarn(warn nullFunc) {
	lm.Warning = warn
}

func RegisterWarnF(warn nullFuncF) {
	lmf.Warning = warn
}

func RegisterInfo(info nullFunc) {
	lm.Info = info
}

func RegisterInfoF(info nullFuncF) {
	lmf.Info = info
}

func RegisterErr(err nullFunc) {
	lm.Error = err
}

func RegisterErrF(err nullFuncF) {
	lmf.Error = err
}

func RegisterFatal(fatal nullFunc) {
	lm.Fatal = fatal
}

func RegisterFatalF(fatal nullFuncF) {
	lmf.Fatal = fatal
}

func RegisterUnified(unified nullFunc) {
	lm.Unified = unified
}

func RegisterUnifiedF(unified nullFuncF) {
	lmf.Unified = unified
}

func SetGlobalTrigger(occurLevel, outLevel logLevel) {
	occurLev = occurLevel
	outLev = outLevel
}

func SetLocalTrigger(ctx context.Context, occurLevel, outLevel logLevel) {
	root := ctx.Value(ctxKeyName)
	if nil == root {
		return
	}

	rootCallStack := root.(*traceRoot)

	(*rootCallStack).OwnLevel = true
	(*rootCallStack).OwnOccurLev = occurLevel
	(*rootCallStack).OwnOutLev = outLevel

	//fmt.Println("outLevel", (*rootCallStack).OwnOccurLev, (*rootCallStack).OwnOutLev, occurLevel, outLevel)
}

func SetFuncSignal(ctx context.Context, s string, depth int) context.Context {
	var (
		_, file, line, _ = runtime.Caller(depth)
		nowCallStack     = traceInfo{FilePath: file, Line: line, FuncName: s}
	)

	root := ctx.Value(ctxKeyName)
	if nil == root {
		nowCallStack.Title = s

		uid := ctx.Value(UUIDKeyName)
		if nil == uid {
			uid = uuid.New().String()
		}

		r := traceRoot{
			fmt.Sprint(uid),
			s,
			Debug,
			[][]interface{}{},
			[]*traceInfo{},
			false,
			Debug,
			Debug,
		}

		ctx = context.WithValue(ctx, ctxKeyName, &r)
		return context.WithValue(ctx, ctxSubKeyName, nowCallStack)
	}

	rootCallStack := root.(*traceRoot)

	sub := ctx.Value(ctxSubKeyName)
	if nil == sub {
		nowCallStack.Title = (*rootCallStack).Title + "-" + s
	} else {
		subCallStack := sub.(traceInfo)
		nowCallStack.Title = subCallStack.Title + "-" + s
	}

	(*rootCallStack).callStack = append((*rootCallStack).callStack, &nowCallStack)
	return context.WithValue(ctx, ctxSubKeyName, nowCallStack)
}

func Log(ctx context.Context, level logLevel, args ...interface{}) {
	root := ctx.Value(ctxKeyName)
	if nil == root {
		lmLog(level, formatFunc(LogFields{UUID: uuid.New().String(), Level: levelString[level], Title: "-"}, args...)...)
		return
	}

	sub := ctx.Value(ctxSubKeyName)
	if nil == sub {
		return
	}

	rootCallStack := root.(*traceRoot)
	subCallStack := sub.(traceInfo)

	if level > (*rootCallStack).MaxLevel {
		(*rootCallStack).MaxLevel = level
	}

	var (
		value     = []interface{}{level, logTypeSprint}
		logFields = LogFields{
			UUID:     uuid.New().String(),
			Level:    levelString[level],
			Title:    subCallStack.Title,
			FuncName: subCallStack.FuncName,
			FilePath: subCallStack.FilePath,
			Line:     subCallStack.Line,
		}
	)
	//if 0 != len(args) {
	//	args[0] = fmt.Sprintf("[%s] %v", subCallStack.Title, args[0])
	//	value = append(value, args...)
	//} else {
	//	value = append(value, fmt.Sprintf("[%s]", subCallStack.Title))
	//}

	value = append(value, formatFunc(logFields, args...)...)

	(*rootCallStack).DoList = append((*rootCallStack).DoList, value)
}

func Logf(ctx context.Context, level logLevel, format string, args ...interface{}) {
	root := ctx.Value(ctxKeyName)
	if nil == root {
		f, a := formatFuncF(LogFields{UUID: uuid.New().String(), Level: levelString[level], Title: "-"}, format, args...)
		lmfLog(level, f, a...)
		return
	}

	sub := ctx.Value(ctxSubKeyName)
	if nil == sub {
		return
	}

	rootCallStack := root.(*traceRoot)
	subCallStack := sub.(traceInfo)

	if level > (*rootCallStack).MaxLevel {
		(*rootCallStack).MaxLevel = level
	}

	var (
		logFields = LogFields{
			UUID:     uuid.New().String(),
			Level:    levelString[level],
			Title:    subCallStack.Title,
			FuncName: subCallStack.FuncName,
			FilePath: subCallStack.FilePath,
			Line:     subCallStack.Line,
		}

		f, a  = formatFuncF(logFields, format, args...)
		value = []interface{}{level, logTypeSprintF, f}
	)
	//if 0 != len(args) {
	//	args[0] = fmt.Sprintf("[%s] %v", subCallStack.Title, args[0])
	//	value = append(value, args...)
	//} else {
	//	value = append(value, fmt.Sprintf("[%s]", subCallStack.Title))
	//}

	value = append(value, a...)

	(*rootCallStack).DoList = append((*rootCallStack).DoList, value)
}

func CatchInfo(ctx context.Context) {
	root := ctx.Value(ctxKeyName)
	if nil == root {
		return
	}

	rootCallStack := root.(*traceRoot)

	var (
		nowOccurLev logLevel
		nowOutLev   logLevel
	)

	if (*rootCallStack).OwnLevel {
		nowOccurLev = (*rootCallStack).OwnOccurLev
		nowOutLev = (*rootCallStack).OwnOutLev
		//fmt.Println((*rootCallStack).OwnOccurLev, (*rootCallStack).OwnOutLev)
	} else {
		nowOccurLev = occurLev
		nowOutLev = outLev
	}

	//if (*rootCallStack).MaxLevel < nowOccurLev {
	//	return
	//}

	for _, value := range (*rootCallStack).DoList {
		level := value[0].(logLevel)
		logFormatType := value[1].(int)

		d := value
		//d[0] = levelString[level]

		if logTypeSprint == logFormatType {
			if (*rootCallStack).MaxLevel >= nowOccurLev && level >= nowOutLev {
				lm.Unified(d[2:]...)
				continue
			}

			lmLog(level, d[2:]...)
			continue
		}

		format := fmt.Sprint(d[2])
		if (*rootCallStack).MaxLevel >= nowOccurLev && level >= nowOutLev {
			lmf.Unified(format, d[3:]...)
			continue
		}

		lmfLog(level, format, d[3:]...)
		//fmt.Println(nowOutLev, value)
	}
}

func lmLog(level logLevel, args ...interface{}) {
	switch level {
	case Debug:
		{
			lm.Debug(args...)
		}
	case Info:
		{
			lm.Info(args...)
		}
	case Warning:
		{
			lm.Warning(args...)
		}
	case Error:
		{
			lm.Error(args...)
		}
	case Fatal:
		{
			lm.Fatal(args...)
		}
	}
}

func lmfLog(level logLevel, format string, args ...interface{}) {
	switch level {
	case Debug:
		{
			lmf.Debug(format, args...)
		}
	case Info:
		{
			lmf.Info(format, args...)
		}
	case Warning:
		{
			lmf.Warning(format, args...)
		}
	case Error:
		{
			lmf.Error(format, args...)
		}
	case Fatal:
		{
			lmf.Fatal(format, args...)
		}
	}
}
