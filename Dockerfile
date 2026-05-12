# syntax=docker/dockerfile:1
FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/go-go-hostd ./cmd/go-go-hostd \
 && CGO_ENABLED=1 GOOS=linux go build -o /out/go-go-host ./cmd/go-go-host \
 && CGO_ENABLED=1 GOOS=linux go build -o /out/go-go-host-agent ./cmd/go-go-host-agent

FROM debian:bookworm-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates tzdata \
 && rm -rf /var/lib/apt/lists/* \
 && useradd --system --uid 10001 --home /var/lib/go-go-host --create-home go-go-host
WORKDIR /
COPY --from=build /out/go-go-hostd /go-go-hostd
COPY --from=build /out/go-go-host /go-go-host
COPY --from=build /out/go-go-host-agent /go-go-host-agent
USER go-go-host:go-go-host
EXPOSE 8080
ENTRYPOINT ["/go-go-hostd"]
