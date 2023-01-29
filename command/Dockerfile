FROM golang:1.19.5-alpine
WORKDIR /app
COPY ../go.mod .
COPY ../go.sum .
RUN go mod download

COPY ./ .
RUN go build -o /command ./command/*go
EXPOSE 3000
CMD [ "/command" ]