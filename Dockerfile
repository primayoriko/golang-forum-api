FROM golang:alpine

# RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main

# EXPOSE 8080 8080

ENTRYPOINT ["/app/main"]
