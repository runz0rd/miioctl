FROM python:slim
RUN apt-get update -y && apt-get install golang -y
RUN pip3 install python-miio

WORKDIR /build
COPY . .
RUN go build -o /miioctl cmd/main.go

ENTRYPOINT ["/miioctl"]
