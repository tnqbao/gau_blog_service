FROM golang:1.23-alpine AS builder
WORKDIR /gau_blog
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env .  
RUN go build -o main .

FROM alpine:latest
WORKDIR /gau_user
COPY --from=builder /gau_blog/main .
COPY .env .  
EXPOSE 8085
CMD ["./main"]