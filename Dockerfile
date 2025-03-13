FROM golang:1.22.1-alpine as build
WORKDIR /app
COPY /app .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
CMD ["/app/main"]
EXPOSE 8080


FROM scratch
WORKDIR /app
COPY --from=build /app/main .
ENTRYPOINT [ "./main" ]