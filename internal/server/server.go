package server

import (
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"net/http"

	"github.com/ambardhesi/runnable/pkg/job"
	"github.com/ambardhesi/runnable/pkg/repository"
	"github.com/ambardhesi/runnable/pkg/runnable"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port   int
	LogDir string
	// TODO add cert dir/file attributes
}

type Server struct {
	config Config
	js     runnable.JobService
	lfs    runnable.LogFileService
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

func (s *Server) Start() {
	// Log HTTP server output to console.
	gin.DefaultWriter = io.MultiWriter(os.Stdout)

	router := gin.Default()

	// Wire up routes
	router.POST("/job", s.StartJob)
	router.GET("job/:id", s.GetJob)
	router.POST("/job/:id/stop", s.StopJob)
	router.GET("/job/:id/logs", s.GetJobLogs)

	s.monitorTerminationSignal()

	// Start server on port provided in config
	router.Run(":" + strconv.Itoa(s.config.Port))
}

func (s *Server) monitorTerminationSignal() {
	killChan := make(chan os.Signal)
	signal.Notify(killChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-killChan
		err := s.lfs.DeleteAllLogFiles()
		if err != nil {
			log.Printf("Failed to delete all log files before shutting down %v\n", err)
		}
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

	// TODO grab ownerID from cert CN, hardcoded for now
	cmd := strings.Split(request.Command, " ")
	jobID, err := s.js.Start("ownerID", cmd[0], cmd[1:]...)

	if err != nil {
		if runnable.ErrorCode(err) == runnable.EINVALID {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, StartJobResponse{
		JobID: jobID,
	})
}

func (s *Server) GetJob(ctx *gin.Context) {
	jobID := ctx.Param("id")
	// TODO grab ownerID from cert CN, hardcoded for now
	job, err := s.js.Get("ownerID", jobID)

	if err != nil {
		if runnable.ErrorCode(err) == runnable.ENOTFOUND {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, FromJob(job))
}

func (s *Server) StopJob(ctx *gin.Context) {
	jobID := ctx.Param("id")
	// TODO grab ownerID from cert CN, hardcoded for now
	err := s.js.Stop("ownerID", jobID)

	if err != nil {
		if runnable.ErrorCode(err) == runnable.ENOTFOUND {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if runnable.ErrorCode(err) == runnable.EINVALID {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.String(http.StatusOK, "")
}

func (s *Server) GetJobLogs(ctx *gin.Context) {
	jobID := ctx.Param("id")
	// TODO grab ownerID from cert CN, hardcoded for now
	logs, err := s.js.GetLogs("ownerID", jobID)

	if err != nil {
		if runnable.ErrorCode(err) == runnable.ENOTFOUND {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.String(http.StatusOK, *logs)
}
