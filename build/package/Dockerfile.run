FROM alpine:3.9
MAINTAINER "Stakater Team"

RUN apk add --update ca-certificates

COPY Chowkidar /bin/Chowkidar

ENTRYPOINT ["/bin/Chowkidar"]
