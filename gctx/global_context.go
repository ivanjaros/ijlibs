package gctx

import (
	"context"
	"sync"
)

type features struct {
	wg     sync.WaitGroup
	cancel context.CancelFunc
	once   sync.Once
}

type key struct{}

func New() context.Context {
	var f features
	ctx, cancel := context.WithCancel(context.Background())
	f.cancel = cancel
	ctx = context.WithValue(ctx, key{}, &f)
	return ctx
}

func Cancel(ctx context.Context) {
	f := ctx.Value(key{}).(*features)
	f.cancel()
}

func Add(ctx context.Context, num int) {
	f := ctx.Value(key{}).(*features)
	f.wg.Add(num)
}

func Done(ctx context.Context) {
	f := ctx.Value(key{}).(*features)
	f.wg.Done()
}

func Wait(ctx context.Context) {
	f := ctx.Value(key{}).(*features)
	f.wg.Wait()
}

func Once(ctx context.Context, fn func()) {
	f := ctx.Value(key{}).(*features)
	f.once.Do(fn)
}
