FROM alpine:3.6
MAINTAINER Dalton Hubble <dalton.hubble@coreos.com>
COPY bin/matchbox /matchbox
COPY bin/matchboxd /matchboxd
EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["/matchboxd"]
