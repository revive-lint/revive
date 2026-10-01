package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestStructTag(t *testing.T) {
	testRule(t, "struct_tag", &rule.StructTagRule{})
}

func TestStructTagWithUserOptions(t *testing.T) {
	testRule(t, "struct_tag_user_options", &rule.StructTagRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			"json,inline,outline",
			"bson,gnu",
			"url,myURLOption",
			"datastore,myDatastoreOption",
			"mapstructure,myMapstructureOption",
			"validate,displayName",
			"toml,unknown",
			"spanner,mySpannerOption",
			"codec,myCodecOption",
			"cbor,myCborOption",
		},
	})
}

func TestStructTagWithOmittedTags(t *testing.T) {
	testRule(t, "struct_tag_user_options_omit", &rule.StructTagRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			"!validate",
			"!toml",
			"json,inline,outline",
			"bson,gnu",
			"url,myURLOption",
			"datastore,myDatastoreOption",
			"mapstructure,myMapstructureOption",
			"spanner,mySpannerOption",
		},
	})
}

func TestStructTagAfterGo1_24(t *testing.T) {
	testRule(t, "go1.24/struct_tag", &rule.StructTagRule{})
}

func TestStructTagConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.StructTagRule{}
	testRule(t, "struct_tag_config_reset", r, &lint.RuleConfig{
		Arguments: lint.Arguments{"validate,myOption"},
	})
	testRule(t, "struct_tag_config_reset_ok", r)
}
