FROM golang:1.17-alpine

WORKDIR /app



COPY . .

RUN go get gopkg.in/yaml.v2
RUN go build -o trustify main.go

ENTRYPOINT ["/app/trustify"]