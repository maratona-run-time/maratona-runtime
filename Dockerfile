FROM golang:alpine
# create a working directory
WORKDIR /go/src/app
# add source code
COPY . .
# run main.go
#RUN go run main.go
#CMD ["bash"]

CMD ["go", "run", "comparator/server/server.go"]