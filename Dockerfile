# FROM golang:1.13-alpine as builder
# COPY go.mod go.sum /go/src/gitlab.com/hydra/forum-api/
# WORKDIR /go/src/gitlab.com/hydra/forum-api
# RUN go mod download
# COPY . /go/src/gitlab.com/hydra/forum-api
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/forum-api gitlab.com/hydra/forum-api

# FROM alpine
# RUN apk add --no-cache ca-certificates && update-ca-certificates
# COPY --from=builder /go/src/gitlab.com/hydra/forum-api/build/forum-api /usr/bin/forum-api
# EXPOSE 8080 8080
# ENTRYPOINT ["/usr/bin/forum-api"]

FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main

ENTRYPOINT ["/app/main"]
