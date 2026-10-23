//go:build integration

package testkit

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"testing"
	"time"

	"crypto/rand"

	"github.com/alicebob/miniredis/v2"
	pb "github.com/autumnterror/volha-backend/api/proto/gen"
	authapi "github.com/autumnterror/volha-backend/internal/auth-service/api"
	authrepo "github.com/autumnterror/volha-backend/internal/auth-service/repository"
	authservice "github.com/autumnterror/volha-backend/internal/auth-service/service"
	gwconfig "github.com/autumnterror/volha-backend/internal/gateway/config"
	gwAuth "github.com/autumnterror/volha-backend/internal/gateway/grpc/auth"
	gwProducts "github.com/autumnterror/volha-backend/internal/gateway/grpc/products"
	gwEcho "github.com/autumnterror/volha-backend/internal/gateway/net/echo"
	gwRedis "github.com/autumnterror/volha-backend/internal/gateway/redis"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	productapi "github.com/autumnterror/volha-backend/internal/product-service/api"
	productconfig "github.com/autumnterror/volha-backend/internal/product-service/config"
	"github.com/autumnterror/volha-backend/internal/product-service/infra/psql/psqltx"
	productservice "github.com/autumnterror/volha-backend/internal/product-service/service"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Stack runs the real product and auth gRPC APIs over TCP with PostgreSQL,
// plus the real HTTP gateway and a Redis protocol server. It requires an
// explicitly named, migrated test database and never truncates shared tables.
type Stack struct {
	t          *testing.T
	DB         *sql.DB
	Products   pb.ProductsClient
	Auth       pb.AuthServiceClient
	BaseURL    string
	HTTP       *http.Client
	AdminToken string
	ids        map[string][]string
	close      []func()
}

type Catalog struct {
	Brand, Category, Country, Material, Color string
}

func New(t *testing.T) *Stack {
	t.Helper()
	dbName := os.Getenv("TEST_POSTGRES_DB")
	if dbName == "" || !strings.Contains(strings.ToLower(dbName), "test") {
		t.Fatal("set TEST_POSTGRES_DB to a migrated, dedicated test database")
	}
	cfg := productconfig.Test()
	db, err := sql.Open("postgres", cfg.ConnStr)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, db.PingContext(ctx), "start PostgreSQL and apply backend/migrations")
	var table string
	require.NoError(t, db.QueryRowContext(ctx, "SELECT to_regclass('public.articles')::text").Scan(&table), "apply backend/migrations before integration tests")

	s := &Stack{t: t, DB: db, ids: make(map[string][]string), AdminToken: "integration-admin-token"}
	s.close = append(s.close, func() { _ = db.Close() })
	t.Cleanup(s.shutdown)

	productLogic := productservice.NewProductsService(psqltx.NewTxRunner(db), psqltx.NewRepoProvider(db))
	productAddr := s.startGRPC(func(server *grpc.Server) { productapi.Register(server, productLogic) })
	authLogic := authservice.New(authrepo.New(db), strings.Repeat("integration-secret-", 3), 15*time.Minute, 24*time.Hour)
	authAddr := s.startGRPC(func(server *grpc.Server) { authapi.Register(server, authLogic) })
	s.Products = s.dialProducts(productAddr)
	s.Auth = s.dialAuth(authAddr)

	redisServer, err := miniredis.Run()
	require.NoError(t, err)
	s.close = append(s.close, redisServer.Close)
	metrics := observability.NewMetrics("integration-gateway", "test", "test")
	gwCfg := &gwconfig.Config{
		AddrProducts: productAddr, AddrAuth: authAddr, Timeout: 3 * time.Second,
		Mode: "DEV", DevAdminToken: s.AdminToken, Port: freePort(t),
		RedisAddr: redisServer.Addr(), RedisDictTTL: time.Minute,
		CORSAllowedOrigins: []string{"http://localhost:3000"},
	}
	gwProduct, err := gwProducts.New(gwCfg, observability.Config{Service: "integration-gateway"}, metrics)
	require.NoError(t, err)
	s.close = append(s.close, func() { _ = gwProduct.Close() })
	gwAuthentication, err := gwAuth.New(gwCfg, observability.Config{Service: "integration-gateway"}, metrics)
	require.NoError(t, err)
	s.close = append(s.close, func() { _ = gwAuthentication.Close() })
	cache, err := gwRedis.New(gwCfg, metrics)
	require.NoError(t, err)
	s.close = append(s.close, func() { _ = cache.Close() })

	// File endpoints use relative image/PDF directories. No other test in these
	// packages runs in parallel with a Stack, so changing cwd is safe here.
	t.Chdir(t.TempDir())
	gw := gwEcho.New(cache, gwProduct, gwAuthentication, gwCfg, observability.Config{Service: "integration-gateway"}, metrics)
	s.BaseURL = fmt.Sprintf("http://127.0.0.1:%d", gwCfg.Port)
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	s.HTTP = &http.Client{Jar: jar, Timeout: 5 * time.Second}
	serverErr := make(chan error, 1)
	go func() { serverErr <- gw.Run() }()
	s.close = append(s.close, func() {
		stopCtx, stop := context.WithTimeout(context.Background(), 3*time.Second)
		defer stop()
		_ = gw.Stop(stopCtx)
		select {
		case <-serverErr:
		case <-time.After(3 * time.Second):
			t.Error("gateway did not stop")
		}
	})
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		response, err := s.HTTP.Get(s.BaseURL + "/health/live")
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return s
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("gateway did not become healthy")
	return nil
}

