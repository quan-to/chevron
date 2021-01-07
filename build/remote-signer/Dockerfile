FROM golang:1.14-alpine as build

RUN apk update

RUN apk add git ca-certificates

ADD . /go/src/github.com/quan-to/chevron

ENV GO111MODULE=on

# Compile Server
WORKDIR /go/src/github.com/quan-to/chevron/cmd/server
RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../remote-signer

# Compile Standalone
WORKDIR /go/src/github.com/quan-to/chevron/cmd/cli
RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../standalone

# Compile DBMigrate
WORKDIR /go/src/github.com/quan-to/chevron/cmd/dbmigrate
RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../dbmigrate

FROM alpine:latest

MAINTAINER Lucas Teske <lucas@contaquanto.com.br>

RUN apk --no-cache add ca-certificates

RUN mkdir -p /opt/remote-signer/
WORKDIR /opt/remote-signer

COPY --from=build /go/src/github.com/quan-to/chevron/remote-signer .
COPY --from=build /go/src/github.com/quan-to/chevron/standalone .
COPY --from=build /go/src/github.com/quan-to/chevron/dbmigrate .

RUN mkdir -p /keys

VOLUME ["/keys"]

EXPOSE "5100"

# Common Environment

ENV HTTP_PORT "5100"
ENV PRIVATE_KEY_FOLDER "/keys"
ENV SYSLOG_IP "127.0.0.1"
ENV SYSLOG_FACILITY "LOG_USER"
ENV SKS_SERVER "http://sks:11371"
ENV KEY_PREFIX ""
ENV MAX_KEYRING_CACHE_SIZE "1000"
ENV KEYS_BASE64_ENCODED "true"
ENV READONLY_KEYPATH "false"

# Cluster Mode
ENV MASTER_GPG_KEY_PATH ""
ENV MASTER_GPG_KEY_PASSWORD_PATH ""
ENV MASTER_GPG_KEY_BASE64_ENCODED "true"


# Hashicorp Vault
ENV VAULT_ADDRESS ""
ENV VAULT_ROOT_TOKEN ""
ENV VAULT_PATH_PREFIX ""
ENV VAULT_STORAGE "false"
ENV VAULT_SKIP_VERIFY "false"

# Agent Vars
ENV AGENT_TARGET_URL ""
ENV AGENT_KEY_FINGERPRINT ""
ENV AGENT_BYPASS_LOGIN "false"
ENV AGENT_EXTERNAL_URL "/agent"
ENV AGENTADMIN_EXTERNAL_URL "/agentAdmin"


# Database
ENV ENABLE_DATABASE "false"
ENV DATABASE_NAME "remote_signer"
ENV DATABASE_TOKEN_MANAGER "false"
ENV DATABASE_AUTH_MANAGER "false"
ENV DATABASE_DIALECT ""

# Redis Caching
ENV REDIS_ENABLE "false"
ENV REDIS_TLS_ENABLE "false"
ENV REDIS_HOST "localhost:6379"
ENV REDIS_PASS ""
ENV REDIS_USER ""
ENV REDIS_MAX_LOCAL_TTL "5m"
ENV REDIS_MAX_LOCAL_OBJECTS "100"
ENV REDIS_CLUSTER_MODE "false"

CMD /opt/remote-signer/remote-signer
