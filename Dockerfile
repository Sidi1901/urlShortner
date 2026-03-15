#Stage 1
FROM golang:alpine AS builder 

WORKDIR /build

COPY . .

RUN go build -o main

#Stage 2

FROM alpine:latest


WORKDIR /app 

COPY --from=builder /build/main /app

EXPOSE 3000 

ENTRYPOINT ["./main"]