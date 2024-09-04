package modules

import (
	"context"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZyoGo/default-ddd-http/config"
	"github.com/ZyoGo/default-ddd-http/pkg/bcrypt"
	"github.com/ZyoGo/default-ddd-http/pkg/database"
	"github.com/ZyoGo/default-ddd-http/pkg/ulid"
	"github.com/gin-gonic/gin"

	userCore "github.com/ZyoGo/default-ddd-http/internal/user/core"
	userRouter "github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/http"
	userHttpV1 "github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/http/v1"
	userRepo "github.com/ZyoGo/default-ddd-http/internal/user/infrastructure/repository/postgresql"
	userService "github.com/ZyoGo/default-ddd-http/internal/user/service"
)

type httpHandler struct {
	userHandlerV1 *userHttpV1.Handler
}

type HTTPServer struct {
	httpHandler
	cfg *config.AppConfig

	server *http.Server
	engine *gin.Engine

	closers []func(context.Context) error
}

func New() (h *HTTPServer, err error) {
	h = &HTTPServer{}

	h.cfg = config.GetConfig()

	if err := h.registerModules(); err != nil {
		return nil, err
	}

	if err := h.buildServer(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *HTTPServer) registerModules() (err error) {
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

func (h *HTTPServer) buildServer() (err error) {
	banner := `
██████╗  ██████╗ ██╗██╗     ███████╗██████╗ ██████╗ ██╗      █████╗ ████████╗███████╗    
██╔══██╗██╔═══██╗██║██║     ██╔════╝██╔══██╗██╔══██╗██║     ██╔══██╗╚══██╔══╝██╔════╝    
██████╔╝██║   ██║██║██║     █████╗  ██████╔╝██████╔╝██║     ███████║   ██║   █████╗      
██╔══██╗██║   ██║██║██║     ██╔══╝  ██╔══██╗██╔═══╝ ██║     ██╔══██║   ██║   ██╔══╝      
██████╔╝╚██████╔╝██║███████╗███████╗██║  ██║██║     ███████╗██║  ██║   ██║   ███████╗    
╚═════╝  ╚═════╝ ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝    
                                                                                         
███╗   ██╗ ██████╗  ██████╗ ██████╗ ██╗     ███████╗    ██████╗ ███████╗██╗   ██╗        
████╗  ██║██╔═══██╗██╔═══██╗██╔══██╗██║     ██╔════╝    ██╔══██╗██╔════╝██║   ██║        
██╔██╗ ██║██║   ██║██║   ██║██║  ██║██║     █████╗      ██║  ██║█████╗  ██║   ██║        
██║╚██╗██║██║   ██║██║   ██║██║  ██║██║     ██╔══╝      ██║  ██║██╔══╝  ╚██╗ ██╔╝        
██║ ╚████║╚██████╔╝╚██████╔╝██████╔╝███████╗███████╗    ██████╔╝███████╗ ╚████╔╝         
╚═╝  ╚═══╝ ╚═════╝  ╚═════╝ ╚═════╝ ╚══════╝╚══════╝    ╚═════╝ ╚══════╝  ╚═══╝          
`
	stdLog.Println(banner)

	// init gin engine
	h.engine = gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	h.engine.Use(gin.Recovery())

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
			stdLog.Fatal(err)
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
	stdLog.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		stdLog.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		stdLog.Println("timeout of 5 seconds.")
	}
	stdLog.Println("Server exiting")
}
