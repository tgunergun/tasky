# Building the binary of the App
FROM golang:1.19 AS build

ARG arch=amd64
WORKDIR /go/src/tasky
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH="$ARG" go build -o /go/src/tasky/tasky


FROM alpine:3.19.0 as release

WORKDIR /app
COPY --from=build  /go/src/tasky/tasky .
COPY --from=build  /go/src/tasky/assets ./assets
EXPOSE 8080
ENTRYPOINT ["/app/tasky"]


