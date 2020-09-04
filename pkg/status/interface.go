// Interface definitions
package status

import (
	"context"

	"google.golang.org/grpc"
)

// RPCInterface interface for RPC.
type RPCInterface interface {
	Register(*grpc.Server)
	Check(ctx context.Context, r *HealthCheckRequest) (*HealthCheckResponse, error)
	Update(ctx context.Context, r *KafkaMetricRequest) (*KafkaMetricResponse, error)
}

// ServiceInterface interface for service provider.
type ServiceInterface interface {
	StoreLag(lag int64)
	ElaborateHealthStatus() bool
}
