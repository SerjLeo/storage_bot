FROM golang:1.20-alpine3.16 as builder

RUN apk update && apk add --no-cache gcc musl-dev

COPY . /github.com/SerjLeo/storage_bot/
WORKDIR /github.com/SerjLeo/storage_bot/

RUN go mod download && go get -u ./...
RUN CGO_ENABLED=1 go build -o "bin/main" cmd/bot/main.go

FROM alpine:3.16 as runner

WORKDIR /root/

COPY --from=builder /github.com/SerjLeo/storage_bot/bin/main .
RUN mkdir data && cd data && mkdir sqlite

ENTRYPOINT ["./main"]