package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// With many account tabs, the bar must stay on one row when the terminal is wide
// enough, and wrap onto multiple rows when it is not — so tabs never overflow
// off-screen (the 29-mailbox case).
func TestWrapTabRowsWrapsWhenNarrow(t *testing.T) {
	var tabs []string
	for i := 0; i < 29; i++ {
		tabs = append(tabs, "[ user@example.com ]") // ~20 cells each
	}

	wide := wrapTabRows(tabs, 10000)
	if got := lipgloss.Height(wide); got != 1 {
		t.Fatalf("wide terminal: expected 1 tab row, got %d", got)
	}

	narrow := wrapTabRows(tabs, 120)
	if got := lipgloss.Height(narrow); got <= 1 {
		t.Fatalf("narrow terminal: expected tabs to wrap to >1 row, got %d", got)
	}

	// No single rendered row may exceed the requested width.
	for _, line := range strings.Split(narrow, "\n") {
		if w := lipgloss.Width(line); w > 120 {
			t.Fatalf("row exceeds width: %d > 120 (%q)", w, line)
		}
	}
}

func TestWrapTabRowsEdgeCases(t *testing.T) {
	if got := wrapTabRows(nil, 80); got != "" {
		t.Fatalf("empty input should render empty, got %q", got)
	}
	// A zero/unknown width must not panic and must still produce one row.
	one := wrapTabRows([]string{"ALL", "a@b.c"}, 0)
	if lipgloss.Height(one) != 1 {
		t.Fatalf("two short tabs at default width should be 1 row, got %d", lipgloss.Height(one))
	}
}
