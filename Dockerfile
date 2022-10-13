FROM golang:alpine AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o /miioctl cmd/main.go

FROM scratch
COPY --from=builder /miioctl /miioctl
ENTRYPOINT ["/miioctl"]