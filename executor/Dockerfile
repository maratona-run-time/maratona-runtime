FROM golang:alpine
EXPOSE 8080
WORKDIR /go/src/app
RUN apk add gcc g++ python3
COPY go.sum go.mod ./
RUN go mod download
COPY model/ model/
COPY executor/ executor/
CMD ["go", "run", "executor/main.go"]
