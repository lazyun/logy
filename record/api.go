package record

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

type nullFunc func(i ...interface{})

type ctxKey string
type logLevel int

const (
	offset = 2

	ctxKeyName    ctxKey = "Love"
	ctxSubKeyName ctxKey = "LoveMe"

	Debug   logLevel = 0
	Info    logLevel = 1
	Warning logLevel = 2
	Error   logLevel = 3
	Fatal   logLevel = 4
)

var (
	lm = levelMapping{
		Debug:   func(i ...interface{}) {},
		Info:    func(i ...interface{}) {},
		Warning: func(i ...interface{}) {},
		Error:   func(i ...interface{}) {},
		Fatal:   func(i ...interface{}) {},
	}

	occurLev = Error
	outLev   = Info
)
func RegisterDebug(debug func(...interface{})) {
	lm.Debug = debug
}

func RegisterWarn(warn func(...interface{})) {
	lm.Warning = warn
}

func RegisterInfo(info func(...interface{})) {
	lm.Info = info
}

func RegisterErr(err func(...interface{})) {
	lm.Error = err
}

func RegisterFatal(fatal func(...interface{})) {
	lm.Fatal = fatal
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

	var value []interface{}
	value = append(value, level)

	if 0 != len(args) {
		args[0] = fmt.Sprintf("[%s] %v", subCallStack.Title, args[0])
	}

	value = append(value, args...)

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

		if level < nowOutLev {
			continue
		}

		d := value

		switch level {
		case Debug:
			{
				d[0] = "Debug"
				lm.Debug(d...)
			}
		case Info:
			{
				d[0] = "Info"
				lm.Info(d...)
			}
		case Warning:
			{
				d[0] = "Warning"
				lm.Warning(d...)
			}
		case Error:
			{
				d[0] = "Error"
				lm.Error(d...)
			}
		case Fatal:
			{
				d[0] = "Fatal"
				lm.Fatal(d...)
			}
		}

		//fmt.Println(nowOutLev, value)
	}
}
