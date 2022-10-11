FROM golang:alpine
RUN apk add --no-cache python3 py3-pip py3-netifaces gcc g++ make libffi-dev openssl-dev
RUN pip3 install python-miio

WORKDIR /build
COPY . .
RUN go build -o /miioctl cmd/main.go

ENTRYPOINT ["/miioctl"]
