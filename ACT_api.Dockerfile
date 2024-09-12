FROM golang:1.21.1-bullseye

WORKDIR /application

COPY ./ .

CMD ["go", "run", "/application/main.go"]