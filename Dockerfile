FROM golang:1.12-alpine as build

RUN apk update

RUN apk add git ca-certificates

ADD . /go/src/github.com/quan-to/chevron

ENV GO111MODULE=on

# Compile Server
WORKDIR /go/src/github.com/quan-to/chevron/cmd/server
RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../remote-signer

# Compile Standalone
WORKDIR /go/src/github.com/quan-to/chevron/cmd/standalone
RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../standalone


FROM alpine:latest

MAINTAINER Lucas Teske <lucas@contaquanto.com.br>

RUN apk --no-cache add ca-certificates

RUN mkdir -p /opt/remote-signer/
WORKDIR /opt/remote-signer

COPY --from=build /go/src/github.com/quan-to/chevron/remote-signer .
COPY --from=build /go/src/github.com/quan-to/chevron/standalone .

RUN mkdir -p /keys

VOLUME ["/keys"]

EXPOSE "5100"
ENV HTTP_PORT "5100"
ENV PRIVATE_KEY_FOLDER /keys
ENV SYSLOG_IP "127.0.0.1"
ENV SYSLOG_FACILITY "LOG_USER"
ENV SKS_SERVER "http://sks:11371"
ENV KEY_PREFIX ""
ENV MAX_KEYRING_CACHE_SIZE "1000"
ENV ENABLE_RETHINKDB_SKS "false"
ENV RETHINKDB_HOST "rethinkdb"
ENV RETHINKDB_USERNAME "admin"
ENV RETHINKDB_PASSWORD ""
ENV RETHINKDB_PORT "28015"
ENV RETHINKDB_POOL_SIZE "10"
ENV DATABASE_NAME "remote_signer"
ENV MASTER_GPG_KEY_PATH ""
ENV MASTER_GPG_KEY_PASSWORD_PATH ""
ENV MASTER_GPG_KEY_BASE64_ENCODED "true"
ENV KEYS_BASE64_ENCODED "true"
ENV VAULT_ADDRESS ""
ENV VAULT_ROOT_TOKEN ""
ENV VAULT_PATH_PREFIX ""
ENV VAULT_STORAGE "false"
ENV READONLY_KEYPATH "false"
ENV VAULT_SKIP_VERIFY "false"
ENV AGENT_TARGET_URL ""
ENV AGENT_KEY_FINGERPRINT ""
ENV AGENT_BYPASS_LOGIN "false"
ENV RETHINK_TOKEN_MANAGER "false"
ENV RETHINK_AUTH_MANAGER "false"

CMD /opt/remote-signer/remote-signer

