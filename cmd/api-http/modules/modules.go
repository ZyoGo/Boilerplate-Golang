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
	"github.com/ZyoGo/default-ddd-http/pkg/database"
	"github.com/gin-gonic/gin"

	userCore "github.com/ZyoGo/default-ddd-http/internal/user-v1/core"
	userRouter "github.com/ZyoGo/default-ddd-http/internal/user-v1/infrastructure/http"
	userHttpV1 "github.com/ZyoGo/default-ddd-http/internal/user-v1/infrastructure/http/v1"
	userRepo "github.com/ZyoGo/default-ddd-http/internal/user-v1/infrastructure/repository/postgresql"
	userService "github.com/ZyoGo/default-ddd-http/internal/user-v1/service"

	_ "net/http/pprof"
)

const (
	CodeSuccess = iota
	CodeBadConfig
)

func Run() int {
	m, err := registerModules()
	if err != nil {
		return CodeBadConfig
	}

	return m.start()
}

type server struct {
	userHandlerV1 *userHttpV1.Handler
}

func registerModules() (*server, error) {
	s := &server{}

	// load config
	cfg := config.GetConfig()

	dbConn := database.DatabaseConnection(cfg)

	var (
		userRepository userCore.Repository
	)
	{
		userRepository = userRepo.NewPostgreSQL(dbConn)

	}

	{
		userSvc, err := userService.New(
			userService.WithUserRepository(userRepository),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize user service: %s", err.Error())
		}

		userHTTP := userHttpV1.New(userSvc)
		s.userHandlerV1 = userHTTP
	}

	return s, nil
}

func (s *server) start() int {
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

	// load config
	cfg := config.GetConfig()

	// init gin engine
	router := gin.New()

	addSrv := fmt.Sprintf("%s:%d", cfg.App.Address, cfg.App.Port)
	srv := http.Server{
		Handler:      router,
		Addr:         addSrv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	userRouter.RegisterPath(router, s.userHandlerV1)

	go func() {
		stdLog.Printf("Starting the server on %s", addSrv)
		if err := srv.ListenAndServe(); err != nil {
			stdLog.Fatal(err)
		}
	}()

	// Implement graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	stdLog.Println("Shutting down the server...")

	// Set a timeout for shutdown (for example, 5 seconds).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		stdLog.Fatalf("Server shutdown error: %v", err)
	}
	stdLog.Println("Server gracefully stopped")

	return CodeSuccess
}
