FROM golang:1.7-wheezy
MAINTAINER Maxim Kupriianov <max@kc.vc>

RUN apt-get update -q \
	&& DEBIAN_FRONTEND=noninteractive apt-get install -qy pkg-config python2.7-dev \
	&& apt-get clean \
	&& rm -rf /var/lib/apt

CMD /usr/local/go/bin/go
