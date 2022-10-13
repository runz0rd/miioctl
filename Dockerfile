FROM golang:alpine AS builder

WORKDIR /build
COPY . .
RUN go build -o /miioctl cmd/main.go

FROM scratch
COPY --from=builder /miioctl /miioctl
