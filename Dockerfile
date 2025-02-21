FROM golang:latest as base

COPY . .

RUN go build -o /sfr .

FROM scratch

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /sfr .


EXPOSE 8080

CMD ["./sfr"]