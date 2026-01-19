package helpers

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/noksa/helm-resolve-deps/internal/models"
	"github.com/patrickmn/go-cache"
	"github.com/xxjwxc/gowp/workpool"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var c = cache.New(time.Minute*5, time.Minute*4)

func getMutex(objectName string) *sync.Mutex {
	mutexInterface, found := c.Get(objectName)
	var realMutex *sync.Mutex
	if !found || mutexInterface == nil {
		realMutex = &sync.Mutex{}
		_ = c.Add(objectName, realMutex, time.Minute*5)
	} else {
		realMutex = mutexInterface.(*sync.Mutex)
	}
	return realMutex
}

func cleanJoin(paths ...string) string {
	return filepath.Clean(filepath.Join(paths...))
}

func LoadChartByPath(chartPath string) (chart *models.MiniHelmChart, err error) {
	chart = &models.MiniHelmChart{}
	b, err := os.ReadFile(cleanJoin(chartPath, "Chart.yaml"))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &chart)
	if err != nil {
		return nil, err
	}
	chart.Path = chartPath
	return chart, nil
}

func updateRepositories(shouldSkip bool) error {
	if shouldSkip {
		return nil
	}
	attempt := 1
	var mErr error
	t := time.Now()
	for attempt <= 3 {
		func() {
			shSession := sh.NewSession()
			b := &strings.Builder{}
			fmt.Printf("Fetching updates from all helm repositories, attempt #%v ...\n", attempt)
			shSession.Stderr = b
			shSession.Stdout = b
			shSession.ShowCMD = true
			shSession.Command("helm", "repo", "up")
			err := shSession.Run()
			if err == nil {
				mErr = nil
				return
			}
			mErr = multierr.Append(mErr, err)
			mErr = multierr.Append(mErr, fmt.Errorf("%v", b))
			attempt++
		}()
		if mErr == nil {
			break
		}
	}
	if mErr == nil {
		fmt.Printf("  * Updates have been fetched, took %v\n", time.Since(t).Round(time.Millisecond))
	}
	return mErr
}

var resolvedDeps []string

func resolveDeps(chart *models.MiniHelmChart, chartPath string, wp *workpool.WorkPool, opts models.HelmResolveDepsOptions) (mErr error) {
	fullName := fmt.Sprintf("%v-%v-%v", chart.Name, chart.Version, chart.Repository)
	chartMutex := getMutex(fullName)
	chartMutex.Lock()
	defer chartMutex.Unlock()
	if slices.Contains(resolvedDeps, fullName) {
		return nil
	}
	errMutex := sync.Mutex{}
	var args []string
	args = append(args, "dep", "up")
	if opts.Clean {
		_ = os.RemoveAll(cleanJoin(chartPath, "charts"))
		_ = os.RemoveAll(cleanJoin(chartPath, "tmpcharts"))
		_ = os.RemoveAll(cleanJoin(chartPath, "Chart.lock"))
	}
	if opts.SkipRefresh || slices.Contains(opts.SkipRefreshInCharts, chart.Name) {
		args = append(args, "--skip-refresh")
	}
	for _, additionalArg := range opts.Args {
		if !slices.Contains(args, additionalArg) {
			args = append(args, additionalArg)
		}
	}
	for _, dep := range chart.Dependencies {
		newOptions := models.HelmResolveDepsOptions{
			SkipRefresh: true,
			Clean:       opts.Clean,
			Threads:     opts.Threads,
			Untar:       opts.Untar,
		}
		dep := dep
		if after, ok := strings.CutPrefix(dep.Repository, "file://"); ok {
			dependantChartPath := filepath.Clean(fmt.Sprintf("%v/%v", chartPath, after))
			dependantChart, err := LoadChartByPath(dependantChartPath)
			if err != nil {
				mErr = multierr.Append(mErr, err)
				continue
			}
			if wp != nil {
				wp.Do(func() error {
					depErr := resolveDeps(dependantChart, dependantChartPath, nil, newOptions)
					errMutex.Lock()
					if depErr != nil {
						mErr = multierr.Append(mErr, fmt.Errorf("got an error while was resolving deps for %v chart", dependantChart.Name))
					}
					mErr = multierr.Append(mErr, depErr)
					errMutex.Unlock()
					return nil
				})
			} else {
				depErr := resolveDeps(dependantChart, dependantChartPath, nil, newOptions)
				errMutex.Lock()
				if depErr != nil {
					mErr = multierr.Append(mErr, fmt.Errorf("got an error while was resolving deps for %v chart", dependantChart.Name))
				}
				mErr = multierr.Append(mErr, depErr)
				errMutex.Unlock()
				continue
			}
		}
	}
	if wp != nil {
		_ = wp.Wait()
	}
	if mErr != nil {
		return mErr
	}
	if len(chart.Dependencies) == 0 {
		resolvedDeps = append(resolvedDeps, fullName)
		return nil
	}

	mErr = multierr.Append(mErr, func() error {
		shSession := sh.NewSession()
		shSession.ShowCMD = true
		b := &strings.Builder{}
		shSession.Stderr = b
		shSession.Stdout = b
		shSession.SetDir(chartPath)
		shSession.Command("helm", args)
		miniErr := shSession.Run()
		if miniErr != nil {
			return fmt.Errorf("%v, %v", miniErr.Error(), b.String())
		}
		return nil
	}())
	resolvedDeps = append(resolvedDeps, fullName)
	return mErr
}

func ResolveDependencies(chartPath string, opts models.HelmResolveDepsOptions) error {
	chart, err := LoadChartByPath(chartPath)
	if err != nil {
		return multierr.Append(err, fmt.Errorf("ensure that the chart directory (%v) exists", filepath.Dir(chartPath)))
	}
	err = updateRepositories(opts.SkipRefresh)
	if err != nil {
		return err
	}
	t := time.Now()
	wp := workpool.New(opts.Threads)
	if !opts.SkipRefresh {
		opts.SkipRefresh = true
	}
	fmt.Printf("Resolving dependencies in %v chart ...\n", chart.Name)
	err = resolveDeps(chart, chartPath, wp, opts)
	if err != nil {
		return err
	}
	if opts.Untar {
		err = filepath.WalkDir(cleanJoin(chartPath, "charts"), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".tgz") {
				return nil
			}
			shSession := sh.NewSession()
			shSession.Command("tar", "xzfm", path, "-C", cleanJoin(chartPath, "charts"))
			shSession.Stdout = io.Discard
			b := &strings.Builder{}
			shSession.Stderr = b
			err = shSession.Run()
			if err != nil {
				return err
			}
			err = os.RemoveAll(path)
			return err
		})
	}
	if err == nil {
		fmt.Printf("  * Dependencies have been resolved, took %v\n", time.Since(t).Round(time.Millisecond))
	}
	return err
}
