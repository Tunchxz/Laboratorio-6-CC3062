# Usar la imagen oficial de Golang
FROM golang:latest

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el código fuente
COPY . .

# Descargar e instalar las dependencias
RUN go get -d -v ./...

# Construir la aplicación de Go
RUN go build -o api .

# Exponer el puerto
EXPOSE 8080

# Correr el ejecutable
CMD [ "./api" ]