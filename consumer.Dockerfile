# Build Stage
FROM golang:1.21-alpine AS build

WORKDIR /app

# Copy only the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all Go files from pkg/rabbitmq and consumer directories
COPY ./pkg/rabbitmq/*.go ./pkg/rabbitmq/
COPY ./consumer/*.go ./consumer/

# Build the Go application
RUN go build -o /app/consumer/bin ./consumer
RUN chmod +x /app/consumer/bin

# Final Stage
FROM alpine:3.14

WORKDIR /app

ENV PORT=9876
ENV RABBIT_HOST=rabbit-1
ENV RABBIT_PORT=5672
ENV RABBIT_USERNAME=guest
ENV RABBIT_PASSWORD=guest

EXPOSE ${PORT}

# Copy the built executable from the previous stage
COPY --from=build /app/consumer/bin /app/consumer/

# Use CMD instead of ENTRYPOINT for flexibility
CMD [ "/app/consumer/bin" ]
