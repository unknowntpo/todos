package grpc

import (
	"context"
	"fmt"

	"github.com/unknowntpo/todos/internal/healthcheck/delivery/grpc/healthcheck"
	"github.com/unknowntpo/todos/internal/logger"

	"google.golang.org/grpc"
)

type healthcheckGrpcServer struct {
	healthcheck.UnimplementedHealthcheckServer
	logger  logger.Logger
	env     string
	version string
}

func NewHealthcheckGrpc(srv *grpc.Server, env, version string, l logger.Logger) {
	h := &healthcheckGrpcServer{}
	healthcheck.RegisterHealthcheckServer(srv, h)
}

// Healthcheck returns the status of our service.
func (h *healthcheckGrpcServer) Healthcheck(ctx context.Context, r *healthcheck.HealthcheckRequest) (*healthcheck.HealthcheckResponse, error) {
	h.logger.PrintInfo(fmt.Sprintf("healthcheckRequest: %s", r.String()), nil)

	resp := &healthcheck.HealthcheckResponse{
		Environment: h.env,
		Status:      "available",
		Version:     h.version,
	}

	return resp, nil
}
