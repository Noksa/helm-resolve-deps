package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/noksa/helm-resolve-deps/internal/helpers"
	"github.com/noksa/helm-resolve-deps/internal/models"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	var opts models.HelmResolveDepsOptions

	cpuDefault := runtime.NumCPU() - 1
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		cpuDefault = cpuDefault - 1
	}
	if cpuDefault < 1 {
		cpuDefault = 1
	}

	rootCmd := &cobra.Command{
		Use:     "resolve-deps [PATH]",
		Short:   "Resolve local and remote dependencies in a proper, fast, concurrent way",
		Version: version,
		Long: `A Helm plugin that properly resolves local chain dependencies in charts.

Helm's built-in dependency update doesn't handle transitive local dependencies well.
This plugin recursively resolves all local dependencies in the correct order.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chartPath := "."
			if len(args) >= 1 {
				chartPath = args[0]
			}

			if strings.HasPrefix(chartPath, "~") {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				chartPath = filepath.Join(homeDir, chartPath[1:])
			}

			chartPath = filepath.Clean(chartPath)
			absPath, err := filepath.Abs(chartPath)
			if err != nil {
				return err
			}
			chartPath = absPath

			return helpers.ResolveDependencies(chartPath, opts)
		},
	}

	rootCmd.SetVersionTemplate(fmt.Sprintf("helm-resolve-deps version %s\n", version))

	rootCmd.Flags().BoolVar(&opts.SkipRefresh, "skip-refresh", false, "Skip fetching updates from helm repositories")
	rootCmd.Flags().BoolVarP(&opts.Clean, "clean", "c", false, "Remove charts, tmpcharts directories and Chart.lock file in each chart before running the dependency update command")
	rootCmd.Flags().BoolVarP(&opts.Untar, "untar", "u", false, "Untar/unpack all (including external) dependent charts. They will be present as directories instead of .tgz archives inside charts/ directory")
	rootCmd.Flags().StringSliceVar(&opts.SkipRefreshInCharts, "skip-refresh-in", []string{}, "Skip fetching updates from helm repositories before resolving dependencies in specific charts (comma-separated)")
	rootCmd.Flags().IntVar(&opts.Threads, "threads", cpuDefault, "Number of CPUs to be used")

	// Deprecated flag
	rootCmd.Flags().BoolVar(&opts.Untar, "unpack-dependencies", false, "Deprecated: use -u|--untar instead")
	_ = rootCmd.Flags().MarkDeprecated("unpack-dependencies", "use -u|--untar instead")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
