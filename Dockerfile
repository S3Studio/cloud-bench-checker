# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -ldflags="-w -s" -v ./bin/cmd/main.go

# Deploy the application binary into a clean image
FROM golang:1.22-alpine AS release-stage

WORKDIR /app
RUN adduser -D nonroot
USER nonroot:nonroot

COPY --chown=nonroot --from=build-stage /app/main /app/main

ENTRYPOINT ["/app/main"]
