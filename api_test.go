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

func TestLogImmediately(t *testing.T) {
	LogImmediately()
}

func TestLogImmediatelyTitle(t *testing.T) {
	LogImmediatelyTitle()
}