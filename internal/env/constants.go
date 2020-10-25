package env

const (
	// EnvFZFIgnore describes the environment variable to set to disable
	// interactive config selection when fzf is installed.
	EnvFZFIgnore = "CONFIG_REGISTRY_IGNORE_FZF"

	// EnvForceColor describes the environment variable to disable color usage
	// when printing current context in a list.
	EnvNoColor = `NO_COLOR`

	// EnvForceColor describes the "internal" environment variable to force
	// color usage to show current config in a list.
	EnvForceColor = `_CONFIG_REGISTRY_FORCE_COLOR`

	// EnvDebug describes the internal environment variable for more verbose logging.
	EnvDebug = `DEBUG`
)
