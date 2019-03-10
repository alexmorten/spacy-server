FROM golang:1.11 as builder
WORKDIR /go/src/github.com/alexmorten/spacy-server
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o spacy-server main/spacy.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/alexmorten/spacy-server/spacy-server .
EXPOSE 4000
CMD ["./spacy-server"]
