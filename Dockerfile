FROM alpine:3.7 as base_go

ARG ARG_GITHUB_ACCESS_TOKEN

ADD . /go/src/github.com/veritone/task-status-update

RUN apk update && \
    apk add -U build-base go git curl libstdc++ && \
                cd /go/src/github.com/veritone/task-conductor-v3.1 && \
    git config --global url."https://${ARG_GITHUB_ACCESS_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/" && \
                go env && go list all | grep cover && \
    GOPATH=/go make -f Makefile.inContainer docker && \
    mkdir /app && \
    mv /go/src/github.com/veritone/task-status-update/task-status-update /app && \
    mv /go/src/github.com/veritone/task-status-update/build-manifest.yml /app && \
    rm -rf /go && apk del build-base go git libstdc++

FROM alpine:3.7

RUN mkdir /app && apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=base_go /app /app

WORKDIR /app
ENTRYPOINT ["/app/task-status-update"]
