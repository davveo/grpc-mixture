package service

import (
	"context"
	"fmt"
	grpc_mixture "github.com/davveo/grpc-mixture"
	"github.com/opentracing/opentracing-go"
	"time"
)

const serviceNameDbQueryUser = "db query user"

type CacheService struct {
	Storage map[string][]byte
}

func (c CacheService) Store(ctx context.Context, req *grpc_mixture.StoreReq) (*grpc_mixture.StoreResp, error) {
	key := req.Key
	value := req.Value
	if oldValue, ok := c.Storage[key]; ok {
		c.Storage[key] = value
		fmt.Printf(" key=%v already exist, old vale=%v|replaced with new value=%v\n", key, oldValue, c.Storage)
	} else {
		c.Storage[key] = value
		fmt.Printf(" key=%v not existing, add new value=%v\n", key, c.Storage)
	}
	r := &grpc_mixture.StoreResp{}
	return r, nil
}

func (c CacheService) Get(ctx context.Context, req *grpc_mixture.GetReq) (*grpc_mixture.GetResp, error) {
	time.Sleep(5 * time.Millisecond)
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			mysqlSpan := tracer.StartSpan(serviceNameDbQueryUser, opentracing.ChildOf(pctx))
			defer mysqlSpan.Finish()
			// do some op
			time.Sleep(time.Millisecond * 10)
		}
	}
	k := req.GetKey()
	value := c.Storage[k]
	fmt.Println("get called with return of value: ", value)
	resp := &grpc_mixture.GetResp{Value: value}
	return resp, nil
}

func (c CacheService) Dump(req *grpc_mixture.DumpReq, server grpc_mixture.CacheService_DumpServer) error {
	return nil
}


