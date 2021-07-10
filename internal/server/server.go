package server

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"net/http"

	"github.com/ambardhesi/runnable/pkg/job"
	"github.com/ambardhesi/runnable/pkg/repository"
	"github.com/ambardhesi/runnable/pkg/runnable"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port           int
	LogDir         string
	CertFilePath   string
	KeyFilePath    string
	CaCertFilePath string
	TestMode       bool
}

type Server struct {
	config Config
	js     runnable.JobService
	lfs    runnable.LogFileService
	server http.Server
}

func NewServer(config Config) (*Server, error) {
	db := repository.NewInMemoryDB()
	lfs, err := repository.NewLocalFileSystem(config.LogDir)
	if err != nil {
		return nil, err
	}

	js := job.NewJobService(db, lfs)

	return &Server{
		config: config,
		js:     js,
		lfs:    lfs,
	}, nil
}

// extracts the client ID from the client cert CN and sets it as the key
func certMiddleware(ctx *gin.Context) {
	tls := ctx.Request.TLS
	if len(tls.PeerCertificates) == 0 {
		log.Printf("No cert found in request")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("No cert found")})
		return
	}

	clientCert := tls.PeerCertificates[0]
	ownerID := clientCert.Subject.CommonName

	ctx.Set("ownerID", ownerID)
	ctx.Next()
}

func (s *Server) Start() {
	if !s.config.TestMode {
		// Log HTTP server output to console.
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	} else {
		gin.DefaultWriter = ioutil.Discard
	}

	router := gin.Default()
	router.Use(certMiddleware)

	// Wire up routes
	router.POST("/job", s.StartJob)
	router.GET("job/:id", s.GetJob)
	router.POST("/job/:id/stop", s.StopJob)
	router.GET("/job/:id/logs", s.GetJobLogs)

	s.monitorTerminationSignal()

	tlsConfig, err := GetTLSConfig(s.config.CertFilePath, s.config.KeyFilePath, s.config.CaCertFilePath)
	if err != nil {
		log.Printf("Failed to get TLSConfig %v\n", tlsConfig)
		os.Exit(1)
	}

	// Start server on port provided in config
	//router.Run(":" + strconv.Itoa(s.config.Port))
	server := http.Server{
		Addr:      "localhost:" + strconv.Itoa(s.config.Port),
		Handler:   router,
		TLSConfig: tlsConfig,
	}
	s.server = server

	log.Fatal(server.ListenAndServeTLS("", ""))
}

func (s *Server) Stop() {
	err := s.lfs.DeleteAllLogFiles()
	if err != nil {
		log.Printf("Failed to delete all log files before shutting down %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}

func (s *Server) monitorTerminationSignal() {
	killChan := make(chan os.Signal)
	signal.Notify(killChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-killChan
		s.Stop()
		close(killChan)
		os.Exit(1)
	}()
}

func (s *Server) StartJob(ctx *gin.Context) {
	var request StartJobRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		log.Printf("failed to create start job request %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerID := ctx.GetString("ownerID")
	cmd := strings.Split(request.Command, " ")
	jobID, err := s.js.Start(ownerID, cmd[0], cmd[1:]...)

	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, StartJobResponse{
		JobID: jobID,
	})
}

func (s *Server) GetJob(ctx *gin.Context) {
	jobID := ctx.Param("id")
	ownerID := ctx.GetString("ownerID")

	job, err := s.js.Get(ownerID, jobID)

	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, FromJob(job))
}

func (s *Server) StopJob(ctx *gin.Context) {
	jobID := ctx.Param("id")
	ownerID := ctx.GetString("ownerID")

	err := s.js.Stop(ownerID, jobID)

	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.String(http.StatusOK, "")
}

func (s *Server) GetJobLogs(ctx *gin.Context) {
	jobID := ctx.Param("id")
	ownerID := ctx.GetString("ownerID")

	logs, err := s.js.GetLogs(ownerID, jobID)

	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.String(http.StatusOK, *logs)
}

func writeError(ctx *gin.Context, err error) {
	if runnable.ErrorCode(err) == runnable.ENOTFOUND {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if runnable.ErrorCode(err) == runnable.EINVALID {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	} else if runnable.ErrorCode(err) == runnable.EUNAUTHORIZED {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
