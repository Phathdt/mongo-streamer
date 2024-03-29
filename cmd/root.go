package cmd

import (
	"fmt"
	"github.com/phathdt/service-context/component/natspub"
	"mongo-streamer/plugins/mongoc"
	"mongo-streamer/streamer"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phathdt/service-context/component/fiberc"

	"mongo-streamer/shared/common"

	sctx "github.com/phathdt/service-context"

	"github.com/spf13/cobra"

	_ "github.com/joho/godotenv/autoload"
)

const (
	serviceName = "mongo-streamer"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName(serviceName),
		sctx.WithComponent(fiberc.New(common.KeyCompFiber)),
		sctx.WithComponent(mongoc.New(common.KeyMongo)),
		sctx.WithComponent(natspub.New(common.KeyNatsPub)),
	)
}

var rootCmd = &cobra.Command{
	Use:   serviceName,
	Short: fmt.Sprintf("start %s", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		sc := newServiceCtx()

		logger := sctx.GlobalLogger().GetLogger("service")

		time.Sleep(time.Second * 1)

		NewRouter(sc)

		if err := sc.Load(); err != nil {
			logger.Fatal(err)
		}

		go func() {
			s := streamer.New()

			if err := s.Run(sc); err != nil {
				logger.Error(err)

				os.Exit(1)
			}
		}()

		// gracefully shutdown
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		_ = sc.Stop()
		logger.Info("Server exited")
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
