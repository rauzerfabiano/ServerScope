# Use a imagem oficial do Go
FROM golang:1.17

# Configura o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos locais para o container
COPY . .

# Instala as dependências
RUN go mod download

# Compila o programa
RUN go build -o servercope .

# Comando para executar o aplicativo
CMD ["/app/servercope"]
