FROM golang:1.11.4-alpine

ADD . /saas-interview-challenge1
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV REDIS_HOST=""
        
WORKDIR /saas-interview-challenge1

RUN apk add git 

ENTRYPOINT [ "go","run", "cmd/main.go","--redis-host", "${REDIS_HOST}]
# CMD ["go", "run","cmd/main.go", "--redis-host", "${REDIS_HOST}"]
