package formatter_test

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mgechev/revive/formatter"
	"github.com/mgechev/revive/lint"
)

//revive:disable:line-length-limit
func TestFormatter(t *testing.T) {
	for name, td := range map[string]struct {
		formatter lint.Formatter
		failures  []lint.Failure
		want      string
		wantXML   bool
	}{
		"checkstyle": {
			formatter: &formatter.Checkstyle{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: `<?xml version='1.0' encoding='UTF-8'?>
<checkstyle version="5.0">
    <file name="err.go">
      <error line="33" column="4" message="replace fmt.Errorf by errors.New (confidence 0)" severity="error" source="revive/use-errors-new"/>
      <error line="38" column="4" message="replace fmt.Errorf by errors.New (confidence 0)" severity="error" source="revive/use-errors-new"/>
    </file>
    <file name="file.go">
      <error line="2" column="5" message="error var Exp should have name of the form ErrFoo (confidence 0)" severity="warning" source="revive/error-naming"/>
    </file>
</checkstyle>`,
			wantXML: true,
		},
		"checkstyle with XML special characters": {
			formatter: &formatter.Checkstyle{},
			failures: []lint.Failure{
				{
					Failure:  `error var Exp should have name of the form ErrFoo & not "Exp"`,
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: `./q"dir/a&b.go`,
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: `./q"dir/a&b.go`,
							Line:     2,
							Column:   10,
						},
					},
				},
			},
			want: `<?xml version='1.0' encoding='UTF-8'?>
<checkstyle version="5.0">
    <file name="./q&#34;dir/a&amp;b.go">
      <error line="2" column="5" message="error var Exp should have name of the form ErrFoo &amp; not &#34;Exp&#34; (confidence 0)" severity="warning" source="revive/error-naming"/>
    </file>
</checkstyle>`,
			wantXML: true,
		},
		"default": {
			formatter: &formatter.Default{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: `file.go:2:5: error var Exp should have name of the form ErrFoo
err.go:33:4: replace fmt.Errorf by errors.New
err.go:38:4: replace fmt.Errorf by errors.New`,
		},
		"friendly": {
			formatter: &formatter.Friendly{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "since Go 1.18 'interface{}' can be replaced by 'any'",
					RuleName: "use-any",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "any.go",
							Line:     4,
							Column:   6,
						},
						End: token.Position{
							Filename: "any.go",
							Line:     4,
							Column:   10,
						},
					},
				},
				{
					Failure:  "since Go 1.18 'interface{}' can be replaced by 'any'",
					RuleName: "use-any",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "any.go",
							Line:     8,
							Column:   16,
						},
						End: token.Position{
							Filename: "any.go",
							Line:     8,
							Column:   20,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: `  ⚠  https://revive.run/r#error-naming  error var Exp should have name of the form ErrFoo
  file.go:2:5

  ⚠  https://revive.run/r#use-any  since Go 1.18 'interface{}' can be replaced by 'any'
  any.go:4:6

  ⚠  https://revive.run/r#use-any  since Go 1.18 'interface{}' can be replaced by 'any'
  any.go:8:16

  ✘  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New
  err.go:33:4

  ✘  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New
  err.go:38:4

✘ 5 problems (2 errors, 3 warnings)

Errors:
  2  use-errors-new

Warnings:
  2  use-any
  1  error-naming

`,
		},
		"friendly no failures": {
			formatter: &formatter.Friendly{},
			failures:  []lint.Failure{},
			want:      "",
		},
		"friendly error": {
			formatter: &formatter.Friendly{},
			failures: []lint.Failure{
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
			},
			want: `  ✘  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New
  err.go:33:4

✘ 1 problem (1 error, 0 warnings)

Errors:
  1  use-errors-new

`,
		},
		"friendly warning": {
			formatter: &formatter.Friendly{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
			},
			want: `  ⚠  https://revive.run/r#error-naming  error var Exp should have name of the form ErrFoo
  file.go:2:5

⚠ 1 problem (0 errors, 1 warning)

Warnings:
  1  error-naming

`,
		},
		"json": {
			formatter: &formatter.JSON{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: "[" +
				`{"Severity":"warning","Failure":"error var Exp should have name of the form ErrFoo","RuleName":"error-naming","Category":"naming","Position":{"Start":{"Filename":"file.go","Offset":0,"Line":2,"Column":5},"End":{"Filename":"file.go","Offset":0,"Line":2,"Column":10}},"Confidence":0,"ReplacementLine":""}` +
				"," +
				`{"Severity":"error","Failure":"replace fmt.Errorf by errors.New","RuleName":"use-errors-new","Category":"errors","Position":{"Start":{"Filename":"err.go","Offset":0,"Line":33,"Column":4},"End":{"Filename":"err.go","Offset":0,"Line":33,"Column":8}},"Confidence":0,"ReplacementLine":""}` +
				"," +
				`{"Severity":"error","Failure":"replace fmt.Errorf by errors.New","RuleName":"use-errors-new","Category":"errors","Position":{"Start":{"Filename":"err.go","Offset":0,"Line":38,"Column":4},"End":{"Filename":"err.go","Offset":0,"Line":38,"Column":9}},"Confidence":0,"ReplacementLine":""}` +
				"]",
		},
		"ndjson": {
			formatter: &formatter.NDJSON{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: `{"Severity":"warning","Failure":"error var Exp should have name of the form ErrFoo","RuleName":"error-naming","Category":"naming","Position":{"Start":{"Filename":"file.go","Offset":0,"Line":2,"Column":5},"End":{"Filename":"file.go","Offset":0,"Line":2,"Column":10}},"Confidence":0,"ReplacementLine":""}` +
				"\n" +
				`{"Severity":"error","Failure":"replace fmt.Errorf by errors.New","RuleName":"use-errors-new","Category":"errors","Position":{"Start":{"Filename":"err.go","Offset":0,"Line":33,"Column":4},"End":{"Filename":"err.go","Offset":0,"Line":33,"Column":8}},"Confidence":0,"ReplacementLine":""}` +
				"\n" +
				`{"Severity":"error","Failure":"replace fmt.Errorf by errors.New","RuleName":"use-errors-new","Category":"errors","Position":{"Start":{"Filename":"err.go","Offset":0,"Line":38,"Column":4},"End":{"Filename":"err.go","Offset":0,"Line":38,"Column":9}},"Confidence":0,"ReplacementLine":""}` +
				"\n",
		},
		"plain": {
			formatter: &formatter.Plain{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: `file.go:2:5: error var Exp should have name of the form ErrFoo https://revive.run/r#error-naming` +
				"\n" +
				`err.go:33:4: replace fmt.Errorf by errors.New https://revive.run/r#use-errors-new` +
				"\n" +
				`err.go:38:4: replace fmt.Errorf by errors.New https://revive.run/r#use-errors-new` +
				"\n",
		},
		"sarif": {
			formatter: &formatter.Sarif{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: `{
  "runs": [
    {
      "results": [
        {
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "file.go"
                },
                "region": {
                  "startColumn": 5,
                  "startLine": 2
                }
              }
            }
          ],
          "message": {
            "text": "error var Exp should have name of the form ErrFoo"
          },
          "ruleId": "error-naming"
        },
        {
          "level": "error",
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "err.go"
                },
                "region": {
                  "startColumn": 4,
                  "startLine": 33
                }
              }
            }
          ],
          "message": {
            "text": "replace fmt.Errorf by errors.New"
          },
          "ruleId": "use-errors-new"
        },
        {
          "level": "error",
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "err.go"
                },
                "region": {
                  "startColumn": 4,
                  "startLine": 38
                }
              }
            }
          ],
          "message": {
            "text": "replace fmt.Errorf by errors.New"
          },
          "ruleId": "use-errors-new"
        }
      ],
      "tool": {
        "driver": {
          "informationUri": "https://revive.run",
          "name": "revive",
          "rules": [
            {
              "helpUri": "https://revive.run/r#use-errors-new",
              "id": "use-errors-new",
              "properties": {
                "severity": "error"
              }
            }
          ]
        }
      }
    }
  ],
  "version": "2.1.0"
}`,
		},
		"stylish": {
			formatter: &formatter.Stylish{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: `err.go
  (33, 4)  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New
  (38, 4)  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New

file.go
  (2, 5)  https://revive.run/r#error-naming  error var Exp should have name of the form ErrFoo


 ✖ 3 problems (2 errors) (1 warning)`,
		},
		"stylish no failures": {
			formatter: &formatter.Stylish{},
			failures:  []lint.Failure{},
			want:      "",
		},
		"stylish error": {
			formatter: &formatter.Stylish{},
			failures: []lint.Failure{
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: `err.go
  (38, 4)  https://revive.run/r#use-errors-new  replace fmt.Errorf by errors.New


 ✖ 1 problem (1 error) (0 warnings)`,
		},
		"stylish warning": {
			formatter: &formatter.Stylish{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
			},
			want: `file.go
  (2, 5)  https://revive.run/r#error-naming  error var Exp should have name of the form ErrFoo


 ✖ 1 problem (0 errors) (1 warning)`,
		},
		"unix": {
			formatter: &formatter.Unix{},
			failures: []lint.Failure{
				{
					Failure:  "error var Exp should have name of the form ErrFoo",
					RuleName: "error-naming",
					Category: lint.FailureCategoryNaming,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   5,
						},
						End: token.Position{
							Filename: "file.go",
							Line:     2,
							Column:   10,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     33,
							Column:   8,
						},
					},
				},
				{
					Failure:  "replace fmt.Errorf by errors.New",
					RuleName: "use-errors-new",
					Category: lint.FailureCategoryErrors,
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   4,
						},
						End: token.Position{
							Filename: "err.go",
							Line:     38,
							Column:   9,
						},
					},
				},
			},
			want: "file.go:2:5: [error-naming] error var Exp should have name of the form ErrFoo" +
				"\n" +
				"err.go:33:4: [use-errors-new] replace fmt.Errorf by errors.New" +
				"\n" +
				"err.go:38:4: [use-errors-new] replace fmt.Errorf by errors.New" +
				"\n",
		},
	} {
		//nolint:reassign // reassigning os.Stdout is necessary to test that formatters don't write to stdout.
		t.Run(name, func(t *testing.T) {
			t.Setenv("NO_COLOR", "true")
			realStdout := os.Stdout
			fakeStdout, err := os.Create(filepath.Join(t.TempDir(), "fakeStdout"))
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = fakeStdout
			t.Cleanup(func() { os.Stdout = realStdout })
			failures := make(chan lint.Failure, len(td.failures))
			for _, f := range td.failures {
				failures <- f
			}
			close(failures)
			output, err := td.formatter.Format(failures, lint.Config{
				Confidence: 0.8,
				Rules: lint.RulesConfig{
					"use-errors-new": lint.RuleConfig{
						Severity: lint.SeverityError,
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = realStdout
			err = fakeStdout.Close()
			if err != nil {
				t.Fatal(err)
			}
			stdout, err := os.ReadFile(fakeStdout.Name())
			if err != nil {
				t.Fatal(err)
			}
			if len(stdout) > 0 {
				t.Errorf("formatter wrote to stdout: %q", stdout)
			}
			if td.want != output {
				t.Errorf("got:\n%s\nwant:\n%s\n", output, td.want)
			}
			if td.wantXML {
				decoder := xml.NewDecoder(strings.NewReader(output))
				for {
					_, err := decoder.Token()
					if errors.Is(err, io.EOF) {
						break
					}
					if err != nil {
						t.Fatalf("output is not well-formed XML: %v\ngot:\n%s", err, output)
					}
				}
			}
		})
	}
}

type sarifTestLog struct {
	Runs []sarifTestRun `json:"runs"`
}

type sarifTestRun struct {
	Results []sarifTestResult `json:"results"`
	Tool    sarifTestTool     `json:"tool"`
}

type sarifTestResult struct {
	RuleID string `json:"ruleId"`
}

type sarifTestTool struct {
	Driver sarifTestDriver `json:"driver"`
}

type sarifTestDriver struct {
	Rules []sarifTestRule `json:"rules"`
}

type sarifTestRule struct {
	ID string `json:"id"`
}

func TestSARIFContainsAllConfiguredRules(t *testing.T) {
	sarifFormatter := &formatter.Sarif{}
	failures := make(chan lint.Failure, 1)
	failures <- lint.Failure{Failure: "example", RuleName: "second-rule"}
	close(failures)

	output, err := sarifFormatter.Format(failures, lint.Config{
		Rules: lint.RulesConfig{
			"first-rule":  {Severity: lint.SeverityWarning},
			"second-rule": {Severity: lint.SeverityError},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var log sarifTestLog
	if err := json.Unmarshal([]byte(output), &log); err != nil {
		t.Fatal(err)
	}
	if len(log.Runs) != 1 {
		t.Fatalf("got %d runs, want 1", len(log.Runs))
	}

	registeredRules := make(map[string]bool, len(log.Runs[0].Tool.Driver.Rules))
	for _, rule := range log.Runs[0].Tool.Driver.Rules {
		registeredRules[rule.ID] = true
	}
	if !registeredRules["first-rule"] || !registeredRules["second-rule"] {
		t.Fatalf("registered rules = %v, want first-rule and second-rule", registeredRules)
	}
	for _, result := range log.Runs[0].Results {
		if !registeredRules[result.RuleID] {
			t.Errorf("result references unregistered rule %q", result.RuleID)
		}
	}
}
