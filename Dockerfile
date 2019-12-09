FROM golang:latest

RUN mkdir /applications
ENV GO111MODULE=on
WORKDIR /applications

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ftpReader
EXPOSE 8181
ENTRYPOINT ["/applications/ftpReader"]