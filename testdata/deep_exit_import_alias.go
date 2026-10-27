package fixtures

import (
	fl "flag"
	stdlog "log"
	sys "os"
	sc "syscall"
)

func exitThroughRenamedImports() {
	sys.Exit(1)                          // MATCH /calls to os.Exit only in main() or init() functions/
	sc.Exit(1)                           // MATCH /calls to syscall.Exit only in main() or init() functions/
	stdlog.Fatal("boom")                 // MATCH /calls to log.Fatal only in main() or init() functions/
	stdlog.Panicf("%s", "boom")          // MATCH /calls to log.Panicf only in main() or init() functions/
	fl.Parse()                           // MATCH /calls to flag.Parse only in main() or init() functions; move the call or refactor to use flag.NewFlagSet with flag.ContinueOnError/
	fl.NewFlagSet("cmd", fl.ExitOnError) // MATCH /calls to flag.NewFlagSet with flag.ExitOnError only in main() or init() functions/
}

func noExitThroughRenamedImports() {
	fl.NewFlagSet("cmd", fl.ContinueOnError)
	stdlog.Print("not an exit")
}

func init() {
	sys.Exit(1)
}

func main() {
	stdlog.Fatal("boom")
}
