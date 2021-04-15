package logy

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
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
	ctx = SetFuncSignal(ctx, "A", 1)
	SetLocalTrigger(ctx, Error, Debug)
	RegisterLogFormat(func(fields LogFields, args ...interface{}) []interface{} {
		value := []interface{}{
			fmt.Sprintf("%s %v <%s><%s><%s>[%s]",
				time.Now().Format("2006-01-02 15:04:05.999"),
				fields.Level,
				strings.TrimPrefix(fields.FilePath, "/root/janus-api/src/janus/"),
				fields.FuncName,
				fields.UUID,
				fields.Title),
		}

		if 0 == len(args) {
			return value
		} else {
			return append(value, args...)
		}
	})

	RegisterLogFormatF(func(fields LogFields, format string, args ...interface{}) (string, []interface{}) {
		return fmt.Sprintf("%s %v <%s><%s><%s>[%s]%v",
			time.Now().Format("2006-01-02 15:04:05.999"),
			fields.Level,
			strings.TrimPrefix(fields.FilePath, "/root/janus-api/src/janus/"),
			fields.FuncName,
			fields.UUID,
			fields.Title,
			format), args
	})

	defer CatchInfo(ctx)
	Log(ctx, Debug)
	Log(ctx, Debug, "123")
	Log(ctx, Info, "Begin A. %s", "la~la~la~")
	B(ctx)
	Log(ctx, Warning, "End A. %s", "biu~biu~biu~")
	fmt.Println("uuid", GetUUID(ctx))
}

func B(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "B", 1)
	Log(ctx, Debug, "Begin B. %s", "la~la~la~")
	C(ctx)
	Log(ctx, Warning, "End B. %s", "biu~biu~biu~")
}

func C(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "C", 1)
	Log(ctx, Error, "Begin C. %s", "la~la~la~")

	Log(ctx, Warning, "End C. %s", "biu~biu~biu~")
}

func AAAF(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "AF", 1)
	SetLocalTrigger(ctx, Error, Debug)
	defer CatchInfo(ctx)
	Logf(ctx, Debug, "")
	Logf(ctx, Info, "Begin AF. %s", "la~la~la~")
	BF(ctx)
	Logf(ctx, Warning, "End AF. %s", "biu~biu~biu~")
}

func BF(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "BF", 1)
	Logf(ctx, Debug, "Begin BF. %s", "la~la~la~")
	CF(ctx)
	Logf(ctx, Warning, "End BF. %s", "biu~biu~biu~")
}

func CF(ctx context.Context) {
	ctx = SetFuncSignal(ctx, "CF", 1)
	Logf(ctx, Error, "Begin CF. %s", "la~la~la~")

	Logf(ctx, Warning, "End CF. %s", "biu~biu~biu~")
}

func LogImmediately() {
	Log(context.Background(), Info, "la~la~la~", "biu~biu~biu~")
}

func LogImmediatelyTitle() {
	ctx := SetFuncSignal(context.Background(), "LogImmediatelyTitle", 1)
	Log(ctx, Info, "la~la~la~", "biu~biu~biu~")
}

func TestAAA(t *testing.T) {
	ctx := context.Background()
	AAA(ctx)

	//RegisterUnified(nil)
}

func TestLogf(t *testing.T) {
	ctx := context.Background()
	AAAF(ctx)
}

func TestLogImmediately(t *testing.T) {
	LogImmediately()
}

func TestLogImmediatelyTitle(t *testing.T) {
	LogImmediatelyTitle()
}

func Test123(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4}
	)

	t.Log(a[1:])
}