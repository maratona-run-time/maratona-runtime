# Image to test the project on an alpine environment
FROM golang:alpine 
WORKDIR /go/src/app
RUN apk add g++ gcc go python3
COPY go.mod go.sum ./
RUN go mod download
# add source code
COPY . .
# run project tests
CMD ["./test.sh"]