# sturdy-octo-meme

![ci-tests](https://github.com/shipperizer/sturdy-octo-meme/workflows/Unit%20Tests/badge.svg)
[![codecov](https://codecov.io/gh/shipperizer/sturdy-octo-meme/branch/master/graph/badge.svg)](https://codecov.io/gh/shipperizer/sturdy-octo-meme)

GRPC sidecar app for headless kafka consumers, it provides the standard Health GRPC check and an Update GRPC that is used to update the metrics used to determine the status of the kafka consumer.


### Test behaviour

To test it, simply build the image and run it, then use [`grpc_cli`](https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md) and point it to the `proto` files


```
shipperizer@arbalester  on  master!23:30:01 π       docker build -t health-kafka-sidecar:latest .
shipperizer@arbalester  on  master!23:32:01 π       docker run -it -p 18000:8000 -e LOG_LEVEL=DEBUG health-kafka-sidecar:latest
shipperizer@arbalester  on  master!23:32:07 π      grpc_cli call --noremotedb --proto_path pkg/status --protofiles update.proto localhost:18000 KafkaMetric.Update 'lag: 1000'
connecting to localhost:18000
accepted: true
Rpc succeeded with OK status
shipperizer@arbalester  on  master!23:32:28 π   grpc_cli call --noremotedb --proto_path pkg/status --protofiles status.proto localhost:18000 Health.Check 'service: "1000"'
connecting to localhost:18000
status: SERVING
Rpc succeeded with OK status
```

### Kubernetes integration

The idea is to run it alongside other containers in the same pod:

```
containers:
- name: health-kafka-sidecar
  image: health-kafka-sidecar:latest
  imagePullPolicy: IfNotPresent
  ports:
    - containerPort: 8000
      name: grpc
  livenessProbe:
    exec:
      command: ["/bin/grpc_health_probe", "-addr=:8000"]
    initialDelaySeconds: 5
    failureThreshold: 10
    timeoutSeconds: 5
    periodSeconds: 10
  readinessProbe:
    exec:
      command: ["/bin/grpc_health_probe", "-addr=:8000"]
    initialDelaySeconds: 5
    failureThreshold: 10
    timeoutSeconds: 5
    periodSeconds: 10
  lifecycle:
    preStop:
      exec:
        command: ["bash", "-c" , "sleep 5"]
  env:
    - name: LOG_LEVEL
      value: DEBUG # optional
```       

while the other container(s) will simply point at it on the pod ip:, see snippet below to use downstream api to get the pod id inside the other contianers:

```
env:
  - name: POD_IP
    valueFrom:
      fieldRef:
        fieldPath: status.podIP
```


### Integration example

Gist of if is to push Topic Lag metrics tot he Update RPC and the have the Health RPC to use those metrics to evaluate if the consumer is working or not.
Below an example of a quick integration:

```
package main

import (
	"context"
	"fmt"

	status "github.com/shipperizer/sturdy-octo-meme/pkg/status"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	var healthGRPC status.KafkaMetricClient
	conn, err := grpc.DialContext(
		context.Background(),
		":8000",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Error("Error setting up Health GRPC client: ", err.Error())
		return
	}

  healthGRPC = status.NewKafkaMetricClient(conn)
  log.Info("All good with GRPC setup: ", healthGRPC)
  defer conn.Close()

	r, e := healthGRPC.Update(context.Background(), &status.KafkaMetricRequest{Lag: int64(1000)})

	if e != nil {
		log.Error(e)
	}

	log.Info(r)
}
```   
