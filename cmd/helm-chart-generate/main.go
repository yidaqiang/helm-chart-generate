package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	helmchartgenerate "github.com/yidaqiang/helm-chart-generate"
	"github.com/yidaqiang/helm-chart-generate/pkg/helm"
	"io/fs"
	"strconv"

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
		ChartVersion       string
		appVersion         string
		chartName          string
		config             string
		username           string
		password           string
		checkHelmVersion   bool
		insecureSkipVerify bool
		dependencyUpdate   bool
		out                io.Writer
		timeout            int64
	}
	confg struct {
		CurrentContext string             `json:"current-context"`
		Contexts       map[string]context `json:"contexts"`
	}
	context struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	}
)

var (
	v2settings  v2environment.EnvSettings
	settings    = cli.New()
	globalUsage = `Helm plugin to generate helm chart from templates

Examples

  $ helm chart-gen <chart name>                      # generate helm chart from default templates configuration
  $ helm chart-gen <chart name> -c config.yml          # generate helm chart from custom templates configuration
`
)

func newChartGenCmd(args []string) *cobra.Command {
	g := &genCmd{}
	cmd := &cobra.Command{
		Use:          "helm chart-gen",
		Short:        "Helm plugin to generate helm chart from templates",
		Long:         globalUsage,
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If the --check-helm-version flag is provided, short circuit
			if g.checkHelmVersion {
				fmt.Println(helm.HelmMajorVersionCurrent())
				return nil
			}
			g.out = cmd.OutOrStdout()

			if len(args) != 2 {
				return errors.New("this command needs 1 arguments: name of chart")
			}
			g.chartName = args[1]
			fmt.Println(g.chartName)
			return g.gen()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&g.config, "config", "c", "", "Specify the configuration file for the generate helm chart")

	f.StringVarP(&g.ChartVersion, "version", "v", "", "Override chart version pre-gen")
	f.StringVarP(&g.appVersion, "app-version", "a", "", "Override app version pre-gen")
	f.StringVarP(&g.username, "username", "u", "", "Override HTTP basic auth username [$HELM_REPO_USERNAME]")
	f.StringVarP(&g.password, "password", "p", "", "Override HTTP basic auth password [$HELM_REPO_PASSWORD]")
	f.BoolVarP(&g.insecureSkipVerify, "insecure", "", false, "Connect to server with an insecure way by skipping certificate verification [$HELM_REPO_INSECURE]")
	f.BoolVarP(&g.dependencyUpdate, "dependency-update", "d", false, `update dependencies from "requirements.yaml" to dir "charts/" before packaging`)
	f.BoolVarP(&g.checkHelmVersion, "check-helm-version", "", false, `outputs either "2" or "3" indicating the current Helm major version`)
	f.Int64VarP(&g.timeout, "timeout", "t", 30, "The duration (in seconds) Helm will wait to get response from chartmuseum")
	f.BoolP("help", "h", false, "")

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

func (g *genCmd) setFieldsFromEnv() {
	if v, ok := os.LookupEnv("HELM_REPO_USERNAME"); ok && g.username == "" {
		g.username = v
	}
	if v, ok := os.LookupEnv("HELM_REPO_PASSWORD"); ok && g.password == "" {
		g.password = v
	}
	if v, ok := os.LookupEnv("HELM_REPO_INSECURE"); ok {
		g.insecureSkipVerify, _ = strconv.ParseBool(v)
	}

}

func (g *genCmd) gen() error {
	templatesFS := helmchartgenerate.GetTemplatesFS()
	debugFS(templatesFS)
	return nil
}

func debugFS(fsys fs.FS) {
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("Found: %s (Dir: %v)\n", path, d.IsDir())
		return nil
	})
}

func main() {
	cmd := newChartGenCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
