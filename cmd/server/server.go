package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/mattn/go-sqlite3"
	"github.com/viniciusidacruz/go-graphql/graph"
	"github.com/viniciusidacruz/go-graphql/internal/database"
)

const defaultPort = "8080"

func main() {
	// open the database connection
	db, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		log.Fatal("Failed to open database: %v", err)
	}
	defer db.Close()

	categoryDB := database.NewCategory(db)
	courseDB := database.NewCourse(db)
	// finish the database connection

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Implement database category in Resolver
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CategoryDB: categoryDB,
		CourseDB:   courseDB,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
