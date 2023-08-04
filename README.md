# ServerScope

ServerScope é uma ferramenta de monitoramento em tempo real que exibe o uso de CPU, memória, disco e tráfego de rede de um servidor. É fácil de usar e pode ser executada localmente ou em um contêiner Docker.

## Requisitos

- Go 1.16 ou superior
- Docker (opcional para execução em contêiner)

## Instalação

### Usando Go

1. Clone o repositório:
   ```bash
   git clone https://github.com/rauzerfabiano/ServerScope.git

2. Entre no diretório do projeto:
   ```bash
   cd ServerScope

3. Instale as dependências:
   ```bash
   - go get github.com/gizak/termui/v3
   - go get github.com/shirou/gopsutil/v3/cpu
   - go get github.com/shirou/gopsutil/v3/mem
   - go get github.com/shirou/gopsutil/v3/disk
   - go get github.com/shirou/gopsutil/v3/net

4. Execute o projeto:
   ```bash
   go run main.go

### Usando Docker

1. Clone o repositório e entre no diretório:
   ```bash
   git clone https://github.com/rauzerfabiano/ServerScope.git
   cd ServerScope

2. Construa a imagem Docker:
   ```bash
   docker build -t serverscope .

3. Execute o container:
   ```bash
   docker run -it servercope
