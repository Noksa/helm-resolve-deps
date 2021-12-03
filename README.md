# helm-resolve-deps

### A Helm plugin to resolve local and remote dependencies in a proper way

If you have local charts that have dependencies as `file://` and they also have other local/external chain dependencies than you probably want to resolve all local chain dependencies automatically. 

This plugin does it for you.

**Note** that this plugin also does `untar` so all subcharts will be unpacked as directories.
If it is not convenient let me know and I'll add a flag to enable/disable this feature.

## Installation

```
helm plugin install --version "main" https://github.com/Noksa/helm-resolve-deps.git
```

## Upgrade
```
helm plugin update resolve-deps
```


## Usage
```
helm resolve-deps directory_with_a_chart
```
