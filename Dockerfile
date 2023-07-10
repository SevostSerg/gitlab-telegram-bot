FROM golang:alpine AS builder

WORKDIR /app

ADD go.mod .

COPY ./ ./

RUN go mod download

ENV storage=""

RUN go build -o .

FROM alpine

WORKDIR /app

COPY --from=builder /app/GitlabTgBot /app/GitlabTgBot

EXPOSE 8080

CMD [ "/app/GitlabTgBot" ]
