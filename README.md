# config-registry

![Latest GitHub release](https://img.shields.io/github/release/mumoshu/config-registry.svg)
![GitHub stars](https://img.shields.io/github/stars/mumoshu/config-registry.svg?label=github%20stars)
[![CI](https://github.com/mumoshu/config-registry/workflows/Test/badge.svg)](https://github.com/mumoshu/config-registry/actions?query=workflow%3A"Test")

`config-registry` is a utility to manage and switch between kubeconfigs.

Use this to **avoid unintentional operation on your production clusters**, while easing the management of multiple clusters.

```
Usage:
  config-registry [flags]
  config-registry [command]

Available Commands:
  cp          copy config OLD to NEW
  current     show the current config name
  help        Help about any command
  import      import existing kubeconfig at PATH as NAME
  init        initialize kubeconf
  locate      print the path to config NAME
  ls          list the configs
  mv          rename config OLD to NEW
  rm          delete config NAME
  use         switch to config NAME

Flags:
  -h, --help   help for ./kubeconf

Use "config-registry [command] --help" for more information about a command.
```

### Usage

```sh
$ config-registry init
✔ Config default created.

$ config-registry cp . custom1
Copied config "default" to "custom1".

$ $EDITOR "$(config-registry locate custom1)"

$ config-registry mv custom1 prod
Renamed config "custom1" to "prod".

$ config-registry ls
default
prod

$ config-registry use prod
Switched to config "prod".

$ config-registry use -
Switched to config "default".

$ config-registry locate prod
$HOME/.kube/kubeconf/registry/prod
```

Switching to the production config by using `config-registry use prod` is strongly NOT RECOMMEND.
That's because doing so would result in future you running a disruption operation on the production cluster without noticing the current config is production.

Instead, use `config-registry locate prod`.

With that you can grab the kubeconfig path without switching,
so that you will never end up running unexpected operations in production:

```
# Avoid unintentional operation on prod by using `config-registry locate`

$ KUBECONFIG=$(config-registry locate prod) kubectl version
```

-----

## Installation

There two installation options:

- As kubectl plugins (macOS/Linux)
- Manual installation

### Kubectl Plugins (macOS and Linux)

You can install and use [Krew](https://github.com/kubernetes-sigs/krew/) kubectl
plugin manager to get `config-registry`.

```sh
kubectl krew install config-registry
```

After installing, the plugin will be available as `kubectl config-registry`.

-----

### Manual installation

----

## Interactive mode

If you want `config-registry` command to present you an interactive menu
with fuzzy searching, you just need to [install
`fzf`](https://github.com/junegunn/fzf) in your PATH.

If you have `fzf` installed, but want to opt out of using this feature, set the environment variable `CONFIG_REGISTRY_IGNORE_FZF=1`.

If you want to keep `fzf` interactive mode but need the default behavior of the command, you can do it using Unix composability:
```
config-registry | cat
```


-----

### Customizing colors

If you like to customize the colors indicating the current config, set the environment variables `CONFIG_REGISTRY_CURRENT_FGCOLOR` and `CONFIG_REGISTRY_CURRENT_BGCOLOR` (refer color codes [here](https://linux.101hacks.com/ps1-examples/prompt-color-using-tput/)):

```
export CONFIG_REGISTRY_CURRENT_FGCOLOR=$(tput setaf 6) # blue text
export CONFIG_REGISTRY_CURRENT_BGCOLOR=$(tput setab 7) # white background
```

Colors in the output can be disabled by setting the
[`NO_COLOR`](http://no-color.org/) environment variable.

-----

# Acknowledgement

The initial version of `config-registry` codebase has been roughly 50% derived from @ahmetb's awesome [kubectx](https://github.com/ahmetb/kubectx). You can see which source files are still kept without major changes today by seeing the license header comments in the source files. A big thanks to @ahmetb and the contributors for all the hard work, and sharing it as an opensource project!
