package fixtures

import (
	stdlog "log"
	sys "os"
)

func exitsThroughRenamedOS(n int) {
	if n == 0 {
		sys.Exit(0)
	}
	exitsThroughRenamedOS(n - 1)
}

func exitsThroughRenamedLog(n int) {
	if n == 0 {
		stdlog.Fatal("done")
	}
	exitsThroughRenamedLog(n - 1)
}

func noExitThroughRenamedLog(n int) {
	if n == 0 {
		stdlog.Print("not an exit")
	}
	noExitThroughRenamedLog(n - 1) // MATCH /unconditional recursive call/
}
