package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/renatomh/api-simplechat/db/sqlc"
)

// Server serves HTTP requests for the application
type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUser)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (serve *Server) Start(address string) error {
	return serve.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
