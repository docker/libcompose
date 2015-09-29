package project

import (
	"fmt"
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
