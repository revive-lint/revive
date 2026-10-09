package fixtures

func testAllowlistBlocklist() {
	// JSON is in allowlist, so both readJSON and readJson are accepted
	readJSON := "json"
	readJson := "json"
	_ = readJSON
	_ = readJson

	// URL is in blocklist, so it is enforced as initialism (apiURL)
	apiURL := "url"
	apiUrl := "url" // MATCH /var apiUrl should be apiURL/
	_ = apiURL
	_ = apiUrl

	// ID is not in blocklist, so it is treated as a word (userId)
	userId := 1
	userID := 2 // MATCH /var userID should be userId/
	_ = userId
	_ = userID
}
