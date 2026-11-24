// Same names as var_naming_config_reset.go, but linted with the default configuration.

package fixtures

const VAR_NAMING_CONFIG_RESET_DEFAULT = 1 // MATCH /don't use ALL_CAPS in Go names; use CamelCase/

func varNamingConfigResetDefault() {
	someIdDefault := 1 // MATCH /var someIdDefault should be someIDDefault/
	someGrpcDefault := 2
	_, _ = someIdDefault, someGrpcDefault
}
