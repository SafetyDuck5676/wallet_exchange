FROM golang:1.22

WORKDIR /app
ADD ./wallet .

#RUN go get github.com/SafetyDuck5676/grpc_duck
RUN go mod tidy
RUN go build -o main .

CMD ["./main"]