# BWA API

Esta é uma API simples escrita em Go, configurada com Docker para facilitar o build e execução. O projeto utiliza um `Makefile` para automação das etapas de inicialização do módulo Go, construção da imagem Docker e execução do container.

## Requisitos

- [Go](https://golang.org/) (v1.16 ou superior)
- [Docker](https://www.docker.com/)

## Estrutura do Projeto

```bash

├── Dockerfile           # Arquivo Docker para construir a imagem da aplicação
├── go.mod               # Arquivo de dependências Go
├── main.go              # Código fonte principal da aplicação
├── Makefile             # Arquivo Make para automação do build e execução

```
## Instalação

Clone o repositório:

```bash
git clone https://github.com/bl4ckw1d0w/bwaapi.git
cd bwaapi
```

### (Opcional) Modifique o nome do módulo no Makefile, se necessário:

O nome padrão do módulo é `github.com/bl4ckw1d0w/bwaapi`. Se quiser mudar, basta ajustar o valor da variável `MODULE` no Makefile ou passar o valor como argumento:

```bash
make MODULE=github.com/usuario/novoprojeto
```

## Uso

### Construção e Execução com Docker

Para construir a aplicação Go e a imagem Docker, execute:

```bash
make build
```

Para rodar o container Docker:

```bash
make up
```

Isso iniciará a API e ela estará disponível em `http://localhost:8080`.


>Antes de iniciar o Kibana, certifique-se de que o usuário `kibana_user` possui as permissões necessárias. Você pode verificar e configurar as permissões do usuário com os seguintes passos:  

## Verificação de Permissões do Kibana

1. **Acesse o Elasticsearch:**
   Execute o seguinte comando para verificar as permissões do usuário `kibana_user`:
   ```bash
   curl -u elastic:pass -X GET "http://localhost:9200/_security/user/kibana_user"
   ```

2. **Verifique o Papel Atribuído:**
   O usuário deve ter um papel que contém as permissões necessárias, como `custom_kibana_admin`. Para verificar o papel, execute:
   ```bash
   curl -u elastic:pass -X GET "http://localhost:9200/_security/role/custom_kibana_admin"
   ```

3. **Atualize as Permissões, se Necessário:**
   Caso as permissões não estejam adequadas, você pode atualizar o papel com o seguinte comando:
   ```bash
   curl -u elastic:pass -X PUT "http://localhost:9200/_security/role/custom_kibana_admin" -H 'Content-Type: application/json' -d '{
     "cluster": ["all"],
     "indices": [
       {
         "names": [".kibana*", "*"],
         "privileges": ["all"]
       }
     ]
   }'
   ```

4. **Reinicie o Kibana:**
   Após garantir que as permissões estão corretas, reinicie o Kibana para aplicar as configurações.
```

Dessa forma, fica claro como verificar e corrigir as permissões necessárias para o funcionamento do Kibana. Se precisar de mais ajustes, é só avisar!

### Criando Usuários

Para criar um usuário, envie uma requisição `POST` para o endpoint `/create-user`:

```bash
curl -X POST http://localhost:8080/create-user \
     -H "Content-Type: application/json" \
     -d '{"username": "seu_usuario", "password": "mypassword"}'
```

Isso criará um novo usuário no banco de dados, com a senha criptografada.

## Contribuição

Sinta-se à vontade para abrir issues ou pull requests para contribuir com o projeto.