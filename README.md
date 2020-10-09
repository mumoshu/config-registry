# kubeconf

![Latest GitHub release](https://img.shields.io/github/release/mumoshu/kubeconf.svg)
![GitHub stars](https://img.shields.io/github/stars/mumoshu/kubeconf.svg?label=github%20stars)
[![CI](https://github.com/mumoshu/kubeconf/workflows/Test/badge.svg)](https://github.com/mumoshu/kubeconfactions?query=workflow%3A"Test")

`kubeconf` is a utility to manage and switch between kubeconfigs.

Use this to **avoid unintentional operation on your production clusters**, while easing the management of multiple clusters.

```
Usage:
  kubeconf [flags]
  kubeconf [command]

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

Use "kubeconf [command] --help" for more information about a command.
```

### Usage

```sh
$ kubecong init
✔ Config default created.

$ kubeconf cp . custom1
Copied config "default" to "custom1".

$ $EDITOR "$(kubeconf locate custom1)"

$ kubeconf mv custom1 prod
Renamed config "custom1" to "prod".

$ kubeconf ls
default
prod

$ kubeconf use prod
Switched to config "prod".

$ kubeconf use -
Switched to config "default".

$ kubeconf locate prod
$HOME/.kube/kubeconf/registry/prod
```

Switching to the production config by using `kubeconf use prod` is strongly NOT RECOMMEND.
That's because doing so would result in future you running a disruption operation on the production cluster without noticing the current config is production.

Instead, use `kubeconf locate prod`.

With that you can grab the kubeconfig path without switching,
so that you will never end up running unexpected operations in production:

```
# Avoid unintentional operation on prod by using `kubeconf locate`

$ KUBECONFIG=$(kubeconf locate prod) kubectl version
```

-----

## Installation

There two installation options:

- As kubectl plugins (macOS/Linux)
- Manual installation

### Kubectl Plugins (macOS and Linux)

You can install and use [Krew](https://github.com/kubernetes-sigs/krew/) kubectl
plugin manager to get `conf`.

```sh
kubectl krew install conf
```

After installing, the tools will be available as `kubectl conf`.

-----

### Manual installation

----

## Interactive mode

If you want `kubeconf` command to present you an interactive menu
with fuzzy searching, you just need to [install
`fzf`](https://github.com/junegunn/fzf) in your PATH.

If you have `fzf` installed, but want to opt out of using this feature, set the environment variable `KUBECONF_IGNORE_FZF=1`.

If you want to keep `fzf` interactive mode but need the default behavior of the command, you can do it using Unix composability:
```
kubeconf | cat
```


-----

### Customizing colors

If you like to customize the colors indicating the current config, set the environment variables `KUBECONF_CURRENT_FGCOLOR` and `KUBECONF_CURRENT_BGCOLOR` (refer color codes [here](https://linux.101hacks.com/ps1-examples/prompt-color-using-tput/)):

```
export KUBECONF_CURRENT_FGCOLOR=$(tput setaf 6) # blue text
export KUBECONF_CURRENT_BGCOLOR=$(tput setab 7) # white background
```

Colors in the output can be disabled by setting the
[`NO_COLOR`](http://no-color.org/) environment variable.

-----

# Acknowledgement

The initial version of `kubeconf` codebase has been roughly 50% derived from @ahmetb's awesome [kubectx](https://github.com/ahmetb/kubectx). You can see which source files are still kept without major changes today by seeing the license header comments in the source files. A big thanks to @ahmetb and the contributors for all the hard work, and sharing it as an opensource project!
