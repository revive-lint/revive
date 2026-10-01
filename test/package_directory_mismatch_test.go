package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageDirectoryMismatch(t *testing.T) {
	// Configure rule to ignore no directories so our tests in 'testdata/' can run
	config := &lint.RuleConfig{Arguments: lint.Arguments{map[string]any{"ignore-directories": []any{}}}}

	testRule(t, "package_directory_mismatch/good/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/bad/bad", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/maincmd/main", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/mixed/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/mixed/bad", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/go-good/good1", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/go-good/good2", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/go-good/bad", &rule.PackageDirectoryMismatchRule{}, config)

	// Test normalization cases
	testRule(t, "package_directory_mismatch/normalization/fo-ob_ar/foobar", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/foo_bar/foobar", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/foo.b_ar/foobar", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/go-foo_bar/foo_bar", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/go-foo_bar/foobar", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/go-foo-bar/foo_bar", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/go-foo-bar/foobar", &rule.PackageDirectoryMismatchRule{}, config)

	// Edge cases that are not logic, but we decided to ignore to have a simpler implementation
	testRule(t, "package_directory_mismatch/normalization/go-od/go_od", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/go-od/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/normalization/go-od/od", &rule.PackageDirectoryMismatchRule{}, config)

	// Test version directories (v1, v2, etc.)
	testRule(t, "package_directory_mismatch/api/v1/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/api/V1/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/api/V1/api_test", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/api/v1v/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/api/v2/v2", &rule.PackageDirectoryMismatchRule{}, config)

	// Test internal directory variations
	testRule(t, "package_directory_mismatch/internal/good/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/bad/bad", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/any", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/api/v1/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/api/v2/api", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/internal/v1/api", &rule.PackageDirectoryMismatchRule{}, config)

	// Test handling of testfiles
	testRule(t, "package_directory_mismatch/test/good_test", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/test/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/test/bad_test", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/test/bad", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/test/main_test", &rule.PackageDirectoryMismatchRule{}, config)

	// Test handling of root directories with go.mod
	testRule(t, "package_directory_mismatch/rootdir/good/client", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/rootdir/bad/client", &rule.PackageDirectoryMismatchRule{}, config)

	// Test handling of root directories with .git
	mkdirTempDotGit(t, "package_directory_mismatch/rootdir/withgit")
	testRule(t, "package_directory_mismatch/rootdir/withgit/client", &rule.PackageDirectoryMismatchRule{}, config)
}

func TestPackageDirectoryMismatchWithDefaultConfig(t *testing.T) {
	// Test with default configuration (should ignore testdata directories by default)
	testRule(t, "package_directory_mismatch/testdata/ignored", &rule.PackageDirectoryMismatchRule{})
}

func TestPackageDirectoryMismatchWithTestDirectories(t *testing.T) {
	// Test with new format that excludes specific test directories
	config := &lint.RuleConfig{Arguments: lint.Arguments{map[string]any{"ignore-directories": []any{"testinfo", "testutils"}}}}

	testRule(t, "package_directory_mismatch/testinfo/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/testutils/good", &rule.PackageDirectoryMismatchRule{}, config)

	// Ignoring a directory must not suppress similarly named sibling directories.
	prefixConfig := &lint.RuleConfig{Arguments: lint.Arguments{map[string]any{"ignore-directories": []any{"target"}}}}
	testRule(t, "package_directory_mismatch/targeted/wrong", &rule.PackageDirectoryMismatchRule{}, prefixConfig)
}

func TestPackageDirectoryMismatchWithMultipleDirectories(t *testing.T) {
	// Test with new named argument format that excludes both "testutils" and "testinfo"
	config := &lint.RuleConfig{Arguments: lint.Arguments{map[string]any{"ignoreDirectories": []any{"testutils", "testinfo"}}}}

	testRule(t, "package_directory_mismatch/testutils/good", &rule.PackageDirectoryMismatchRule{}, config)
	testRule(t, "package_directory_mismatch/testinfo/good", &rule.PackageDirectoryMismatchRule{}, config)
}
