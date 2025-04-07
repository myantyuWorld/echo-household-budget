FROM golang:1.23 AS aws
ENV TZ Asia/Tokyo
WORKDIR /go/src/app
RUN go install github.com/rubenv/sql-migrate/...@latest
COPY ./ ./
RUN go mod download
WORKDIR /go/src/app//cmd
CMD ["go", "run", "main.go"]
