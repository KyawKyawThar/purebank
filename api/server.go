package api

import (
	"github.com/gin-gonic/gin"
	db "purebank/db/sqlc"
	"purebank/db/util"
	"purebank/pasetotoken"
	"purebank/worker"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config          util.Config
	maker           pasetotoken.Maker
	router          *gin.Engine
	taskdistributor worker.TaskDistributor
	//store connect to real db
	store db.Store
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store, taskdistributor worker.TaskDistributor) (*Server, error) {

	paaetoTokenMaker, err := pasetotoken.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, err
	}
	server := &Server{
		config:          config,
		maker:           paaetoTokenMaker,
		store:           store,
		taskdistributor: taskdistributor,
	}

	server.setUpRouter()
	return server, nil
}

func (s *Server) setUpRouter() {

	r := gin.Default()

	r.GET("/verify_email", s.verifyEmail)
	r.POST("/user", s.createUser)
	r.POST("/user/login", s.loginUser)
	r.POST("/tokens/renew_access", s.renewAccessToken)
	s.router = r
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {

	return s.router.Run(address)

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
