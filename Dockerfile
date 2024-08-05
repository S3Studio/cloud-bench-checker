# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v ./bin/cmd/main.go

# Deploy the application binary into a clean image
FROM alpine AS release-stage

WORKDIR /app
RUN adduser -D nonroot && chown nonroot /app
USER nonroot:nonroot

COPY --chown=nonroot --from=build-stage /app/main /app/main
RUN mkdir output && touch output/output.csv && ln -s output/output.csv test.csv

ENTRYPOINT ["/app/main", "-c", "config.conf"]
