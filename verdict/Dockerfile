FROM golang:alpine AS build
EXPOSE 8083
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY model/ model/
COPY utils/ utils/
COPY verdict/ verdict/
RUN go build verdict/main.go

FROM alpine
WORKDIR /go/src/app
COPY --from=build /go/src/app/main .
CMD ["./main"]
