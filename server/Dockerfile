FROM golang:1.24.3-alpine

RUN apk add --no-cache git tzdata
ENV TZ=Asia/Jakarta

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./main.go

CMD ["go", "run", "./main.go"]