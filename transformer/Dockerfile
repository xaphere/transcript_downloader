FROM golang:1

ENV PROJECT=extractor

COPY . /
WORKDIR /

RUN go build -mod=readonly -a -o /${PROJECT}

CMD ["/extractor"]