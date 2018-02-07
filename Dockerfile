FROM ubuntu:16.04

MAINTAINER Lucas Teske <lucas@contaquanto.com.br>

COPY RemoteSigner /usr/bin/
RUN chmod +x /usr/bin/RemoteSigner
RUN mkdir /keys
RUN ln -s /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so
ENV PRIVATE_KEY_FOLDER /keys
WORKDIR /

CMD /usr/bin/RemoteSigner

