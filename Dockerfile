FROM golang:1.16 as builder
WORKDIR /go/src/app
COPY . .
RUN make clean build

FROM debian:buster-slim
COPY --from=builder /go/src/app/bin/waterbucket .
CMD ["./waterbucket"]
