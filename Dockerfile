FROM golang
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o playground
EXPOSE 9000
ENTRYPOINT ["./playground"]
