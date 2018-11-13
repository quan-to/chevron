FROM mono:5.16

MAINTAINER Lucas Teske <lucas@contaquanto.com.br>

#ARG DEBIAN_FRONTEND=noninteractive
#RUN apt-get update && apt-get install -y --no-install-recommends mono-complete && rm -rf /var/lib/apt/lists/*

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

RUN mkdir -p /opt/remote-signer/
RUN mkdir -p /tmp/remote-signer/
RUN mkdir -p /keys
RUN ln -s /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so

COPY ./ /tmp/remote-signer
WORKDIR /tmp/remote-signer
RUN ./build-nix.sh
RUN cp ./RemoteSigner/bin/Release/* /opt/remote-signer/
WORKDIR /
RUN rm -fr /tmp/remote-signer

CMD /usr/bin/mono /opt/remote-signer/RemoteSigner.exe

