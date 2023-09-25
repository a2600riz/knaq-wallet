FROM golang:latest AS wallet
WORKDIR /app
COPY ./ ./
RUN go get -d -v ./...
RUN go build -o wallet

FROM node:latest
COPY --from=wallet /app/wallet /usr/local/bin
RUN npm install -g near-cli && apt-get update && apt-get install -y vim
EXPOSE 8080
CMD wallet