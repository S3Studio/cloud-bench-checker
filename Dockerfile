# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -ldflags="-w -s" -v ./bin/cmd/main.go

# Deploy the application binary into a clean image
FROM alpine AS build-release-stage

WORKDIR /app

COPY --from=build-stage main main

USER nonroot:nonroot

ENTRYPOINT ["/app/main"]
