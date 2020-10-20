FROM golang:alpine
# create a working directory
WORKDIR /go/src/app
# add source code

COPY go.sum go.mod ./

RUN go mod download

COPY . .
# run main.go
#RUN go run main.go
#CMD ["bash"]

EXPOSE 8080

CMD ["go", "run", "comparator/server/server.go"]