#Build stage
FROM golang:1.17.2-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go

#Run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
RUN apk add --no-cache bash
EXPOSE 7000
CMD ./main