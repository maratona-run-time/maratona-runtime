FROM golang:buster
EXPOSE 8080
WORKDIR /go/src/app
#RUN apk add python3
COPY go.sum go.mod ./
RUN go mod download
COPY executor/ executor/
CMD ["go", "run", "executor/main.go"]