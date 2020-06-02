FROM golang

WORKDIR /go/src/desafio-b2w
COPY . .

WORKDIR /go/src/desafio-b2w/cmd/planet-api
RUN go install -v

CMD ["planet-api"]