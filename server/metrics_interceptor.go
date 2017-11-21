package server

import (
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)

// unaryMetricsInterceptor returns a new unary server interceptor for our own OpenTracing-based metrics.
func unaryMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		span := opentracing.StartSpan("api")
		span.SetTag("fn_operation", info.FullMethod)
		newContext := opentracing.ContextWithSpan(ctx, span)
		defer span.Finish()
		resp, err := handler(newContext, req)
		return resp, err
	}
}

// streamMetricsInterceptor returns a new streaming server interceptor for our own OpenTracing-based metrics.
func streamMetricsInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		span := opentracing.StartSpan("api")
		span.SetTag("fn_operation", info.FullMethod)
		newContext := opentracing.ContextWithSpan(stream.Context(), span)
		defer span.Finish()
		wrappedStream := grpc_middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = newContext
		err := handler(srv, wrappedStream)
		return err
	}
}
