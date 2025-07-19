package api

import (
	"fmt"

	db "github.com/akkahshh24/go-simple-bank-app/db/sqlc"
	"github.com/akkahshh24/go-simple-bank-app/token"
	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (s *Server) setupRouter() {
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)
	s.router.POST("/tokens/renew_access", s.renewAccessToken)

	authRoutes := s.router.Group("/").Use(authMiddleware(s.tokenMaker))
	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.listAccounts)

	authRoutes.POST("/transfers", s.createTransfer)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
