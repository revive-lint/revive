package fixtures

const HttpRes = 200
const HTTPRes = 200 // MATCH /const HTTPRes should be HttpRes/

func readJson() {
	return
}

func readJSON() { // MATCH /func readJSON should be readJson/
	return
}

type SshConfig struct {
	keyPath string
}

type SSHConfig struct { // MATCH /type SSHConfig should be SshConfig/
	keyPath string
}

func testIdentifiers() {
	apiUrl := "http://example.com"
	apiURL := "http://example.com" // MATCH /var apiURL should be apiUrl/
	_ = apiUrl
	_ = apiURL

	userId := 1
	userID := 2 // MATCH /var userID should be userId/
	_ = userId
	_ = userID

	httpMethod := "GET"
	HTTPMethod := "POST" // MATCH /var HTTPMethod should be HttpMethod/
	_ = httpMethod
	_ = HTTPMethod

	id := 10
	Id := 20
	ID := 30 // MATCH /var ID should be Id/
	_ = id
	_ = Id
	_ = ID
}
