FROM golang:latest as base

COPY . .

RUN go build -o /sfr .

EXPOSE 8080

CMD ["/sfr"]