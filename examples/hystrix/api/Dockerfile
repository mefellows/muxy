FROM golang:1.13-alpine
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o api

ENV PORT 80
EXPOSE 80

CMD ["/go/src/app/api"]