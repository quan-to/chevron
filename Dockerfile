FROM ubuntu:16.04

MAINTAINER Lucas Teske <lucas@contaquanto.com.br>

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends mono-complete && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/remote-signer/

COPY ./RemoteSigner/bin/Release/ /opt/remote-signer/
RUN mkdir /keys
RUN ln -s /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so
ENV PRIVATE_KEY_FOLDER /keys
WORKDIR /

CMD /usr/bin/mono /opt/remote-signer/RemoteSigner.exe

