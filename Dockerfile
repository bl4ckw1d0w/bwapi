# Usar uma imagem base do Go
FROM golang:1.21-alpine

WORKDIR /bwapi

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
