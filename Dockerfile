FROM golang:1.23-alpine AS builder
WORKDIR /gau_blog
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -tags '!dev' -o main .

FROM alpine:latest
WORKDIR /gau_blog
COPY --from=builder /gau_blog/main .
EXPOSE 8085
CMD ["./main"]