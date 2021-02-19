package service

import (
	"context"
	grpc_mixture "github.com/davveo/grpc-mixture"
)

const service_name_call_get = "callGet"

type CacheClient struct {

}

func (cc *CacheClient) CallGet(ctx context.Context, key string, csc grpc_mixture.CacheServiceClient) ( []byte, error) {
	getReq:=&grpc_mixture.GetReq{Key:key}
	getResp, err :=csc.Get(ctx, getReq )
	if err != nil {
		return nil, err
	}
	value := getResp.Value
	return value, err
}

func (cc *CacheClient) CallStore(key string, value []byte, client grpc_mixture.CacheServiceClient) ( *grpc_mixture.StoreResp, error) {
	storeReq := grpc_mixture.StoreReq{Key: key, Value: value}
	storeResp, err := client.Store(context.Background(), &storeReq)
	if err != nil {
		return nil, err
	}
	return storeResp, err
}

//func (cc *CacheClient) CallGet(ctx context.Context, key string, csc pb.CacheServiceClient) ( []byte, error) {
//	span := opentracing.StartSpan(service_name_call_get)
//	sc := span.Context()
//	a :=fmt.Sprintf("context:%+v", span.Context())
//	log.Println("a:", a)
//	scZipkin :=sc.(openzipkin.SpanContext)
//	log.Printf("zipkinTraceId:%v", scZipkin.TraceID)
//	defer span.Finish()
//	time.Sleep(5*time.Millisecond)
//	// Put root span in context so it will be used in our calls to the client.
//	ctx = opentracing.ContextWithSpan(ctx, span)
//	//ctx := context.Background()
//	getReq :=&pb.GetReq{Key:key}
//	getResp, err :=csc.Get(ctx, getReq )
//	if err != nil {
//		return nil, err
//	}
//	value := getResp.Value
//	return value, err
//}

