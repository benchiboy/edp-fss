
FROM alpine:latest AS production
WORKDIR /go/src/
COPY edp-fss .
COPY NotoSansSC-Regular.ttf .
COPY config-ci.json .

EXPOSE 8080
CMD [ "/go/src/edp-fss"]