package fixtures

import (
	"fmt"
	"log"
	"runtime/trace"
	"testing"
)

func unnecessaryFormat(t *testing.T, b *testing.B, f *testing.F) {
	logger := log.New(nil, "", 0)

	fmt.Appendf(nil, "no format") // MATCH /unnecessary use of formatting function "fmt.Appendf", you can replace it with "fmt.Append"/
	fmt.Errorf("no format")       // MATCH /unnecessary use of formatting function "fmt.Errorf", you can replace it with "errors.New"/
	fmt.Fprintf(nil, "no format") // MATCH /unnecessary use of formatting function "fmt.Fprintf", you can replace it with "fmt.Fprint"/
	fmt.Fscanf(nil, "no format")  // MATCH /unnecessary use of formatting function "fmt.Fscanf", you can replace it with "fmt.Fscan" or "fmt.Fscanln"/
	fmt.Printf("no format")       // MATCH /unnecessary use of formatting function "fmt.Printf", you can replace it with "fmt.Print" or "fmt.Println"/
	fmt.Scanf("no format")        // MATCH /unnecessary use of formatting function "fmt.Scanf", you can replace it with "fmt.Scan"/
	fmt.Sprintf("no format")      // MATCH /unnecessary use of formatting function "fmt.Sprintf", you can replace it with "fmt.Sprint" or just the string itself"/
	fmt.Sscanf("", "no format")   // MATCH /unnecessary use of formatting function "fmt.Sscanf", you can replace it with "fmt.Sscan"/
	// standard logging functions
	log.Fatalf("no format")    // MATCH /unnecessary use of formatting function "log.Fatalf", you can replace it with "log.Fatal"/
	log.Panicf("no format")    // MATCH /unnecessary use of formatting function "log.Panicf", you can replace it with "log.Panic"/
	log.Printf("no format")    // MATCH /unnecessary use of formatting function "log.Printf", you can replace it with "log.Print"/
	logger.Fatalf("no format") // MATCH /unnecessary use of formatting function "logger.Fatalf", you can replace it with "logger.Fatal"/
	logger.Panicf("no format") // MATCH /unnecessary use of formatting function "logger.Panicf", you can replace it with "logger.Panic"/
	logger.Printf("no format") // MATCH /unnecessary use of formatting function "logger.Printf", you can replace it with "logger.Print"/
	// standard testing functions
	t.Errorf("no format") // MATCH /unnecessary use of formatting function "t.Errorf", you can replace it with "t.Error"/
	t.Fatalf("no format") // MATCH /unnecessary use of formatting function "t.Fatalf", you can replace it with "t.Fatal"/
	t.Logf("no format")   // MATCH /unnecessary use of formatting function "t.Logf", you can replace it with "t.Log"/
	t.Skipf("no format")  // MATCH /unnecessary use of formatting function "t.Skipf", you can replace it with "t.Skip"/
	b.Errorf("no format") // MATCH /unnecessary use of formatting function "b.Errorf", you can replace it with "b.Error"/
	b.Fatalf("no format") // MATCH /unnecessary use of formatting function "b.Fatalf", you can replace it with "b.Fatal"/
	b.Logf("no format")   // MATCH /unnecessary use of formatting function "b.Logf", you can replace it with "b.Log"/
	b.Skipf("no format")  // MATCH /unnecessary use of formatting function "b.Skipf", you can replace it with "b.Skip"/
	f.Errorf("no format") // MATCH /unnecessary use of formatting function "f.Errorf", you can replace it with "f.Error"/
	f.Fatalf("no format") // MATCH /unnecessary use of formatting function "f.Fatalf", you can replace it with "f.Fatal"/
	f.Logf("no format")   // MATCH /unnecessary use of formatting function "f.Logf", you can replace it with "f.Log"/
	f.Skipf("no format")  // MATCH /unnecessary use of formatting function "f.Skipf", you can replace it with "f.Skip"/
	// standard trace functions
	trace.Logf(nil, "http", "no format", nil) // MATCH /unnecessary use of formatting function "trace.Logf", you can replace it with "trace.Log"/

	fmt.Appendf(nil, "format %d", 0)
	fmt.Errorf("format %d", 0)
	fmt.Fprintf(nil, "format %d", 0)
	fmt.Fscanf(nil, "format %d", 0)
	fmt.Printf("format %d", 0)
	fmt.Scanf("format %d", 0)
	fmt.Sprintf("format %d", 0)
	fmt.Sscanf("", "format %d", 0)
	// standard logging functions
	log.Fatalf("format %d", 0)
	log.Panicf("format %d", 0)
	log.Printf("format %d", 0)
	logger.Fatalf("format %d", 0)
	logger.Panicf("format %d", 0)
	logger.Printf("format %d", 0)
	// standard testing functions
	t.Errorf("format %d", 0)
	t.Fatalf("format %d", 0)
	t.Logf("format %d", 0)
	t.Skipf("format %d", 0)
	b.Errorf("format %d", 0)
	b.Fatalf("format %d", 0)
	b.Logf("format %d", 0)
	b.Skipf("format %d", 0)
	f.Errorf("format %d", 0)
	f.Fatalf("format %d", 0)
	f.Logf("format %d", 0)
	f.Skipf("format %d", 0)
	// standard trace functions
	trace.Logf(nil, "http", "format %d", nil)

	// test with multiline string argument
	// MATCH:77 /unnecessary use of formatting function "fmt.Appendf", you can replace it with "fmt.Append"/
	fmt.Appendf(nil, `no 
		format`)
	fmt.Appendf(nil, `format 
	%d`, 0)

	// test with directives spelled as escape sequences
	fmt.Sprintf("\x25d", 0)
	fmt.Sprintf("\u0025d", 0)
	fmt.Sprintf("\045d", 0)
	fmt.Sprintf("\x6eo format") // MATCH /unnecessary use of formatting function "fmt.Sprintf", you can replace it with "fmt.Sprint" or just the string itself"/
}
