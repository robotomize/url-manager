FROM golang:1.21.0 AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

COPY . .

RUN go build \
  -trimpath \
  -ldflags "-s -w -X main.BuildTag=$(git describe --tags --abbrev=0) -X main.BuildName=urlmanager -extldflags '-static'" \
  -o /bin/urlmanager \
  ./cmd/url-manager

FROM scratch
COPY --from=builder /bin/urlmanager /bin/urlmanager

ENTRYPOINT ["/bin/urlmanager"]