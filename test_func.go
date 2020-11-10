package logy

import (
	"context"
	"fmt"
)

func init() {
	SetGlobalTrigger(Error, Warning)
	RegisterDebug(fff)
	RegisterInfo(fff)
	RegisterWarn(fff)
	RegisterErr(fff)
	RegisterFatal(fff)
	RegisterUnified(fff)

	RegisterDebugF(fffF)
	RegisterInfoF(fffF)
	RegisterWarnF(fffF)
	RegisterErrF(fffF)
	RegisterFatalF(fffF)
	RegisterUnifiedF(fffF)
}

func fff(a ...interface{}) {
	fmt.Println(a...)
}

func fffF(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func AAA(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "A")
	SetLocalTrigger(ctx, Error, Debug)
	defer CatchInfo(ctx)
	Log(ctx, Debug)
	Log(ctx, Info, "Begin A. %s", "la~la~la~")
	B(ctx)
	Log(ctx, Warning, "End A. %s", "biu~biu~biu~")
}

func B(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "B")
	Log(ctx, Debug, "Begin B. %s", "la~la~la~")
	C(ctx)
	Log(ctx, Warning, "End B. %s", "biu~biu~biu~")
}

func C(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "C")
	Log(ctx, Error, "Begin C. %s", "la~la~la~")

	Log(ctx, Warning, "End C. %s", "biu~biu~biu~")
}

func AAAF(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "AF")
	SetLocalTrigger(ctx, Error, Debug)
	defer CatchInfo(ctx)
	Logf(ctx, Debug, "")
	Logf(ctx, Info, "Begin AF. %s", "la~la~la~")
	BF(ctx)
	Logf(ctx, Warning, "End AF. %s", "biu~biu~biu~")
}

func BF(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "BF")
	Logf(ctx, Debug, "Begin BF. %s", "la~la~la~")
	CF(ctx)
	Logf(ctx, Warning, "End BF. %s", "biu~biu~biu~")
}

func CF(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "CF")
	Logf(ctx, Error, "Begin CF. %s", "la~la~la~")

	Logf(ctx, Warning, "End CF. %s", "biu~biu~biu~")
}
