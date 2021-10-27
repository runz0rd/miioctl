FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o /miioctl .

FROM python:alpine
RUN apk add --no-cache gcc libc-dev libffi-dev libressl-dev git make musl-dev go
RUN pip3 install python-miio
COPY --from=builder /miioctl /miioctl
RUN chmod u+x /miioctl
ENTRYPOINT ["/miioctl"]

