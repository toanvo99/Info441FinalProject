FROM alpine
RUN apk add --no-cache ca-certificates
COPY messaging /messaging
EXPOSE 80
ENTRYPOINT ["/teams"]