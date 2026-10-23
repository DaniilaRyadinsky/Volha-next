package requestid

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestContextHTTPAndGRPCPropagation(t *testing.T) {
	require.NotEmpty(t, New())
	require.Empty(t, FromContext(context.Background()))
	ctx := WithContext(context.Background(), "request-123")
	require.Equal(t, "request-123", FromContext(ctx))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(Header, "supplied")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	require.NoError(t, EchoMiddleware()(func(c echo.Context) error {
		require.Equal(t, "supplied", FromContext(c.Request().Context()))
		return c.NoContent(204)
	})(c))
	require.Equal(t, "supplied", rec.Header().Get(Header))
	err := UnaryClientInterceptor()(ctx, "/service/method", nil, nil, nil, func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		require.True(t, ok)
		require.Equal(t, []string{"request-123"}, md.Get(MetadataKey))
		return nil
	})
	require.NoError(t, err)
	incoming := metadata.NewIncomingContext(context.Background(), metadata.Pairs(MetadataKey, "incoming"))
	_, err = UnaryServerInterceptor()(incoming, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, _ any) (any, error) {
		require.Equal(t, "incoming", FromContext(ctx))
		return nil, nil
	})
	require.NoError(t, err)
}
