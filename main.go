package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	PopulateConfig("credentials.json")
}

/*
	Steps (After getting the credentials from the json file)

1) Attempts to open an MySql connection with the info from the json
2) Attaches a Gin engine to the server
3) Sets up a CORS Middleware
4) Adds the custom defined validators and the available routes
5) Runs the server
*/
func main() {
	var err error

	mySqlInfo := fmt.Sprintf("%s:%s@/%s",
		Cred.User, Cred.Password, Cred.Dbname)
	Db, err = sql.Open("mysql", mySqlInfo)
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}(Db)

	err = Db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connection")

	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost:4200"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization",
			"accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router := server.Group("/api")
	addValidators()
	AddRoutes(router)

	err = server.Run("localhost:8080")
	if err != nil {
		return
	}
}
