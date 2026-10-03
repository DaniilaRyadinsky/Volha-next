package handlers

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapGRPCErrorSemantics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code codes.Code
		want int
	}{
		{name: "invalid argument", code: codes.InvalidArgument, want: http.StatusBadRequest},
		{name: "not found", code: codes.NotFound, want: http.StatusNotFound},
		{name: "already exists", code: codes.AlreadyExists, want: http.StatusConflict},
		{name: "unavailable", code: codes.Unavailable, want: http.StatusServiceUnavailable},
		{name: "deadline exceeded", code: codes.DeadlineExceeded, want: http.StatusGatewayTimeout},
		{name: "internal", code: codes.Internal, want: http.StatusBadGateway},
	}

	api := &Apis{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, _ := api.mapGRPCError(context.Background(), "test", status.Error(tt.code, "message"))
			if got != tt.want {
				t.Fatalf("mapGRPCError(%s) = %d, want %d", tt.code, got, tt.want)
			}
		})
	}
}
