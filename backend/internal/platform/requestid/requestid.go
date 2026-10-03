package requestid

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const Header = "X-Request-ID"
const MetadataKey = "x-request-id"

type contextKey struct{}

func New() string {
	return xid.New().String()
}

func WithContext(ctx context.Context, id string) context.Context {
	if id == "" {
		id = New()
	}
	return context.WithValue(ctx, contextKey{}, id)
}

func FromContext(ctx context.Context) string {
	if id, ok := ctx.Value(contextKey{}).(string); ok {
		return id
	}
	return ""
}

func EchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get(Header)
			if id == "" {
				id = New()
			}

			c.Response().Header().Set(Header, id)
			req := c.Request().WithContext(WithContext(c.Request().Context(), id))
			c.SetRequest(req)
			return next(c)
		}
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if id := FromContext(ctx); id != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, MetadataKey, id)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get(MetadataKey); len(values) > 0 {
				ctx = WithContext(ctx, values[0])
			}
		}
		return handler(ctx, req)
	}
}
