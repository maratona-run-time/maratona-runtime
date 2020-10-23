FROM golang:alpine
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY . .
CMD ["go", "run", "main.go"]