package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/renatomh/api-simplechat/db/sqlc"
	"github.com/renatomh/api-simplechat/token"
	"github.com/renatomh/api-simplechat/util"
)

// Server serves HTTP requests for the application
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Adding routes to the router
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// Defining group of routes which require authentication
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.GET("/users", server.listUser)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (serve *Server) Start(address string) error {
	return serve.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
