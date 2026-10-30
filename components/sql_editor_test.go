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
// completion when the cursor sits right after "table." with an empty prefix
// (previously acceptCompletion bailed out and just closed the popup).
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

	// The first suggested column should replace the whole "tenant_llm." token.
	completion := editor.acItems[0].Text
	editor.acceptCompletion()

	want := "SELECT " + completion + " FROM tenant_llm;"
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

	want := "SELECT " + completion + " FROM tenant_llm;"
	if got := editor.GetText(); got != want {
		t.Fatalf("after Enter via handler: got %q, want %q", got, want)
	}
}
