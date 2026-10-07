package fixtures

type customOptionTagAfterReset struct {
	Field string `validate:"myOption"` // MATCH /unknown option "myOption" in validate tag/
}
