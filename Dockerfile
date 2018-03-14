FROM mono:5.8

MAINTAINER Lucas Teske <lucas@contaquanto.com.br>

#ARG DEBIAN_FRONTEND=noninteractive
#RUN apt-get update && apt-get install -y --no-install-recommends mono-complete && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/remote-signer/
RUN mkdir -p /tmp/remote-signer/
RUN mkdir -p /keys
RUN ln -s /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so
ENV PRIVATE_KEY_FOLDER /keys

COPY ./ /tmp/remote-signer
WORKDIR /tmp/remote-signer
RUN ./build-nix.sh
RUN cp ./RemoteSigner/bin/Release/* /opt/remote-signer/
WORKDIR /
RUN rm -fr /tmp/remote-signer

CMD /usr/bin/mono /opt/remote-signer/RemoteSigner.exe

