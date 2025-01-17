FROM alpine:3.19.0

COPY .env .

ADD survey /

CMD ["/survey"]