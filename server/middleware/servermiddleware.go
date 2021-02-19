package middleware

import (
	"context"
	grpc_mixture "github.com/davveo/grpc-mixture"
)

type CacheServiceMiddleware struct {
	Next grpc_mixture.CacheServiceServer
}

func BuildGetMiddleware(s grpc_mixture.CacheServiceServer) grpc_mixture.CacheServiceServer {
	tm := ThrottleMiddleware{s}
	cs := CacheServiceMiddleware{&tm}
	return &cs
}

func (csm *CacheServiceMiddleware) Get(ctx context.Context, req *grpc_mixture.GetReq) (*grpc_mixture.GetResp, error) {
	return csm.Next.Get(ctx, req)
}

func (csm *CacheServiceMiddleware) Store(ctx context.Context, req *grpc_mixture.StoreReq) (*grpc_mixture.StoreResp, error) {
	return csm.Next.Store(ctx, req)
}

func (csm *CacheServiceMiddleware) Dump(dr *grpc_mixture.DumpReq, csds grpc_mixture.CacheService_DumpServer) error {
	return csm.Next.Dump(dr, csds)
}