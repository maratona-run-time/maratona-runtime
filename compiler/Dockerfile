FROM golang:alpine
EXPOSE 8080
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY compiler/ compiler/
CMD ["go", "run", "main.go"]