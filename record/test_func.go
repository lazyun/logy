package record

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
}

func fff(a ...interface{}) {
	fmt.Println(a...)
}

func AAA(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "A")
	//SetLocalTrigger(ctx, Error, Debug)
	defer CatchInfo(ctx)
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
