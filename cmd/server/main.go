package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"ergo-tools-go/internal/api"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) { c.Writer.Header().Set("Access-Control-Allow-Origin", "*"); c.Next() })

	api.RegisterRoutes(r)

	srv := &http.Server{Addr: ":" + port, Handler: r, ReadHeaderTimeout: 10 * time.Second}
	log.Printf("→ ergo-tools-go listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
