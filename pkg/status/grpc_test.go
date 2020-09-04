package status

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/shipperizer/sturdy-octo-meme/internal/config"
)

// Blueprint struct for db service and auth provider.
func TestingGRPCServer() (*grpc.Server, *bufconn.Listener) {
	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	return grpc.NewServer(), listener
}

func TestingBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestGRPCUpdate(t *testing.T) {
	config.Load()

	srv, listener := TestingGRPCServer()
	ctx := context.Background()

	NewRPC(
		NewService(1000, 10000),
	).Register(srv)

	go func() {
		err := srv.Serve(listener)
		if err != nil {
			t.Fatal("failed to start grpc server: ", err)
		} else {
		}
	}()
	defer srv.Stop()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(TestingBufDialer(listener)), grpc.WithInsecure())

	if err != nil {
		t.Fatal("failed to dial: ", err)
	}

	defer conn.Close()

	client := NewKafkaMetricClient(conn)

	resp, err := client.Update(ctx, &KafkaMetricRequest{Lag: int64(10000)})
	assert := assert.New(t)
	assert.Nil(err, "should have been nil")
	assert.True(resp.GetAccepted(), "accepted should have been true")
}

func TestGRPCCheck(t *testing.T) {
	config.Load()

	srv, listener := TestingGRPCServer()
	ctx := context.Background()

	NewRPC(
		NewService(1000, 10000),
	).Register(srv)

	go func() {
		err := srv.Serve(listener)
		if err != nil {
			t.Fatal("failed to start grpc server: ", err)
		} else {
		}
	}()
	defer srv.Stop()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(TestingBufDialer(listener)), grpc.WithInsecure())

	if err != nil {
		t.Fatal("failed to dial: ", err)
	}

	defer conn.Close()

	client := NewHealthClient(conn)

	resp, err := client.Check(ctx, &HealthCheckRequest{Service: "whatever"})
	assert := assert.New(t)
	assert.Nil(err, "should have been nil")
	assert.Equal(HealthCheckResponse_SERVING, resp.GetStatus(), "status should have been different")
}
