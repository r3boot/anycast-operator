FROM alpine:3.7
MAINTAINER Lex van Roon <r3boot@r3blog.nl>

RUN apk update \
    && apk upgrade

COPY build/anycast-operator /usr/sbin/anycast-operator
COPY files/run_anycast-operator /run_anycast-operator

ENTRYPOINT ["/run_anycast-operator"]
CMD [""]