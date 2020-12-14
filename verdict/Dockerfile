FROM golang:alpine
EXPOSE 8083
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY model/ model/
COPY utils/ utils/
COPY verdict/ verdict/
CMD ["go", "run", "verdict/main.go"]