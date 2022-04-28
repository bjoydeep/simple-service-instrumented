FROM golang:latest 
RUN mkdir /app 
#RUN git clone  https://github.com/prometheus/client_golang.git /usr/local/go/src/github.com/prometheus/client_golang
#RUN go get github.com/prometheus/client_golang/prometheus/promhttp
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 
EXPOSE 8080
CMD ["/app/main"]
