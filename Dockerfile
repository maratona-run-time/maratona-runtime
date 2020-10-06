FROM golang:1.8.5-jessie
# create a working directory
WORKDIR /go/src/app
# add source code
COPY . .
# run main.go
RUN go run main.go
#CMD ["bash"]