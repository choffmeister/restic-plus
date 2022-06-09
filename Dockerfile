FROM alpine:3.16
RUN apk add --no-cache ca-certificates
COPY restic-plus /bin/restic-plus
ENTRYPOINT ["/bin/restic-plus"]
