# helm-resolve-deps

### A Helm plugin to resolve local dependencies in a proper way

If you have local charts that have dependencies as `file://` and they also have other local/external chain dependencies than you probably want to resolve all local chain dependencies automatically. 

This plugin does it for you.

## Installation

```
helm plugin install --version "main" https://github.com/Noksa/helm-resolve-deps.git
```

## Upgrade
```
helm plugin update resolve-deps
```


## Usage
Run this command to receive all available options:
```
helm resolve-deps -h
```
You can pass all flags from `helm dependency update` command to the plugin's command.

They  all will be substituted to `helm dependency update`.

Examples:
```
helm resolve-deps . --skip-refresh
helm resolve-deps 
helm resolve-deps ~/charts/my-chart --skip-refresh --unpack-dependencies
```

## Custom flags
This plugin has its own flags. You can pass them in addition to `helm dep up` flags or without them.
```
-u[--unpack-dependencies] - untar/unpack dependent charts. They will be present as directories instead of .tgz archieves
```