func (s *Stack) startGRPC(register func(*grpc.Server)) string {
	s.t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.t, err)
	server := grpc.NewServer()
	register(server)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(server, healthServer)
	go func() { _ = server.Serve(listener) }()
	s.close = append(s.close, server.Stop)
	return listener.Addr().String()
}

func (s *Stack) dialProducts(addr string) pb.ProductsClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(s.t, err)
	s.close = append(s.close, func() { _ = conn.Close() })
	return pb.NewProductsClient(conn)
}

func (s *Stack) dialAuth(addr string) pb.AuthServiceClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(s.t, err)
	s.close = append(s.close, func() { _ = conn.Close() })
	return pb.NewAuthServiceClient(conn)
}

func freePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	return port
}

func (s *Stack) Track(table, id string) {
	s.t.Helper()
	switch table {
	case "products", "articles", "slides", "users", "brands", "categories", "countries", "materials", "colors":
		s.ids[table] = append(s.ids[table], id)
	default:
		s.t.Fatalf("unsupported cleanup table %q", table)
	}
}

func (s *Stack) shutdown() {
	for _, table := range []string{"products", "articles", "slides", "users", "brands", "categories", "countries", "materials", "colors"} {
		for _, id := range s.ids[table] {
			if _, err := s.DB.Exec("DELETE FROM "+table+" WHERE id = $1", id); err != nil {
				s.t.Errorf("cleanup %s %s: %v", table, id, err)
			}
		}
	}
	for i := len(s.close) - 1; i >= 0; i-- {
		s.close[i]()
	}
}

func ID() string { return xid.New().String() }

func SKU(t *testing.T) string {
	t.Helper()
	n, err := rand.Int(rand.Reader, big.NewInt(100000000))
	require.NoError(t, err)
	return fmt.Sprintf("%08d", n.Int64())
}

func (s *Stack) SeedCatalog() Catalog {
	s.t.Helper()
	ctx := context.Background()
	id := ID()
	c := Catalog{Brand: ID(), Category: ID(), Country: ID(), Material: ID(), Color: ID()}
	_, err := s.Products.CreateBrand(ctx, &pb.Brand{Id: c.Brand, Title: "Brand " + id})
	require.NoError(s.t, err)
	s.Track("brands", c.Brand)
	_, err = s.Products.CreateCategory(ctx, &pb.Category{Id: c.Category, Title: "Category " + id, Slug: "category-" + id})
	require.NoError(s.t, err)
	s.Track("categories", c.Category)
	_, err = s.Products.CreateCountry(ctx, &pb.Country{Id: c.Country, Title: "Country " + id})
	require.NoError(s.t, err)
	s.Track("countries", c.Country)
	_, err = s.Products.CreateMaterial(ctx, &pb.Material{Id: c.Material, Title: "Material " + id})
	require.NoError(s.t, err)
	s.Track("materials", c.Material)
	_, err = s.Products.CreateColor(ctx, &pb.Color{Id: c.Color, Title: "Color " + id, Hex: "#112233"})
	require.NoError(s.t, err)
	s.Track("colors", c.Color)
	return c
}

func (s *Stack) SeedProduct(c Catalog, photo string) *pb.ProductId {
	s.t.Helper()
	id := ID()
	p := &pb.ProductId{Id: id, Title: "Product " + id, Slug: "product-" + id, Article: SKU(s.t),
		Brand: c.Brand, Category: c.Category, Country: c.Country,
		Width: 10, Height: 20, Depth: 30, Price: 100, Pricesale: 80,
		Materials: []string{c.Material}, Colors: []string{c.Color}, Photos: []string{photo}}
	_, err := s.Products.CreateProduct(context.Background(), p)
	require.NoError(s.t, err)
	s.Track("products", id)
	return p
}

func (s *Stack) JSON(method, path string, body any, admin bool) (int, []byte) {
	s.t.Helper()
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		require.NoError(s.t, err)
		payload = bytes.NewReader(encoded)
	}
	req, err := http.NewRequest(method, s.BaseURL+path, payload)
	require.NoError(s.t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if admin {
		req.Header.Set("Authorization", "Bearer "+s.AdminToken)
	}
	response, err := s.HTTP.Do(req)
	require.NoError(s.t, err)
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	require.NoError(s.t, err)
	return response.StatusCode, data
}
