package main

import (
    "database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, MDF!")
}

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
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
    log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
    }

    err = db.Ping()
    if err != nil {
    log.Fatalf("Erro ao pingar o banco de dados: %v", err)
    }
    log.Println("Conectado ao banco de dados com sucesso!")

    http.HandleFunc("/", helloHandler)
    http.ListenAndServe(":8080", nil)
}