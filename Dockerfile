FROM golang as builder

WORKDIR app

COPY . .

RUN "go build -o 'bin' cmd/bot/main"

FROM alpine as runner

WORKDIR app

COPY --from=builder ./app/bin bin

CMD ["bin", "-token=$TOKEN"]