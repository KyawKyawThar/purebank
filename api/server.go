package api

import (
	"github.com/gin-gonic/gin"
	db "purebank/db/sqlc"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	router *gin.Engine
	//store connect to real db
	store  db.Store
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) (*Server, error) {

	server := &Server{
		store: store,
	}



	server.setUpRouter()
	return server, nil
}

func (s *Server) setUpRouter() {

	r := gin.Default()

	r.POST("/user", s.createUser)
	r.POST("/user/login", s.loginUser)
	s.router = r
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {

	return s.router.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
