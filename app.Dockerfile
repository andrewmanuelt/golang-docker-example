FROM golang:alpine3.15

WORKDIR /app/src 

COPY . .

ENV GO111MODULE=on

RUN go get 

RUN go build -o /golang-app .

CMD [ "/golang-app"]