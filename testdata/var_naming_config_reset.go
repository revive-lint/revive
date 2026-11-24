// Configured with allowlist, blocklist, skipInitialismNameChecks and upperCaseConst,
// so none of these names fail.

package fixtures

const VAR_NAMING_CONFIG_RESET = 1

func varNamingConfigReset() {
	someId := 1
	someGrpc := 2
	_, _ = someId, someGrpc
}
