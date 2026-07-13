package getlongopts_test

import (
	"testing"

	"github.com/dxloc/go-getlongopts"
)

func TestParser(t *testing.T) {
	var a = ""
	var b = ""
	var c = ""
	var d = false

	opts := []getlongopts.LongOption{
		{
			Long:    "long-a",
			Short:   "a",
			ArgType: getlongopts.ArgTypeFile,
			SetFn: func(value string) {
				a = value
			},
			Description: "Option (-a|--long-a), file argument",
		},
		{
			Long:    "long-b",
			ArgType: getlongopts.ArgTypeDir,
			SetFn: func(value string) {
				b = value
			},
			Description: "Option (--long-b), directory argument",
		},
		{
			Long:    "long-c",
			Short:   "c",
			ArgType: getlongopts.ArgTypeOther,
			SetFn: func(value string) {
				c = value
			},
			Description: "Option (--long-c), text argument",
		},
		{
			Long:    "long-d",
			Short:   "d",
			ArgType: getlongopts.ArgTypeNone,
			SetFn: func(value string) {
				d = true
			},
			Description: "Option (-d|--long-d), no argument",
		},
	}

	p := getlongopts.NewParser(opts)

	// In real application, we run p.Parse(os.Args) instead of a []string
	testArgs := []string{"progname", "-a", "A", "--long-b", "B", "-c", "C", "-d"}

	if _, e := p.Parse(testArgs); e != nil {
		t.Error("Error: ", e)
	}
	if a != "A" {
		t.Error("Expected a = A")
	}
	if b != "B" {
		t.Error("Expected b = B")
	}
	if c != "C" {
		t.Error("Expected c = C")
	}
	if d != true {
		t.Error("Expected true, got ", d)
	}
}
