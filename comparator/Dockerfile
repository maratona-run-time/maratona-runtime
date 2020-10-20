FROM golang:alpine

WORKDIR /go/src/app

COPY comparator.go .

COPY server/ ./server

RUN go get github.com/go-martini/martini
RUN go get github.com/martini-contrib/binding

CMD ["go", "run", "server/server.go"]