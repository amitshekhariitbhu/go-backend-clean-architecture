FROM golang:1.19-alpine as builder
WORKDIR /
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main main.go

FROM gcr.io/distroless/static:nonroot
RUN mkdir /app
WORKDIR /app
COPY --from=builder /main .
USER 65532:65532
ENTRYPOINT ["/app/main"]