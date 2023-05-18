FROM golang:1.18.1-alpine as builder

# Make Build dir
WORKDIR /go/src/github.com/nacrytchuk/item-service

# Copy golang dependency manifests
COPY go.mod go.sum ./

# Cache the downloaded dependency in the layer.
RUN go mod download

# add the source code
COPY . .
# Build
RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o item-service cmd/main.go

FROM scratch as runner

COPY --from=builder /go/src/github.com/nacrytchuk/item-service .

ENTRYPOINT ["/item-service"]
