FROM python:alpine
WORKDIR /build
RUN apk add --no-cache gcc libc-dev libffi-dev libressl-dev git make musl-dev go
RUN pip3 install python-miio
COPY . .
RUN go build -o /miioctl cmd/main.go
RUN chmod +x /miioctl
