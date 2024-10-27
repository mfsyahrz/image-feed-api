FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o image_feed_api .

FROM alpine:latest

WORKDIR /root/

RUN mkdir migration
COPY --from=builder /app/image_feed_api .
COPY --from=builder /app/migration/ migration/       

EXPOSE 3333 3333

CMD ["./image_feed_api"]