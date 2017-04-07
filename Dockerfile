FROM golang

RUN apt-get install -y git

RUN go get github.com/czerwonk/udp-mirror

ENTRYPOINT [ "udp-mirror", "-receivers" ]
EXPOSE 9999
