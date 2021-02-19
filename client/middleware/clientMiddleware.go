package middleware

import (
	grpc_mixture "github.com/davveo/grpc-mixture"
	"golang.org/x/net/context"
)

type callGetter interface {
	CallGet(ctx context.Context, key string, c grpc_mixture.CacheServiceClient) ( []byte, error)
}
type CallGetMiddleware struct {
	Next callGetter
}
func BuildGetMiddleware(cc callGetter) callGetter {
	cbcg := CircuitBreakerCallGet{cc}
	tcg := TimeoutCallGet{&cbcg}
	rcg := RetryCallGet{&tcg}
	return &rcg
}

func (cg *CallGetMiddleware) CallGet(ctx context.Context, key string, csc grpc_mixture.CacheServiceClient) ( []byte, error) {
	return cg.Next.CallGet(ctx, key, csc)
}

//func BuildGetMiddleware() callGetter {
//	cc := service.CacheClient{}
//	//tcg := TimeoutCallGet{&cc}
//	//rcg := RetryCallGet{&tcg}
//	return &cc
//}



