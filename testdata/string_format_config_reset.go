package fixtures

import "fmt"

func stringFormatConfigReset() {
	fmt.Errorf("bad") // MATCH /old configuration was applied/
}
