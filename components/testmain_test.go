package components

import (
	"os"
	"testing"

	"github.com/jorgerojas26/lazysql/internal/i18n"
)

// TestMain switches the UI language to English for all component tests.
// Assertions in these tests were written against the original English UI
// strings; running them under the (now default) Chinese locale would break
// text-based assertions unrelated to the feature under test.
func TestMain(m *testing.M) {
	if err := i18n.SetLanguage("en"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
