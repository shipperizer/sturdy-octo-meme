FROM golang:1.14 AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO_BIN=/go/bin/app

RUN apt update
RUN apt install -y awscli
RUN GRPC_HEALTH_PROBE_VERSION=v0.1.0-alpha.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /var/app

COPY . .

ARG app_name=grpc

ENV APP_NAME=$app_name

RUN make build

FROM gcr.io/distroless/base

COPY --from=builder /go/bin/app /app
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

CMD ["/app"]
