package logy

import (
	"context"
	"fmt"
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
	Title    string
	MaxLevel logLevel

	DoList    [][]interface{}
	callStack []*traceInfo

	OwnLevel    bool
	OwnOccurLev logLevel
	OwnOutLev   logLevel
}

type traceInfo struct {
	Title  string
	DoList [][]interface{}
}

type nullFunc func(...interface{})
type nullFuncF func(string, ...interface{})

type FormatLog func(fields LogFields, args ...interface{}) []interface{}
type FormatLogF func(fields LogFields, format string, args ...interface{}) (string, []interface{})

type ctxKey string
type logLevel int
type LogFields []interface{}

const (
	offset = 2

	ctxKeyName    ctxKey = "Love"
	ctxSubKeyName ctxKey = "LoveMe"

	Debug   logLevel = 0
	Info    logLevel = 1
	Warning logLevel = 2
	Error   logLevel = 3
	Fatal   logLevel = 4

	logTypeSprint int = iota
	logTypeSprintF
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
			args[0] = fmt.Sprintf("[%s] %v", fields[1], args[0])
		}

		if 0 != len(args) {
			args[0] = fmt.Sprintf("[%s] %v", fields[1], args[0])
			return append([]interface{}{fields[0]}, args...)
		} else {
			return append([]interface{}{fields[0]}, fmt.Sprintf("[%s]", fields[1]))
		}
	}

	formatFuncF FormatLogF = func(fields LogFields, format string, args ...interface{}) (string, []interface{}) {
		return fmt.Sprintf("%v [%s] %v", fields[0], fields[1], format), args
	}

	_ LogFields = []interface{}{Info, "title"}
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

func SetFuncSignal(ctx context.Context, s string) context.Context {
	var nowCallStack = traceInfo{}

	root := ctx.Value(ctxKeyName)
	if nil == root {
		nowCallStack.Title = s

		r := traceRoot{
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

	var value = []interface{}{level, logTypeSprint}
	//if 0 != len(args) {
	//	args[0] = fmt.Sprintf("[%s] %v", subCallStack.Title, args[0])
	//	value = append(value, args...)
	//} else {
	//	value = append(value, fmt.Sprintf("[%s]", subCallStack.Title))
	//}

	value = append(value, formatFunc(LogFields{levelString[level], subCallStack.Title}, args...)...)

	(*rootCallStack).DoList = append((*rootCallStack).DoList, value)
}

func Logf(ctx context.Context, level logLevel, format string, args ...interface{}) {
	root := ctx.Value(ctxKeyName)
	if nil == root {
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
		f, a  = formatFuncF(LogFields{levelString[level], subCallStack.Title}, format, args...)
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

	if (*rootCallStack).MaxLevel < nowOccurLev {
		return
	}

	for _, value := range (*rootCallStack).DoList {
		level := value[0].(logLevel)
		logFormatType := value[1].(int)

		d := value
		//d[0] = levelString[level]

		if logTypeSprint == logFormatType {
			if level >= nowOutLev {
				lm.Unified(d[2:]...)
				continue
			}

			switch level {
			case Debug:
				{
					lm.Debug(d[2:]...)
				}
			case Info:
				{
					lm.Info(d[2:]...)
				}
			case Warning:
				{
					lm.Warning(d[2:]...)
				}
			case Error:
				{
					lm.Error(d[2:]...)
				}
			case Fatal:
				{
					lm.Fatal(d[2:]...)
				}
			}

			return
		}

		format := fmt.Sprint(d[2])
		if level >= nowOutLev {
			lmf.Unified(format, d[3:]...)
			continue
		}

		switch level {
		case Debug:
			{
				lmf.Debug(format, d[3:]...)
			}
		case Info:
			{
				lmf.Info(format, d[3:]...)
			}
		case Warning:
			{
				lmf.Warning(format, d[3:]...)
			}
		case Error:
			{
				lmf.Error(format, d[3:]...)
			}
		case Fatal:
			{
				lmf.Fatal(format, d[3:]...)
			}
		}

		//fmt.Println(nowOutLev, value)
	}
}
