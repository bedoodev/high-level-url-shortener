FROM golang:1.24

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go install github.com/air-verse/air@latest
CMD ["air"]