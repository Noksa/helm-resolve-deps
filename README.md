# helm-resolve-deps

### A Helm plugin to properly resolve local chain dependencies in charts

If you have charts that have 'local' dependencies (charts that have repository starts with `file://`) and those dependencies also may have other local/external chain dependencies than you probably want to resolve all of those chain dependencies automatically and in a proper way. 

---
This plugin:
* Does it for you.
* Requires Helm3
* Can be used as replacement for default `helm dep up` command.

---

Why should you use it as replacement for `helm dep up`?

Because it has the `-u|--unpack-dependencies` flag which allows you to automatically unpack all dependent charts and see what manifests are inside them. 

Moreover in this case you are able to edit dependent charts right inside `charts/` directory. 

It can be helpful for debugging purposes.

Also the plugin has `-c|--clean` flag which allows you to remove charts, tmpcharts directories and Chart.lock file automatically.


#### And of course because the plugin does proper resolution of local chain dependencies.

---

## Installation

```
helm plugin install --version "main" https://github.com/Noksa/helm-resolve-deps.git
```

---

## Upgrade
```
helm plugin update resolve-deps
```

---

## Usage
Run this command to receive all available options:
```
helm resolve-deps -h
```
You can pass all flags from `helm dependency update` command to the plugin's command.

They  all will be substituted to `helm dependency update`.

---

## Custom flags
This plugin has its own flags. You can pass them in addition to `helm dep up` flags or without them.
```
-u[--unpack-dependencies] - untar/unpack dependent charts. They will be present as directories instead of .tgz archieves
-c[--clean]               - remove charts, tmpcharts directories and Chart.lock file in each chart before running the dependency update command
```

---

## A few examples:
```
helm resolve-deps . --skip-refresh
helm resolve-deps --clean
helm resolve-deps ~/charts/my-chart --skip-refresh --unpack-dependencies
helm resolve-deps ~/charts/my-chart --skip-refresh --unpack-dependencies --clean
```
