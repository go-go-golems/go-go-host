# syntax=docker/dockerfile:1
FROM golang:1.24-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/go-go-hostd ./cmd/go-go-hostd \
 && CGO_ENABLED=0 GOOS=linux go build -o /out/go-go-host ./cmd/go-go-host \
 && CGO_ENABLED=0 GOOS=linux go build -o /out/go-go-host-agent ./cmd/go-go-host-agent

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/go-go-hostd /go-go-hostd
COPY --from=build /out/go-go-host /go-go-host
COPY --from=build /out/go-go-host-agent /go-go-host-agent
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/go-go-hostd"]
