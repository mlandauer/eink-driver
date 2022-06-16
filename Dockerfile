FROM gcc:8.4.0
# TODO: Use standard go image instead

# Install go as well
RUN wget https://dl.google.com/go/go1.14.1.linux-arm64.tar.gz
RUN tar -C /usr/local -xzf go1.14.1.linux-arm64.tar.gz
RUN rm go1.14.1.linux-arm64.tar.gz

RUN apt-get update && apt-get install -y chromium

WORKDIR /usr/src/myapp

COPY web/go.* web/
RUN cd web; /usr/local/go/bin/go mod download

COPY web/main.go web/
RUN cd web; /usr/local/go/bin/go build main.go

COPY finn.png .
ENV URL http://192.168.5.10/solar
CMD ["web/main"]
