# TODO: Don't use latest.
FROM alpine:latest

COPY bin/lander /lander

CMD ["./lander"]