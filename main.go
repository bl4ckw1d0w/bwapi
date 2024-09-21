package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
	"github.com/elastic/go-elasticsearch/v7"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB
var es *elasticsearch.Client

func main() {
	// Pegando as variáveis de ambiente
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	// String de conexão
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	// Conectando ao banco de dados
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializar o cliente Elasticsearch
	es, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_HOST"),
		},
		Username: os.Getenv("ELASTICSEARCH_USER"),
		Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("Erro ao criar cliente Elasticsearch: %v", err)
	}

	// Criar a tabela de usuários
	createUserTable()

	http.HandleFunc("/create-user", createUserHandler)

	log.Println("API rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createUserTable() {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Erro ao criar a tabela de usuários: %v", err)
	}
	log.Println("Tabela de usuários criada ou já existe.")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	// Validar dados do usuário
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Nome de usuário e senha são obrigatórios", http.StatusBadRequest)
		return
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao gerar o hash da senha", http.StatusInternalServerError)
		return
	}

	// Inserindo o usuário no banco de dados
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Erro ao inserir o usuário no banco de dados", http.StatusInternalServerError)
		return
	}

	// Registrar o log no Elasticsearch
	logToElasticsearch(fmt.Sprintf("Usuário %s criado com sucesso", user.Username))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Usuário criado com sucesso"))
}

func logToElasticsearch(message string) {
	doc := map[string]interface{}{
		"message": message,
		"time":    time.Now(),
	}
	docJSON, err := json.Marshal(doc)
	if err != nil {
		log.Printf("Erro ao marshaller o documento: %v", err)
		return
	}

	req := bytes.NewReader(docJSON)
	res, err := es.Index("logs", req)
	if err != nil {
		log.Printf("Erro ao enviar log para Elasticsearch: %v", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Erro ao indexar log: %s", res.String())
	}
}