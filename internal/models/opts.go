package models

type HelmResolveDepsOptions struct {
	Clean               bool
	Untar               bool
	SkipRefresh         bool
	SkipRefreshInCharts []string
	Threads             int
}
