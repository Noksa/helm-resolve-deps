package models

import "testing"

func TestHelmResolveDepsOptions_Defaults(t *testing.T) {
	opts := HelmResolveDepsOptions{}

	if opts.Clean {
		t.Error("expected Clean to be false by default")
	}
	if opts.Untar {
		t.Error("expected Untar to be false by default")
	}
	if opts.SkipRefresh {
		t.Error("expected SkipRefresh to be false by default")
	}
	if opts.Threads != 0 {
		t.Errorf("expected Threads to be 0 by default, got %d", opts.Threads)
	}
	if len(opts.SkipRefreshInCharts) != 0 {
		t.Errorf("expected SkipRefreshInCharts to be empty, got %v", opts.SkipRefreshInCharts)
	}
	if len(opts.Args) != 0 {
		t.Errorf("expected Args to be empty, got %v", opts.Args)
	}
}

func TestHelmResolveDepsOptions_WithValues(t *testing.T) {
	opts := HelmResolveDepsOptions{
		Clean:               true,
		Untar:               true,
		SkipRefresh:         true,
		SkipRefreshInCharts: []string{"chart1", "chart2"},
		Threads:             4,
		Args:                []string{"--debug", "--dry-run"},
	}

	if !opts.Clean {
		t.Error("expected Clean to be true")
	}
	if !opts.Untar {
		t.Error("expected Untar to be true")
	}
	if !opts.SkipRefresh {
		t.Error("expected SkipRefresh to be true")
	}
	if opts.Threads != 4 {
		t.Errorf("expected Threads to be 4, got %d", opts.Threads)
	}
	if len(opts.SkipRefreshInCharts) != 2 {
		t.Errorf("expected 2 charts in SkipRefreshInCharts, got %d", len(opts.SkipRefreshInCharts))
	}
	if len(opts.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(opts.Args))
	}
}
