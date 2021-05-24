package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Aswin/TodoList/graph"
	"github.com/Aswin/TodoList/graph/generated"
	database "github.com/Aswin/TodoList/internal/pkg/db/migrations/mysql"
	"github.com/gorilla/mux"
)

const defaultPort = "10000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := mux.NewRouter()
	database.InitDB()
	database.Migrate()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver()}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
