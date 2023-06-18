FROM golang:1.20.3-alpine3.17 AS dev

ENV APP_VERSION ${APP_VERSION:-dev}
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY internal internal
COPY main.go main.go

ENTRYPOINT [ "go", "run", "main.go" ]

CMD [ "serve" ]

FROM dev AS builder

RUN go build -o bin/kube-apiserver-proxy main.go

ENTRYPOINT []

FROM scratch AS bin

COPY --from=builder /app/bin/kube-apiserver-proxy /kube-apiserver-proxy

ENTRYPOINT [ "/kube-apiserver-proxy" ]

CMD [ "serve" ]
