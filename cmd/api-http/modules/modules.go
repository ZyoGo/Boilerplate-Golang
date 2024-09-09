package modules

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZyoGo/default-ddd-http/config"
	"github.com/ZyoGo/default-ddd-http/pkg/bcrypt"
	"github.com/ZyoGo/default-ddd-http/pkg/database"
	"github.com/ZyoGo/default-ddd-http/pkg/logger"
	"github.com/ZyoGo/default-ddd-http/pkg/ulid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	userCore "github.com/ZyoGo/default-ddd-http/internal/user/core"
	userRouter "github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/http"
	userHttpV1 "github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/http/v1"
	userRepo "github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/repository/postgresql"
	userService "github.com/ZyoGo/default-ddd-http/internal/user/service"

	mwGin "github.com/ZyoGo/default-ddd-http/pkg/gin"
)

type userHandler struct {
	userHandlerV1 *userHttpV1.Handler
}

type HTTPServer struct {
	userHandler
	cfg    *config.AppConfig
	logger zerolog.Logger

	server *http.Server
	engine *gin.Engine
}

func New() (h *HTTPServer, err error) {
	h = &HTTPServer{}

	if err := h.initConfig(); err != nil {
		return nil, err
	}

	if err := h.initLogging(); err != nil {
		return nil, err
	}

	if err := h.initModules(); err != nil {
		return nil, err
	}

	if err := h.initHTTPServer(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *HTTPServer) initConfig() error {
	h.cfg = config.GetConfig()
	return nil
}

func (h *HTTPServer) initLogging() error {
	h.logger = logger.Get()
	return nil
}

func (h *HTTPServer) initModules() (err error) {
	dbConn := database.DatabaseConnection(h.cfg)
	ulidSvc := ulid.NewGenerator()
	hashSvc := bcrypt.New()

	var (
		userRepository userCore.Repository
	)
	{
		userRepository = userRepo.NewPostgreSQL(dbConn)
	}

	{
		userSvc, err := userService.New(
			userService.WithUserRepository(userRepository),
			userService.WithIDGenerator(ulidSvc),
			userService.WithHash(hashSvc),
		)
		if err != nil {
			return fmt.Errorf("%s for user service", err.Error())
		}

		userHTTP := userHttpV1.New(userSvc)
		h.userHandlerV1 = userHTTP
	}

	return
}

// shouldSkipLogging determines if logging should be skipped for a given request.
func shouldSkipLogging(path, method string) bool {
	return false
}

func (h *HTTPServer) initHTTPServer() (err error) {
	// init gin engine
	h.engine = gin.New()
	if h.cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Use middleware for recovery from panics and logging.
	h.engine.Use(gin.Recovery())
	h.engine.Use(mwGin.ZerologLoggerWithSkipper(h.logger, func(c *gin.Context) bool {
		return shouldSkipLogging(c.FullPath(), c.Request.Method)
	}))

	h.engine.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	userRouter.RegisterPath(h.engine, h.userHandlerV1)

	addSrv := fmt.Sprintf("%s:%d", h.cfg.App.Address, h.cfg.App.Port)
	h.server = &http.Server{
		Handler:      h.engine.Handler(),
		Addr:         addSrv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return
}

func (h *HTTPServer) Run() (err error) {
	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.logger.Fatal().Msg(err.Error())
		}
	}()

	h.Shutdown()
	return
}

func (h *HTTPServer) Shutdown() {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	h.logger.Warn().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Fatal().Str("Server Shutdown: ", err.Error())
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		h.logger.Warn().Msg("timeout of 5 seconds.")
	}
	h.logger.Warn().Msg("Server exiting")
}
