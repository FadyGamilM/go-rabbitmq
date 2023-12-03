FROM golang:1.21-alpine AS build

COPY go.mod go.sum /app/

WORKDIR /app/

RUN go mod download

COPY ./producer/producer.go /app/

RUN go build -o /app/producer .

###
FROM alpine:3.14

ENV PORT=9876
ENV RABBIT_HOST=rabbit-1
ENV RABBIT_PORT=5672
ENV RABBIT_USERNAME=guest
ENV RABBIT_PASSWORD=guest

EXPOSE ${PORT}

# i will use CMD not ENTRYPOINT so i have no modification tolerant later via the terminal args
CMD [ "/app/producer" ]

COPY --from=build /app/producer /app/producer
RUN chmod +x /app/producer
