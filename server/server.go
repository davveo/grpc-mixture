package main

import (
	grpc_mixture "github.com/davveo/grpc-mixture"
	"github.com/davveo/grpc-mixture/server/middleware"
	"github.com/davveo/grpc-mixture/server/service"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	openZipkin "github.com/openzipkin/zipkin-go-opentracing"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	serviceNameCacheServer = "cache service"
	hostUrl                = "localhost:5051"
	endpointUrl            =  "http://localhost:9411/api/v1/spans"
)

func main() {
	var (
		err error
		conn net.Listener
		tracer opentracing.Tracer
	)

	log.Println("starting server...")
	if conn, err = net.Listen("tcp", hostUrl); err != nil {
		panic(err)
	}
	if tracer, err = newTracer(); err != nil {
		panic(err)
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads())),
	}
	server := grpc.NewServer(opts...)
	cs := initCache()

	grpc_mixture.RegisterCacheServiceServer(server, cs)

	if err = server.Serve(conn); err != nil {
		panic(err)
	}

	log.Println("server listening on port 5051")
}

func initCache() grpc_mixture.CacheServiceServer {
	s := make(map[string][]byte)
	s["123"] = []byte{123}
	s["231"] = []byte{231}
	cs := service.CacheService{Storage: s}
	return middleware.BuildGetMiddleware(&cs)
}

func newTracer() (opentracing.Tracer, error)  {
	var (
		err error
		collector openZipkin.Collector
		recorder openZipkin.SpanRecorder
		tracer opentracing.Tracer
	)
	if collector, err = openZipkin.NewHTTPCollector(endpointUrl); err != nil {
		return nil, err
	}

	recorder = openZipkin.NewRecorder(collector, true, hostUrl, serviceNameCacheServer)
	if tracer, err = openZipkin.NewTracer(recorder, openZipkin.ClientServerSameSpan(true)); err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, nil
}