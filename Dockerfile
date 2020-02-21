FROM golang:alpine AS development
WORKDIR $GOPATH/src/edp-fss
COPY . .
RUN go build -o edp-fss

FROM alpine:latest AS production
WORKDIR /root/
COPY --from=development /go/src/edp-fss .
EXPOSE 8080
ENTRYPOINT ["nohup ./edp-fss &"]