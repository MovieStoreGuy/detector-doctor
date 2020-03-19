FROM golang:1.14-alpine AS builder 

WORKDIR /src/detector-doctor 
COPY . . 

ENV CGO_ENABLED=0 GOOS=linux 

RUN set -x && \
    apk add --no-cache git && \
    go get ./... && \
    go test --cover -v ./... && \
    go build -o /detdoc \
        --ldflags="-s -w -X main.Version=`git describe --tags` -X main.GitHash=`git rev-parse --verify HEAD`" \
        -v *.go && \
    set +x

FROM alpine:3.9

LABEL Author='Sean (MovieStoreGuy) Marciniak'

COPY --from=builder /detdoc /detdoc 

WORKDIR /home/detdoc

RUN set -x && \
    apk add --no-cache ca-certificates bash && \
    addgroup -S detdoc && \
    adduser  -S detdoc -G detdoc && \
    chown detdoc:detdoc /home/detdoc && \
    set +x 

USER detdoc:detdoc
ENTRYPOINT [ "/detdoc" ]