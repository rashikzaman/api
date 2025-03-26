FROM golang:1.23.2
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0

CMD ["make", "run"]