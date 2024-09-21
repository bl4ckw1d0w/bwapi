IMAGE_NAME=bwapi
MODULE ?= github.com/bl4ckw1d0w/bwaapi


mod:
	@if [ ! -f go.mod ]; then \
		echo "Inicializando o módulo Go"; \
		go mod init $(MODULE); \
	fi
	@if [ -f go.mod ]; then \
		echo "Executando go mod tidy"; \
		go mod tidy; \
	fi

# Comando para construir a aplicação e a imagem Docker
build: mod
	@echo "Construindo a aplicação Go"
	go build -o main .
	@echo "Construindo a imagem Docker"
	docker build -t $(IMAGE_NAME) .
	

# Comando para rodar o container Docker
up: build
	@echo "Iniciando o ambiente com docker-compose"
	docker-compose up --build -d
	@sleep 2 # Esperar um pouco para o container iniciar
	@echo "A API está rodando em: http://localhost:8080"
