ARG user
ARG token
FROM golang:1.23.1
WORKDIR /app
ENV GOPRIVATE=github.com/mayye4ka
RUN git config --global url."https://$user:$token@github.com".insteadOf "https://github.com"
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /pinder
EXPOSE 8080
CMD ["/pinder"]