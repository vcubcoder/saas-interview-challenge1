
FROM golang:1.11.4-alpine as build
WORKDIR /saas-interview-challenge1
ADD . /saas-interview-challenge1
RUN apk add git 
RUN CGO_ENABLED=0 GOOS=linux go build -a -o saasinterviewapp cmd/main.go

FROM alpine:latest  
WORKDIR /saas-interview-challenge1
COPY --from=build /saas-interview-challenge1/saasinterviewapp .
ENTRYPOINT ./saasinterviewapp --redis-host="${REDIS_HOST}"
         