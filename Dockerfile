FROM golang as builder
WORKDIR /go/src/github.com/jstorer/gannett
COPY main.go .
COPY /api ./api
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/src/github.com/jstorer/gannett/app .
CMD ["./app"]