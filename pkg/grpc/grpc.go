package grpc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/shipperizer/sturdy-octo-meme/pkg/status"
)

// InitGRPC initializes RPCs
func InitGRPC() *grpc.Server {
	server := grpc.NewServer()

	status.NewRPC(
		status.NewService(
			viper.GetInt("health.lag.history"),
			viper.GetInt("health.lag.max"),
		),
	).Register(server)

	return server
}

// NewGRPC creates a GRPC server
func NewGRPC() *grpc.Server {
	return InitGRPC()
}
