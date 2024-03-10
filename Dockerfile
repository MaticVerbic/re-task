FROM golang:1.22
LABEL authors="Matic Verbic"

COPY . /app
WORKDIR /app

RUN go mod download

RUN mkdir /tmp/build
RUN cp ./config.yaml /tmp/build/config.yaml

RUN CGO_ENABLE=0 GOOS=linux go build -o /tmp/build/re-task cmd/main.go

CMD ["/tmp/build/re-task"]