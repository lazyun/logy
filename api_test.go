package logy

import (
	"context"
	"testing"
)

func TestAAA(t *testing.T) {
	ctx := context.Background()
	AAA(ctx)
}

func TestLogf(t *testing.T) {
	ctx := context.Background()
	AAAF(ctx)
}