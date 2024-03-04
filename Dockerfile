# Build the application binary
FROM golang:alpine AS build-stage
WORKDIR /build
COPY . .
RUN go build -ldflags "-s -w" -trimpath -o ddns-server

# Deploy the application binary into a lean image
FROM alpine:latest AS release-stage
COPY --from=build-stage /build/ddns-server /usr/bin/
WORKDIR /data
COPY config.yaml .
ENTRYPOINT ["ddns-server"]
