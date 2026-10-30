package components

import (
	"sync"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestSQLEditorSubscribeAndPublishConcurrently(t *testing.T) {
	editor := NewSQLEditor("")
	start := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		for range 1000 {
			editor.Subscribe()
		}
	}()

	go func() {
		defer wg.Done()
		<-start
		for range 1000 {
			editor.Publish("key", "value")
		}
	}()

	close(start)
	wg.Wait()

	subscriber := editor.Subscribe()
	editor.Publish("key", "value")

	change := <-subscriber
	if change.Key != "key" || change.Value != "value" {
		t.Fatalf("unexpected state change: %#v", change)
	}
}

// TestAutocompleteEnterOnTableQualified ensures Enter accepts a column
// completion when the cursor sits right after "table." with an empty prefix,
// and that the "table." qualifier is preserved (tenant_llm.id, not id).
func TestAutocompleteEnterOnTableQualified(t *testing.T) {
	editor := NewSQLEditor("")
	editor.SetTables([]string{"tenant_llm"})
	editor.SetColumns("tenant_llm", []string{"id", "create_time", "create_date", "tenant_id"})

	editor.SetText("SELECT tenant_llm. FROM tenant_llm;", true)
	// Place the cursor right after the dot: "SELECT tenant_llm.| FROM tenant_llm;"
	editor.cx = len("SELECT tenant_llm.")

	editor.acTableHint = ""
	editor.triggerAutocomplete()
	if !editor.acVisible || len(editor.acItems) == 0 {
		t.Fatalf("expected autocomplete popup with column items, visible=%v items=%d", editor.acVisible, len(editor.acItems))
	}

	// The first suggested column should be inserted after "tenant_llm.".
	completion := editor.acItems[0].Text
	editor.acceptCompletion()

	want := "SELECT tenant_llm." + completion + " FROM tenant_llm;"
	if got := editor.GetText(); got != want {
		t.Fatalf("after Enter completion: got %q, want %q", got, want)
	}
}

// TestAutocompleteEnterKeepsTableQualifier covers the partial-prefix variant
// "tenant_llm.cre" — only "cre" is replaced, "tenant_llm." stays.
func TestAutocompleteEnterKeepsTableQualifier(t *testing.T) {
	editor := NewSQLEditor("")
	editor.SetTables([]string{"tenant_llm"})
	editor.SetColumns("tenant_llm", []string{"id", "create_time", "create_date"})

	editor.SetText("SELECT tenant_llm.cre FROM tenant_llm;", true)
	editor.cx = len("SELECT tenant_llm.cre")

	editor.acTableHint = ""
	editor.triggerAutocomplete()
	if !editor.acVisible || len(editor.acItems) == 0 {
		t.Fatalf("expected autocomplete popup, visible=%v items=%d", editor.acVisible, len(editor.acItems))
	}

	completion := editor.acItems[0].Text
	editor.acceptCompletion()

	want := "SELECT tenant_llm." + completion + " FROM tenant_llm;"
	if got := editor.GetText(); got != want {
		t.Fatalf("after Enter completion: got %q, want %q", got, want)
	}
}

// TestAutocompleteEnterViaInputHandler drives the real InputHandler with a
// KeyEnter event, as a user would press it while the popup is visible.
func TestAutocompleteEnterViaInputHandler(t *testing.T) {
	editor := NewSQLEditor("")
	editor.SetTables([]string{"tenant_llm"})
	editor.SetColumns("tenant_llm", []string{"id", "create_time", "create_date"})

	editor.SetText("SELECT tenant_llm. FROM tenant_llm;", true)
	editor.cx = len("SELECT tenant_llm.")
	editor.SetIsFocused(true)

	editor.acTableHint = ""
	editor.triggerAutocomplete()
	if !editor.acVisible || len(editor.acItems) == 0 {
		t.Fatalf("expected autocomplete popup, visible=%v items=%d", editor.acVisible, len(editor.acItems))
	}

	completion := editor.acItems[0].Text
	handler := editor.InputHandler()
	handler(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), func(p tview.Primitive) {})

	want := "SELECT tenant_llm." + completion + " FROM tenant_llm;"
	if got := editor.GetText(); got != want {
		t.Fatalf("after Enter via handler: got %q, want %q", got, want)
	}
}

// TestAutocompleteEnterBarePrefix ensures a plain (unqualified) prefix is
// still replaced as before, without a spurious "table." qualifier.
func TestAutocompleteEnterBarePrefix(t *testing.T) {
	editor := NewSQLEditor("")
	editor.SetTables([]string{"tenant_llm"})
	editor.SetColumns("tenant_llm", []string{"id", "create_time", "create_date"})

	editor.SetText("SELECT cre FROM tenant_llm;", true)
	editor.cx = len("SELECT cre")

	editor.acTableHint = ""
	editor.triggerAutocomplete()
	if !editor.acVisible || len(editor.acItems) == 0 {
		t.Fatalf("expected autocomplete popup, visible=%v items=%d", editor.acVisible, len(editor.acItems))
	}

	completion := editor.acItems[0].Text
	editor.acceptCompletion()

	want := "SELECT " + completion + " FROM tenant_llm;"
	if got := editor.GetText(); got != want {
		t.Fatalf("after Enter completion: got %q, want %q", got, want)
	}
}
