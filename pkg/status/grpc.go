package status

import (
	"context"

	"google.golang.org/grpc"
)

// RPC struct for the service
type RPC struct {
	service ServiceInterface
}

// Check defines if the sidecarred service is healthy
func (rpc *RPC) Check(ctx context.Context, r *HealthCheckRequest) (*HealthCheckResponse, error) {
	var status HealthCheckResponse_ServingStatus

	if rpc.service.ElaborateHealthStatus() {
		status = HealthCheckResponse_SERVING
	} else {
		status = HealthCheckResponse_NOT_SERVING
	}

	return &HealthCheckResponse{Status: status}, nil
}

// Update is the place to dump lag metric data
func (rpc *RPC) Update(ctx context.Context, r *KafkaMetricRequest) (*KafkaMetricResponse, error) {
	rpc.service.StoreLag(r.GetLag())
	return &KafkaMetricResponse{Accepted: true}, nil
}

// Register wires up the RPC to the server
func (rpc *RPC) Register(server *grpc.Server) {
	RegisterHealthServer(server, rpc)
	RegisterKafkaMetricServer(server, rpc)
}

// NewRPC returns a new initialized RPC object.
func NewRPC(svc ServiceInterface) RPCInterface {
	rpc := RPC{}
	rpc.service = svc

	return &rpc
}
