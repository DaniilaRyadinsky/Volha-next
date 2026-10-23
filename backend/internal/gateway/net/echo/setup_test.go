package echo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/autumnterror/volha-backend/internal/gateway/config"
	authgrpc "github.com/autumnterror/volha-backend/internal/gateway/grpc/auth"
	"github.com/autumnterror/volha-backend/internal/gateway/grpc/products"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/stretchr/testify/require"
)

func TestRoutesHealthAdminPolicyAndPDFRange(t *testing.T) {
	t.Chdir(t.TempDir())
	require.NoError(t, os.Mkdir("pdfs", 0755))
	require.NoError(t, os.WriteFile("pdfs/guide.pdf", []byte("%PDF-1.7 example"), 0600))
	cfg := &config.Config{Mode: "PROD", Port: -1, CORSAllowedOrigins: []string{"http://frontend"}}
	app := New(nil, &products.Client{}, &authgrpc.Client{}, cfg, observability.Config{Service: "gateway"}, observability.NewMetrics("gateway", "dev", "test"))
	rec := httptest.NewRecorder()
	app.e.ServeHTTP(rec, httptest.NewRequest("GET", "/health/live", nil))
	require.Equal(t, 200, rec.Code)
	require.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
	rec = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/article", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	app.e.ServeHTTP(rec, req)
	require.Equal(t, 403, rec.Code)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/pdfs/guide.pdf", nil)
	req.Header.Set("Range", "bytes=0-4")
	app.e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusPartialContent, rec.Code)
	require.Equal(t, "%PDF-", rec.Body.String())
	require.Contains(t, rec.Header().Get("Content-Type"), "application/pdf")
	routes := map[string]bool{}
	for _, route := range app.e.Routes() {
		routes[route.Method+" "+route.Path] = true
	}
	for _, route := range []string{"GET /api/article/all", "GET /api/article", "POST /api/article", "PUT /api/article", "DELETE /api/article"} {
		require.True(t, routes[route], route)
	}
	require.Error(t, app.Run())
	require.NoError(t, app.Stop(context.Background()))
}
