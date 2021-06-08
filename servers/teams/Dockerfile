FROM alpine
RUN apk add --no-cache ca-certificates
COPY teams /teams
EXPOSE 80
ENTRYPOINT ["/teams"]