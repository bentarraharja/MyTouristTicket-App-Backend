FROM golang:1.21.5-alpine

# membuat direktori folder
RUN mkdir /app

# set working direktory
WORKDIR /app

COPY ./ /app

RUN go mod tidy

# create executable
RUN go build -o server

# run executable file
CMD ["./server"]
