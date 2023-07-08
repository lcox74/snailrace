# Builder
FROM golang:1.18-alpine AS builder

RUN apk update && apk add git gcc musl-dev

# Set working directory
COPY ./ /src
WORKDIR /src

# Get packages and Build
RUN go mod download
RUN GOOS=linux go build -o /bin/snailrace

# Runner
FROM alpine:latest as deploy
LABEL description="Sailrace" Version="0.0.1"

# Copy the binary from the builder
COPY --from=builder /bin/snailrace /bin/snailrace

# Expose port 3000 to the outside world
EXPOSE 3000
RUN chmod +x /bin/snailrace

RUN mkdir -p /bin/res

COPY .env /bin/.env
COPY res/* /bin/res
WORKDIR /bin

ENTRYPOINT ["/bin/snailrace"]