FROM golang AS builder
WORKDIR /go/src/sensor
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -a
FROM scratch
COPY --from=builder /go/src/sensor/sensor /sensor
CMD [ "/sensor" ]
