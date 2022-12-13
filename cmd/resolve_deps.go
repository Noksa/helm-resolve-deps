package main

import (
	"github.com/noksa/helm-resolve-deps/internal/helpers"
	"github.com/noksa/helm-resolve-deps/internal/models"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	opts := models.HelmResolveDepsOptions{}
	pflag.BoolVar(&opts.SkipRefresh, "skip-refresh", false, "Skip fetching updates from helm repositories")
	pflag.BoolVarP(&opts.Clean, "clean", "c", true, "Remove charts, tmpcharts directories and Chart.lock file in each chart before running the dependency update command")
	pflag.BoolVar(&opts.Untar, "unpack-dependencies", false, "untar/unpack all (including external) dependent charts. They will be present as directories instead of .tgz archieves inside chartrs/ directory")
	pflag.BoolVarP(&opts.Untar, "untar", "u", false, "untar/unpack all (including external) dependent charts. They will be present as directories instead of .tgz archieves inside chartrs/ directory")
	pflag.StringSliceVar(&opts.SkipRefreshInCharts, "skip-refresh-in", []string{}, "skip fetching updates from helm repositories before resolving dependencies in specific charts (pass their names in the argument). Use ',' as delimiter if you want to specify more than one chart."+
		"")
	help := false
	pflag.BoolVarP(&help, "help", "h", false, "show usage")
	_ = pflag.CommandLine.MarkDeprecated("unpack-dependencies", "Use -u|--untar instead")
	cpuDefault := runtime.NumCPU() - 1
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		cpuDefault = cpuDefault - 1
	}
	if cpuDefault < 1 {
		cpuDefault = 1
	}
	pflag.IntVar(&opts.Threads, "threads", cpuDefault, "Number of CPUs to be used")
	pflag.Parse()
	if help {
		pflag.Usage()
		os.Exit(0)
	}
	chartPath := "."
	passedArgs := pflag.Args()
	if len(passedArgs) >= 1 {
		chartPath = passedArgs[0]
	}
	if len(passedArgs) > 1 {
		opts.Args = passedArgs[1:]
	}
	if strings.HasPrefix(chartPath, "~") {
		homeDir, err := os.UserHomeDir()
		helpers.Must(err)
		chartPath = chartPath[1:]
		chartPath = filepath.Join(homeDir, chartPath)
	}
	chartPath = filepath.Clean(chartPath)
	absPath, err := filepath.Abs(chartPath)
	helpers.Must(err)
	chartPath = absPath
	err = helpers.ResolveDependencies(chartPath, opts)
	helpers.Must(err)
}
