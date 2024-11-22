FROM golang:1.23.1-alpine3.20

WORKDIR /app

COPY go.mod ./

COPY cmd /app/cmd
COPY config /app/config
COPY internal /app/internal
COPY test /app/test

RUN go mod tidy
RUN go build cmd/app/main.go

EXPOSE 4953

CMD [ "./main" ]