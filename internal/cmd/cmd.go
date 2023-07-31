package cmd

import (
	"ddns-server/internal/dao"
	"ddns-server/internal/logic"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	_ "ddns-server/internal/service"
)

type Cmd struct {
	*cobra.Command

	config string
}

func NewCmd() *Cmd {
	c := &Cmd{
		Command: &cobra.Command{
			Use:  filepath.Base(os.Args[0]),
			Args: cobra.OnlyValidArgs,
		},
	}
	c.Flags().StringVarP(&c.config, "config", "c", "config.yaml", "config file path")
	c.Command.Run = c.Run
	return c
}

func (c *Cmd) Run(cmd *cobra.Command, args []string) {
	dao.InitCofnigPath(c.config)

	if err := logic.Dns.ReloadRecord(); err != nil {
		c.PrintErrln(err)
		os.Exit(1)
	}

	if err := logic.Dns.ReloadService(); err != nil {
		c.PrintErrln(err)
		os.Exit(1)
	}

	if err := logic.Web.ReloadService(); err != nil {
		logic.Dns.CloseService()
		c.PrintErrln(err)
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logic.Dns.CloseService()
	logic.Web.CloseService()
}
