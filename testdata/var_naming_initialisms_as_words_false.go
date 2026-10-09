package fixtures

const HTTPRes = 200

func readJSON() {
	return
}

type SSHConfig struct {
	keyPath string
}

func testIdentifiers() {
	apiURL := "http://example.com"
	_ = apiURL

	userID := 1
	_ = userID

	HTTPMethod := "POST"
	_ = HTTPMethod

	id := 10
	ID := 20
	_ = id
	_ = ID
}
