# --- Build stage ---
FROM golang:1.24.3-alpine AS builder

WORKDIR /app
COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
ENV CGO_ENABLED=0
RUN go build -o s3-reverse-proxy .

# --- Final stage ---
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/s3-reverse-proxy /s3-reverse-proxy
USER nonroot:nonroot

EXPOSE 80
ENTRYPOINT ["/s3-reverse-proxy"]
