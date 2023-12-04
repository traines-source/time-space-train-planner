FROM golang:alpine

ENV PORT=3000
ENV PROTOVER=25.1

WORKDIR /app

RUN apk add unzip && wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOVER}/protoc-${PROTOVER}-linux-x86_64.zip && unzip *.zip && mv bin/protoc /usr/bin  

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o /tstp ./cmd/main.go

EXPOSE 3000

CMD [ "/tstp" ]