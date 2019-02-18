FROM golang:alpine

WORKDIR /placechicken

ADD . /placechicken

RUN mkdir /resized

ENV STATIC ./static/images
ENV RESIZED /resized

RUN go build -o placechicken .

EXPOSE 8000

ENTRYPOINT ["./placechicken"]
