package middleware

import (
	"context"
	"errors"
	grpc_mixture "github.com/davveo/grpc-mixture"
	"log"
	"sync"
)

const (
	service_throttle = 5
)
var tm throttleMutex

type ThrottleMiddleware struct {
	Next grpc_mixture.CacheServiceServer
}

type throttleMutex struct {
	mu       sync.RWMutex
	throttle int
}

func (t *throttleMutex) getThrottle() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.throttle
}

func (t *throttleMutex) changeThrottle(delta int)  {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.throttle += delta
}

func (tg *ThrottleMiddleware) Get(ctx context.Context, req *grpc_mixture.GetReq) (*grpc_mixture.GetResp, error) {
	if tm.getThrottle() >= service_throttle {
		log.Printf("Get throttle=%v reached\n", tm.getThrottle())
		return nil, errors.New("service throttle reached, please try later")
	} else {
		tm.changeThrottle(1)
		resp, err := tg.Next.Get(ctx, req)
		tm.changeThrottle(-1)
		return resp, err
	}
}

func (tg *ThrottleMiddleware) Store(ctx context.Context, req *grpc_mixture.StoreReq) (*grpc_mixture.StoreResp, error) {
	return tg.Next.Store(ctx, req)
}

func (tg *ThrottleMiddleware) Dump(dr *grpc_mixture.DumpReq, csds grpc_mixture.CacheService_DumpServer) error {
	return tg.Next.Dump(dr,csds )

}