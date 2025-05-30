# --- Build stage ---
FROM golang:1.24.3-alpine AS builder

WORKDIR /app
COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
ENV CGO_ENABLED=0
RUN go build -o s3-proxy .

# --- Final stage ---
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/s3-proxy /s3-proxy
USER nonroot:nonroot

EXPOSE 8080
ENTRYPOINT ["/s3-proxy"]
