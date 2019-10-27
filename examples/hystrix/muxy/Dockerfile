FROM golang:1.13-alpine

RUN mkdir /app
COPY . /app
RUN mv /app/muxy_test_bak.go /app/muxy_test.go
WORKDIR /app
