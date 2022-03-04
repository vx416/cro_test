package cmd

import (
	"context"
	"cro_test/internal/app"
	"cro_test/internal/delivery/http"
	"cro_test/pkg/config"
	"cro_test/pkg/logger"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "cro_test/docs"

	cobra "github.com/spf13/cobra"
	"go.uber.org/fx"
)

// ListenerCmd matching server failed
var ServerCmd = &cobra.Command{
	Run: runServer,
	Use: "server",
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func runServer(cmd *cobra.Command, args []string) {
	opts := []fx.Option{
		fx.Invoke(http.Routes), fx.Provide(config.StartEchoServer),
	}
	if strings.ToLower(os.Getenv("RUN_MIGRATIONS")) == "true" {
		opts = append(opts, fx.Invoke(config.RunMigrations))
	}
	app, err := app.NewServerApp("", opts...)
	if err != nil {
		log.Panicf("init config failed, err:%+v", err)
	}

	ctx := context.Background()
	l := logger.Get()
	err = app.Start(ctx)
	if err != nil {
		l.Errorf("start app failed, err:%+v", err)
		os.Exit(1)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	l.Info("shutdown process start")
	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		l.Errorf("shutdown process failed, err:%+v", err)
	}
}
