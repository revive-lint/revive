package fixtures

import (
	"flag"
	"log"
	stdlog "log"
	"os"
	"syscall"
)

func exitAfterDefer() {
	defer println("cleanup")
	os.Exit(1) // MATCH /os.Exit after a defer statement prevents deferred calls from running/
}

func exitBeforeDeferIsFine() {
	os.Exit(1)
	defer println("cleanup")
}

func exitWithoutDeferIsFine() {
	os.Exit(1)
}

func fatalAfterDefer() {
	defer println("cleanup")
	log.Fatal("boom")   // MATCH /log.Fatal after a defer statement prevents deferred calls from running/
	log.Fatalf("%d", 1) // MATCH /log.Fatalf after a defer statement prevents deferred calls from running/
	log.Fatalln("boom") // MATCH /log.Fatalln after a defer statement prevents deferred calls from running/
}

func syscallExitAfterDefer() {
	defer println("cleanup")
	syscall.Exit(1) // MATCH /syscall.Exit after a defer statement prevents deferred calls from running/
}

func panicAfterDeferIsFine() {
	defer println("cleanup")
	// deferred functions still run while a panic unwinds the stack.
	log.Panic("boom")
	panic("boom")
}

func exitInBranchAfterDefer() {
	defer println("cleanup")
	if len(os.Args) > 1 {
		os.Exit(1) // MATCH /os.Exit after a defer statement prevents deferred calls from running/
	}
}

func nestedFuncLitHasOwnScope() {
	go func() {
		defer println("cleanup")
		os.Exit(1) // MATCH /os.Exit after a defer statement prevents deferred calls from running/
	}()
}

func deferredFuncLitExitIsFine() {
	defer func() {
		os.Exit(0)
	}()
	println("work")
}

func exitInFuncLitBeforeOuterDefer() {
	f := func() {
		os.Exit(1)
	}
	defer println("cleanup")
	f()
}

func flagParseAfterDeferIsFine() {
	defer println("cleanup")
	flag.Parse()
	fs := flag.NewFlagSet("name", flag.ExitOnError)
	_ = fs
}

func aliasedImportExitAfterDefer() {
	defer println("cleanup")
	stdlog.Fatal("boom") // MATCH /stdlog.Fatal after a defer statement prevents deferred calls from running/
}
