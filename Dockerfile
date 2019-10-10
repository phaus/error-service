FROM golang:alpine3.8
RUN apk --no-cache update && apk add --upgrade git 
WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/service

FROM alpine:3.8
RUN apk --no-cache update

WORKDIR /
ENTRYPOINT ["/app/server"]
COPY --from=0 /app/bin/service /app/server
