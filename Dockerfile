FROM gcc:8.4.0

# Install go as well
RUN wget https://dl.google.com/go/go1.14.1.linux-armv6l.tar.gz
RUN tar -C /usr/local -xzf go1.14.1.linux-armv6l.tar.gz
RUN rm go1.14.1.linux-armv6l.tar.gz

RUN apt-get update && apt-get install -y chromium

WORKDIR /usr/src/myapp

COPY bcm2835-1.63 bcm2835-1.63
RUN cd bcm2835-1.63; ./configure; make; make install
RUN rm -rf bcm2835-1.63

COPY IT8951 IT8951
RUN cd IT8951; make clean; make

COPY web/go.* web/
RUN cd web; /usr/local/go/bin/go mod download

COPY web/main.go web/
RUN cd web; /usr/local/go/bin/go build main.go

ENV URL http://192.168.5.10/solar
CMD ["web/main"]
