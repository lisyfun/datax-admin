FROM --platform=linux/amd64 golang:1.22.2-alpine

WORKDIR /app

COPY . .

# linux/amd64
RUN export GOOS=linux && export GOARCH=amd64 && go build -o datax-admin .

CMD ["./datax-admin"]
