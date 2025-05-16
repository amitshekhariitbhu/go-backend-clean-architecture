FROM golang:1.19-alpine as builder
WORKDIR /
COPY go.mod .
RUN go mod download
COPY . .
RUN mkdir exec
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o exec/main cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /exec /app
USER 65532:65532
ENTRYPOINT ["/app/main"]
