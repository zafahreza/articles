FROM golang:1.19.5-alpine
WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./ .
RUN go build -o /query ./query/*go
EXPOSE 3001
CMD [ "/query" ]