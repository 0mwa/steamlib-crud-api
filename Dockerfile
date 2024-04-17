
# Create layer from golang image
FROM golang:1.22-alpine as builder

# Set work directory
WORKDIR /app

# Add files from current directory
COPY . .

# Build app, upgrade, ignore baselayout in final image, install required dependencies
RUN apk -U upgrade --ignore alpine-baselayout && apk add --no-cache \
    gcc \
    g++ \
    libxml2

# CGO_ENABLED - needed for system packages(os), GOOS & GOARCH - make sure Go compiles specificaly for Linux & amd64
# Command for Go compile
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -a -installsuffix cgo -o ./cmd/main ./cmd

# Create layer from alpine
FROM alpine

# Build app, upgrade, ignore baselayout in final image, install base tools
RUN apk -U upgrade --ignore alpine-baselayout && apk add --no-cache \
    curl \
    nano \
    vim \
    bash \
    tzdata

# Add timezone because -> tokens works with time expiration
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

# Add app form builder container
COPY --from=builder /app/cmd/main /app/bin/main

# Navigate to work directory
WORKDIR /app/bin

# Start app
CMD ["./main"]
