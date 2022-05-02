FROM golang:1.16.4

WORKDIR /gefence-service
COPY . .
ENTRYPOINT ["./geofence-service"]
