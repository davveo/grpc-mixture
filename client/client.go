package main

import (
	"context"
	"fmt"
	grpc_mixture "github.com/davveo/grpc-mixture"
	"github.com/davveo/grpc-mixture/client/middleware"
	"github.com/davveo/grpc-mixture/client/service"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	openzipkin "github.com/openzipkin/zipkin-go-opentracing"
	zipkintracer "github.com/openzipkin/zipkin-go-opentracing"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	endpointUrl            = "http://localhost:9411/api/v1/spans"
	hostUrl                = "localhost:5051"
	serviceNameCacheClient = "cache service client"
)

func main() {
	var (
		err error
		conn *grpc.ClientConn
		tracer opentracing.Tracer
		collector zipkintracer.Collector
	)

	log.Println("starting server...")
	if tracer, collector, err = newTracer(); err != nil {
		panic(err)
	}
	defer collector.Close()
	if conn, err = grpc.Dial(hostUrl, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads()))); err != nil {
		panic(err)
	}
	defer conn.Close()
	client := grpc_mixture.NewCacheServiceClient(conn)

	callServer(client)
	testCircuitBreaker(client)
}

func callServer(csc grpc_mixture.CacheServiceClient)  {
	ctx :=context.Background()
	key:="123"
	cc := service.CacheClient{}
	cg := middleware.BuildGetMiddleware(&cc)
	value, err := cg.CallGet(ctx, key, csc)
	if err != nil {
		fmt.Println("error call get:", err)
	} else {
		fmt.Printf("value=%v for key=%v\n", value, key)
	}
	key = "231"
	value, err = cg.CallGet(ctx, key, csc)
	if err != nil {
		fmt.Println("error call get:", err)
	} else {
		fmt.Printf("value=%v for key=%v\n", value, key)
	}

}

func testCircuitBreaker(csc grpc_mixture.CacheServiceClient)  {
	callServer(csc)
	time.Sleep(time.Duration(20*1000)*time.Millisecond)

	callServer(csc)
	time.Sleep(time.Duration(5*1000)*time.Millisecond)
	callServer(csc)
	time.Sleep(time.Duration(20*2000)*time.Millisecond)

	callServer(csc)
}

func newTracer() (opentracing.Tracer, zipkintracer.Collector, error)  {
	var (
		err error
		collector zipkintracer.Collector
		tracer opentracing.Tracer
	)

	if collector, err = openzipkin.NewHTTPCollector(endpointUrl); err != nil {
		return nil, nil, err
	}
	recorder :=openzipkin.NewRecorder(collector,
		true, hostUrl, serviceNameCacheClient)
	if tracer, err = openzipkin.NewTracer(
		recorder,
		openzipkin.ClientServerSameSpan(true)); err !=nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer,collector, nil
}
