package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/phathdt/service-context/component/fiberc"

	"mongo-streamer/shared/common"

	sctx "github.com/phathdt/service-context"

	"github.com/spf13/cobra"
)

const (
	serviceName = "mongo-streamer"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName(serviceName),
		sctx.WithComponent(fiberc.New(common.KeyCompFiber)),
	)
}

var rootCmd = &cobra.Command{
	Use:   serviceName,
	Short: fmt.Sprintf("start %s", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		sc := newServiceCtx()

		logger := sctx.GlobalLogger().GetLogger("service")

		time.Sleep(time.Second * 5)

		if err := sc.Load(); err != nil {
			logger.Fatal(err)
		}

		fiberComp := sc.MustGet(common.KeyCompFiber).(fiberc.FiberComponent)

		app := fiberComp.GetApp()

		if err := app.Listen(fmt.Sprintf(":%d", fiberComp.GetPort())); err != nil {
			logger.Fatal(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
