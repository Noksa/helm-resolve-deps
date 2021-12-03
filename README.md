# helm-resolve-deps

### A Helm plugin to resolve local and dependencies in a proper way

Is you have local charts that have dependencies as `file://` and they also use chain dependency model than you probably want to resolve local dependencies automatically. 
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
```
helm resolve-deps directory_with_a_chart
```
