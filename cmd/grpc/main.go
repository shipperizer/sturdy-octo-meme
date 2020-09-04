package main

import (
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/shipperizer/sturdy-octo-meme/internal/config"
	"github.com/shipperizer/sturdy-octo-meme/pkg/grpc"
)

func main() {
	// TODO @shipperizer create a bootstrap pkg where config and logging are init steps
	config.Load()
	config.SetupLogger(viper.GetString("logging.level"))

	// AWS NLB TLS termination.
	listener, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatal("Error: ", err)
		return
	}

	api := grpc.NewGRPC()

	api.Serve(listener)
}
