package rule_test

import (
	"testing"

	"github.com/mgechev/revive/internal/rule"
)

func TestName(t *testing.T) {
	tests := []struct {
		name                     string
		allowlist                []string
		blocklist                []string
		skipInitialismNameChecks bool
		initialismsAsWords       bool
		want                     string
	}{
		{
			name: "foo_bar",
			want: "fooBar",
		},
		{
			name: "foo_bar_baz",
			want: "fooBarBaz",
		},
		{
			name: "Foo_bar",
			want: "FooBar",
		},
		{
			name: "foo_WiFi",
			want: "fooWiFi",
		},
		{
			name: "id",
			want: "id",
		},
		{
			name: "Id",
			want: "ID",
		},
		{
			name: "foo_id",
			want: "fooID",
		},
		{
			name: "fooId",
			want: "fooID",
		},
		{
			name: "fooUid",
			want: "fooUID",
		},
		{
			name: "idFoo",
			want: "idFoo",
		},
		{
			name: "uidFoo",
			want: "uidFoo",
		},
		{
			name: "midIdDle",
			want: "midIDDle",
		},
		{
			name: "APIProxy",
			want: "APIProxy",
		},
		{
			name: "ApiProxy",
			want: "APIProxy",
		},
		{
			name: "apiProxy",
			want: "apiProxy",
		},
		{
			name: "_Leading",
			want: "_Leading",
		},
		{
			name: "___Leading",
			want: "_Leading",
		},
		{
			name: "trailing_",
			want: "trailing",
		},
		{
			name: "trailing___",
			want: "trailing",
		},
		{
			name: "a_b",
			want: "aB",
		},
		{
			name: "a__b",
			want: "aB",
		},
		{
			name: "a___b",
			want: "aB",
		},
		{
			name: "Rpc1150",
			want: "RPC1150",
		},
		{
			name: "case3_1",
			want: "case3_1",
		},
		{
			name: "case3__1",
			want: "case3_1",
		},
		{
			name: "IEEE802_16bit",
			want: "IEEE802_16bit",
		},
		{
			name: "IEEE802_16Bit",
			want: "IEEE802_16Bit",
		},
		{
			name: "IDS",
			want: "IDs",
		},
		// Test skipInitialismChecks functionality
		{
			name:                     "getJson",
			skipInitialismNameChecks: true,
			want:                     "getJson",
		},
		{
			name:                     "userId",
			skipInitialismNameChecks: true,
			want:                     "userId",
		},
		{
			name:                     "myHttpClient",
			skipInitialismNameChecks: true,
			want:                     "myHttpClient",
		},
		// Test allowlist functionality
		{
			name:      "fooId",
			allowlist: []string{"ID"},
			want:      "fooId",
		},
		{
			name:      "fooApi",
			allowlist: []string{"API"},
			want:      "fooApi",
		},
		{
			name:      "fooHttp",
			allowlist: []string{"HTTP"},
			want:      "fooHttp",
		},
		// Test blocklist functionality
		{
			name:      "fooCustom",
			blocklist: []string{"CUSTOM"},
			want:      "fooCUSTOM",
		},
		{
			name:      "mySpecial",
			blocklist: []string{"SPECIAL"},
			want:      "mySPECIAL",
		},
		// Test combination of allowlist and blocklist
		{
			name:      "fooIdCustom",
			allowlist: []string{"ID"},
			blocklist: []string{"CUSTOM"},
			want:      "fooIdCUSTOM",
		},
		// Test combination of allowlist, blocklist and skipInitialismChecks
		{
			name:                     "fooIdCustomHttpJson",
			allowlist:                []string{"ID"},
			blocklist:                []string{"CUSTOM"},
			skipInitialismNameChecks: true,
			want:                     "fooIdCustomHttpJson",
		},
		// Test uppercase runs with initialismsAsWords=false
		{
			name: "HTTPMethod",
			want: "HTTPMethod",
		},
		{
			name: "HttpMethod",
			want: "HTTPMethod",
		},
		{
			name: "apiURL",
			want: "apiURL",
		},
		{
			name: "apiUrl",
			want: "apiURL",
		},
		// Test initialismsAsWords functionality
		{
			name:               "apiUrl",
			initialismsAsWords: true,
			want:               "apiUrl",
		},
		{
			name:               "apiURL",
			initialismsAsWords: true,
			want:               "apiUrl",
		},
		{
			name:               "userId",
			initialismsAsWords: true,
			want:               "userId",
		},
		{
			name:               "userID",
			initialismsAsWords: true,
			want:               "userId",
		},
		{
			name:               "UserId",
			initialismsAsWords: true,
			want:               "UserId",
		},
		{
			name:               "UserID",
			initialismsAsWords: true,
			want:               "UserId",
		},
		{
			name:               "readJson",
			initialismsAsWords: true,
			want:               "readJson",
		},
		{
			name:               "readJSON",
			initialismsAsWords: true,
			want:               "readJson",
		},
		{
			name:               "httpMethod",
			initialismsAsWords: true,
			want:               "httpMethod",
		},
		{
			name:               "HttpMethod",
			initialismsAsWords: true,
			want:               "HttpMethod",
		},
		{
			name:               "HTTPMethod",
			initialismsAsWords: true,
			want:               "HttpMethod",
		},
		{
			name:               "id",
			initialismsAsWords: true,
			want:               "id",
		},
		{
			name:               "Id",
			initialismsAsWords: true,
			want:               "Id",
		},
		{
			name:               "ID",
			initialismsAsWords: true,
			want:               "Id",
		},
		{
			name:               "ids",
			initialismsAsWords: true,
			want:               "ids",
		},
		{
			name:               "Ids",
			initialismsAsWords: true,
			want:               "Ids",
		},
		{
			name:               "IDS",
			initialismsAsWords: true,
			want:               "Ids",
		},
		{
			name:               "userIds",
			initialismsAsWords: true,
			want:               "userIds",
		},
		{
			name:               "userIDs",
			initialismsAsWords: true,
			want:               "userIds",
		},
		{
			name:               "XmlApi",
			initialismsAsWords: true,
			want:               "XmlApi",
		},
		{
			name:               "SqlUrl",
			initialismsAsWords: true,
			want:               "SqlUrl",
		},
		{
			name:               "foo_id",
			initialismsAsWords: true,
			want:               "fooId",
		},
		{
			name:               "foo_bar",
			initialismsAsWords: true,
			want:               "fooBar",
		},
		{
			name:               "foo_WiFi",
			initialismsAsWords: true,
			want:               "fooWiFi",
		},
		// Test initialismsAsWords with allowlist
		{
			name:               "readJSON",
			allowlist:          []string{"JSON"},
			initialismsAsWords: true,
			want:               "readJSON",
		},
		// Test initialismsAsWords with blocklist
		{
			name:               "apiUrl",
			blocklist:          []string{"URL"},
			initialismsAsWords: true,
			want:               "apiURL",
		},
		{
			name:               "ApiUrl",
			blocklist:          []string{"API"},
			initialismsAsWords: true,
			want:               "APIUrl",
		},
		// Test initialismsAsWords with skipInitialismNameChecks
		{
			name:                     "readJSON",
			skipInitialismNameChecks: true,
			initialismsAsWords:       true,
			want:                     "readJSON",
		},
		{
			name:                     "readJson",
			skipInitialismNameChecks: true,
			initialismsAsWords:       true,
			want:                     "readJson",
		},
	}
	for _, test := range tests {
		got := rule.Name(test.name, test.allowlist, test.blocklist, test.skipInitialismNameChecks, test.initialismsAsWords)
		if got != test.want {
			t.Errorf("name(%q, allowlist=%v, blocklist=%v, skipInitialismNameChecks=%v, initialismsAsWords=%v) = %q, want %q",
				test.name, test.allowlist, test.blocklist, test.skipInitialismNameChecks, test.initialismsAsWords, got, test.want)
		}
	}
}
