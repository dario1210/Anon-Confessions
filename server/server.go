package server

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	listenAddr string
	router     *gin.Engine
	db         *sql.DB
}

func NewAPIServer(listenAddr string, router *gin.Engine, db *sql.DB) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		router:     router,
		db:         db,
	}
}
