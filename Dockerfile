#Base Images
FROM golang:1.21.4-alpine3.18

#Environtment Tambahan
#ENV
ENV GO_ENV=DOCKER DB_HOST=host.docker.internal

#Menentukan wworking directory di container
WORKDIR /app

#Copy project ke working directory
COPY . .

#jalankan perintah (instalasi, build, dll) di kontainer
#1. Install dependency
RUN go mod download

#2. Build Application
RUN go build -v -o /app/coffee-shop ./cmd/main.go

#Expose port
EXPOSE 6121

#daftarkan aplikasi
ENTRYPOINT ["/app/coffee-shop"]