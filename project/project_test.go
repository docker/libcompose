package project

import (
	"fmt"
	"strings"
	"testing"
)

func TestEventEquality(t *testing.T) {
	if fmt.Sprintf("%s", EventServiceStart) != "Started" ||
		fmt.Sprintf("%v", EventServiceStart) != "Started" {
		t.Fatalf("EventServiceStart String() doesn't work: %s %v", EventServiceStart, EventServiceStart)
	}

	if fmt.Sprintf("%s", EventServiceStart) != fmt.Sprintf("%s", EventServiceUp) {
		t.Fatal("Event messages do not match")
	}

	if EventServiceStart == EventServiceUp {
		t.Fatal("Events match")
	}
}

func TestParseWithBadContent(t *testing.T) {
	p := NewProject(&Context{
		ComposeBytes: []byte("garbage"),
	})

	err := p.Parse()
	if err == nil {
		t.Fatal("Should have failed parse")
	}

	if !strings.HasPrefix(err.Error(), "yaml: unmarshal errors") {
		t.Fatal("Should have failed parse", err)
	}
}

func TestParseWithGoodContent(t *testing.T) {
	p := NewProject(&Context{
		ComposeBytes: []byte("not-garbage:\n  image: foo"),
	})

	err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}
