package logy

import (
	"context"
	//"fmt"
)

type rootLogy struct {
	UUID string
	Title string
	WithValue map[string]interface{}

	MaxLevel string

	OwnLevel    bool
	OwnOccurLev logLevel
	OwnOutLev   logLevel

	SubInfo map[string]rootLogy
}

// uuid: obj
var callInfo = map[string]rootLogy{}

func Log1(ctx context.Context, level logLevel, args ...interface{}) {
	//root := ctx.Value(ctxKeyName)
	//if nil == root {
	//
	//}
	//
	//dst, ok := callInfo[fmt.Sprint(root)]
	//if !ok {
	//
	//}

	//dst
}