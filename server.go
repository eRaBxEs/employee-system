package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"employee-management-system/graph"
	"employee-management-system/pkg/environment"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/denisenkom/go-mssqldb" // Microsoft SQL Driver
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
)

const defaultPort = "8080"

func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("mssql", connStr)
	if err != nil {
		return nil, err
	}

	if db != nil {
		log.Println("Connected to Azure Sql Edge")
	}

	SelectVersion(db)

	if err := goose.SetDialect("mssql"); err != nil {
		return nil, fmt.Errorf("error in setting dialect: %v", err)
	}

	// Run migrations
	// if err := goose.Up(db, "./terminal/goose/migrations"); err != nil {
	// 	return nil, fmt.Errorf("error applying migrations: %v", err)
	// }

	return db, nil
}

// Gets and prints SQL Server version
func SelectVersion(db *sql.DB) {
	// Use background context
	ctx := context.Background()

	// Ping database to see if it's still alive.
	// Important for handling network issues and long queries.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}

	var result string

	// Run query and scan for result
	err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Printf("%s\n", result)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// get the environment
	env, err := environment.NewLoadFromFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Configure database connection
	// connStr := "server=localhost,57000;user id=SA;password='SUREcollection7!'database=CompanyDB"
	// connStr := "sqlserver://SA:SUREcollection7!@localhost:57000?database=CompanyDB"
	connStr := fmt.Sprintf("sqlserver://%s:%s@%s:%s?%s?sslmode=disable",
		env.Get("MSSQL_USER"),
		env.Get("MSSQL_PASSWORD"),
		env.Get("MSSQL_ADDRESS"),
		env.Get("MSSQL_PORT"),
		env.Get("MSSQL_DATABASE"))

	// Initialize the database and run migrations
	db, err := InitDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(corsMiddleware()) // Add this line to apply the CORS middleware

	// Set up GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/query")))
	r.POST("/query", gin.WrapH(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
