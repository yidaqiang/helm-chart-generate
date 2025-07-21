package main

import (
	"fmt"
	"github.com/spf13/cobra"
	//"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	//"helm.sh/helm/v3/pkg/downloader"
	//"helm.sh/helm/v3/pkg/getter"
	"io"
	//"k8s.io/helm/pkg/helm"
	v2environment "k8s.io/helm/pkg/helm/environment"
	"os"
)

type (
	genCmd struct {
		checkHelmVersion bool
		out              io.Writer
	}
	confg struct {
	}
	context struct {
	}
)

var (
	v2settings  v2environment.EnvSettings
	settings    = cli.New()
	globalUsage = `Helm plugin to generate helm chart from templates

Examples

  $ helm chart-gen <chart name>                      # generate helm chart from default templates configuration
  $ helm cm-push <chart name> -c config.yml          # generate helm chart from custom templates configuration
`
)

func newChartGenCmd(args []string) *cobra.Command {
	p := &genCmd{}
	cmd := &cobra.Command{
		Use:          "helm chart-gen <chart name>",
		Short:        "Helm plugin to generate helm chart from templates",
		Long:         globalUsage,
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If the --check-helm-version flag is provided, short circuit
			if p.checkHelmVersion {
				//fmt.Println(helm.HelmMajorVersionCurrent())
				return nil
			}
			p.out = cmd.OutOrStdout()

			return p.gen()
		},
	}

	f := cmd.Flags()

	err := f.Parse(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot parse flags: %s\n", err)
		// still return other flags , do not panic the following code
		return cmd
	}
	v2settings.AddFlags(f)
	v2settings.Init(f)

	return cmd
}

func (p *genCmd) gen() error {

	return nil
}

func main() {
	cmd := newChartGenCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
