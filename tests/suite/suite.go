package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/TiZir/gRPC_service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc_servicev1 "github.com/TiZir/gRPC_service_protos/protos/gen/go/gRPC_service"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient grpc_servicev1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()
	cfg := config.MustLoadByPath("../config/local.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create gRPC client: %v", err)
	}
	
	return ctx, &Suite{
		T:         t,
        Cfg:        cfg,
        AuthClient: grpc_servicev1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}