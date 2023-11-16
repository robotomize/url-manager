FROM golang:1.21.0 AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

COPY . .

RUN go build \
  -trimpath \
  -ldflags "-s -w -X main.BuildTag=$(git describe --tags --abbrev=0) -X main.BuildName=url-manager -extldflags '-static'" \
  -o /bin/url-manager \
  ./cmd/url-manager

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/url-manager /bin/url-manager

ENTRYPOINT ["/bin/url-manager"]