FROM golang:alpine AS builder

WORKDIR /server

COPY . .

RUN go get -d -v

RUN go build -o app


FROM alpine:3

WORKDIR /server

COPY --from=builder /server/app app
COPY --from=builder /server/templates templates

ENV PORT 3000
EXPOSE 3000
ENV GIN_MODE release

ENTRYPOINT ["/server/app"]