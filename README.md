# helm-resolve-deps

### A Helm plugin to properly resolve local chain dependencies in charts

If you have charts that have 'local' dependencies (charts that have repository starts with `file://`) and those dependencies also may have other local/external chain dependencies than you probably want to resolve all of those chain dependencies automatically and in a proper way. 

Inspired by https://github.com/helm/helm/issues/2247

---
This plugin:
* Does it for you.
* Requires Helm3
* Can be used as replacement for default `helm dep up` command.

---

Why should you use it as replacement for `helm dep up`?

Because it has the `-u|--untar` flag which allows you to automatically unpack all dependent charts and see what manifests are inside them. 

Moreover in this case you are able to edit dependent charts right inside `charts/` directory. 

It can be helpful for debugging purposes.

Also the plugin has `-c|--clean` flag which allows you to remove charts, tmpcharts directories and Chart.lock file automatically.


And of course because the plugin does proper resolution of local chain dependencies.

---

## Installation

```shell
helm plugin install --version "main" https://github.com/Noksa/helm-resolve-deps.git
```

To install old (bash-style) version, use the following command:
```shell
helm plugin install --version "v1.0.0" https://github.com/Noksa/helm-resolve-deps.git
```

---

## Upgrade

The best way to do it - reinstall it

If you encounter a problem during installation, try to remove helm plugins cache first
```
(h plugin uninstall resolve-deps || true) && h plugin install --version "main" https://github.com/Noksa/helm-resolve-deps.git
```

---

## Usage
Run this command to receive all available options:
```shell
helm resolve-deps -h
```
You can pass all flags from `helm dependency update` command to the plugin's command.

They  all will be substituted to `helm dependency update`.

To do that, use `--` as end for flags parsing and pass arguments after it:
```shell
helm resolve-deps path_to_chart -- --kubeconfig myconfig
```

---

## Custom flags
This plugin has its own flags. You can pass them in addition to `helm dep up` flags or without them.
```shell
-u[--untar]                   - untar/unpack dependent charts. They will be present as directories instead of .tgz archieves. Useful for debugging purposes
-c[--clean]                   - remove charts, tmpcharts directories and Chart.lock file in each chart before running the dependency update command
--skip-refresh-in name1,name2 - skip fetching updates from helm repositories before running 'helm dep up' in specific charts (pass their names in the argument)
--skip-refresh                - skip fetching updated from helm repositories
```

---

## A few examples:
```shell
helm resolve-deps . --skip-refresh
# another way to pass --skip-refresh as 'helm dep up' flag directly:
helm resolve-deps . -- --skip-refresh
helm resolve-deps --clean
helm resolve-deps ~/charts/my-chart --skip-refresh --untar
helm resolve-deps ~/charts/my-chart --skip-refresh -u -c
helm resolve-deps --skip-refresh-in my-chart1,my-second-chart
```
