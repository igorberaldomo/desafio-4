FROM golang:1.22.1-alpine
WORKDIR /app
COPY /app .
RUN go mod tidy
RUN go build -o main.go .
CMD ["/app/main.go"]


