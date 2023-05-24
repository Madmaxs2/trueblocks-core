package statePkg

import (
	"strings"
	"testing"
)

func TestParse_Encoding(t *testing.T) {
	// 0xcdba2fd40000000000000000000000000000000000000000000000000000000000007a69
	if p, err := Parse(`0xcdba2fd40000000000000000000000000000000000000000000000000000000000007a69`); err != nil {
		t.Fatal(err)
	} else {
		if s := p.Encoded; s != `0xcdba2fd40000000000000000000000000000000000000000000000000000000000007a69` {
			t.Fatal("wrong encoded value", s)
		}
	}
}

func TestParse_Selector(t *testing.T) {
	// Invalid selector
	if _, err := Parse(`0xcdba()`); !strings.Contains(err.Error(), errInvalidSelector.Error()) {
		t.Fatal("expected errInvalidSelector, got:", err)
	}

	// No arguments
	if parsed, err := Parse(`0xcdba2fd4()`); err != nil {
		t.Fatal(err)
	} else {
		if value := parsed.SelectorCall.Selector.Value; value != `0xcdba2fd4` {
			t.Fatal("wrong selector", value)
		}
		if argsLen := len(parsed.SelectorCall.Arguments); argsLen != 0 {
			t.Fatal("wrong number of arguments", argsLen)
		}
	}

	// Arguments
	if parsed, err := Parse(`0xcdba2fd4(1, true, false, 0xbeef, "string")`); err != nil {
		t.Fatal(err)
	} else {
		if value := parsed.SelectorCall.Selector.Value; value != `0xcdba2fd4` {
			t.Fatal("wrong selector", value)
		}
		if argValue := *parsed.SelectorCall.Arguments[0].Int; argValue != 1 {
			t.Fatal("wrong #1 input value", argValue)
		}
		if argValue := *parsed.SelectorCall.Arguments[1].Boolean; argValue != true {
			t.Fatal("wrong #2 input value", argValue)
		}
		if argValue := *parsed.SelectorCall.Arguments[2].Boolean; argValue != false {
			t.Fatal("wrong #3 input value", argValue)
		}
		if argValue := *parsed.SelectorCall.Arguments[3].Hex.String; argValue != "0xbeef" {
			t.Fatal("wrong #4 input value", argValue)
		}
		if argValue := *parsed.SelectorCall.Arguments[4].String; argValue != "string" {
			t.Fatal("wrong #5 input value", argValue)
		}
	}

	// Hex parsing
	if parsed, err := Parse(`0xcdba2fd4(0xdeadbeef, 0x6982508145454ce325ddbe47a25d4ec3d23119a1)`); err != nil {
		t.Fatal(err)
	} else {
		if argsLen := len(parsed.SelectorCall.Arguments); argsLen != 2 {
			t.Fatal("wrong inputs length:", argsLen)
		}

		if *parsed.SelectorCall.Arguments[0].Hex.String != "0xdeadbeef" {
			t.Fatal("should parse into Hex.String")
		}
		if parsed.SelectorCall.Arguments[1].Hex.Address.Hex() != "0x6982508145454ce325ddbe47a25d4ec3d23119a1" {
			t.Fatal("should parse into Hex.Address")
		}
	}
}

func TestParse_Function(t *testing.T) {
	// Invalid selector
	if _, err := Parse(`111()`); err == nil {
		t.Fatal("expected parsing error")
	}

	// No arguments
	if parsed, err := Parse(`transfer()`); err != nil {
		t.Fatal(err)
	} else {
		if value := parsed.FunctionNameCall.Name; value != `transfer` {
			t.Fatal("wrong function name", value)
		}
		if argsLen := len(parsed.FunctionNameCall.Arguments); argsLen != 0 {
			t.Fatal("wrong number of arguments", argsLen)
		}
	}

	// Correct Solidity identifiers
	if parsed, err := Parse(`$dollar$_underscoreCamelCase125__()`); err != nil {
		t.Fatal(err)
	} else {
		if value := parsed.FunctionNameCall.Name; value != `$dollar$_underscoreCamelCase125__` {
			t.Fatal("wrong function name", value)
		}
	}

	// Arguments
	if parsed, err := Parse(`something(1, true, false, 0xbeef, "string")`); err != nil {
		t.Fatal(err)
	} else {
		if value := parsed.FunctionNameCall.Name; value != `something` {
			t.Fatal("wrong selector", value)
		}
		if argValue := *parsed.FunctionNameCall.Arguments[0].Int; argValue != 1 {
			t.Fatal("wrong #1 input value", argValue)
		}
		if argValue := *parsed.FunctionNameCall.Arguments[1].Boolean; argValue != true {
			t.Fatal("wrong #2 input value", argValue)
		}
		if argValue := *parsed.FunctionNameCall.Arguments[2].Boolean; argValue != false {
			t.Fatal("wrong #3 input value", argValue)
		}
		if argValue := *parsed.FunctionNameCall.Arguments[3].Hex.String; argValue != "0xbeef" {
			t.Fatal("wrong #4 input value", argValue)
		}
		if argValue := *parsed.FunctionNameCall.Arguments[4].String; argValue != "string" {
			t.Fatal("wrong #5 input value", argValue)
		}
	}

	// Hex parsing
	if parsed, err := Parse(`somethingElse(0xdeadbeef, 0x6982508145454ce325ddbe47a25d4ec3d23119a1)`); err != nil {
		t.Fatal(err)
	} else {
		if argsLen := len(parsed.FunctionNameCall.Arguments); argsLen != 2 {
			t.Fatal("wrong inputs length:", argsLen)
		}

		if *parsed.FunctionNameCall.Arguments[0].Hex.String != "0xdeadbeef" {
			t.Fatal("should parse into Hex.String")
		}
		if parsed.FunctionNameCall.Arguments[1].Hex.Address.Hex() != "0x6982508145454ce325ddbe47a25d4ec3d23119a1" {
			t.Fatal("should parse into Hex.Address")
		}
	}
}
