FROM golang:1.22
RUN go install github.com/air-verse/air@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOFLAGS=-buildvcs=false
CMD ["air"]
