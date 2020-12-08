FROM golang:alpine
EXPOSE 8080
WORKDIR /go/src/app
RUN apk add g++
RUN apk add gcc
COPY go.sum go.mod ./
RUN go mod download
COPY utils/ utils/
COPY compiler/ compiler/
CMD ["go", "run", "compiler/main.go"]