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

	userHttpHandler "github.com/ZyoGo/default-ddd-http/internal/user/handler/http"
	userRepo "github.com/ZyoGo/default-ddd-http/internal/user/modules/postgresql"
	userService "github.com/ZyoGo/default-ddd-http/internal/user/service"

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
	userHandler *userHttpHandler.Handler
}

func registerModules() (*server, error) {
	s := &server{}

	// load config
	cfg := config.GetConfig()

	dbConn := database.DatabaseConnection(cfg)

	var (
		userRepository userRepo.Repository
	)
	{
		userRepository = userRepo.NewPostgreSQL(dbConn)

	}

	{
		userSvc, err := userService.New(userRepository)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize user service: %s", err.Error())
		}

		userHTTP := userHttpHandler.New(userSvc)
		s.userHandler = userHTTP
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

	addSrv := fmt.Sprintf("%s:%d", cfg.App.Address, cfg.App.Port)
	srv := http.Server{
		Addr:         addSrv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	router := http.NewServeMux()
	srv.Handler = router

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("health ok"))
	})

	userHttpHandler.RegisterPath(router, s.userHandler)

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
