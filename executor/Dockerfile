FROM golang:alpine AS build
EXPOSE 8082
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY model/ model/
COPY utils/ utils/
COPY executor/ executor/
RUN go build executor/main.go

FROM alpine
WORKDIR /go/src/app
RUN apk add go gcc g++ python3
COPY --from=build /go/src/app/main .
CMD ["./main"]
