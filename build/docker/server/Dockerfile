FROM golang:1.15 as builder

RUN mkdir /server
WORKDIR /app
ADD . .
RUN make build-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root
COPY --from=builder /app/out/bin/server .
CMD ["./server"]
