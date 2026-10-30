package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultLanguageIsChinese(t *testing.T) {
	if got := CurrentLanguage(); got != DefaultLanguage {
		t.Fatalf("CurrentLanguage() = %q, want %q", got, DefaultLanguage)
	}
}

func TestTChinese(t *testing.T) {
	if got := T("Yes"); got != "是" {
		t.Errorf(`T("Yes") = %q, want "是"`, got)
	}
	if got := T("Save"); got != "保存" {
		t.Errorf(`T("Save") = %q, want "保存"`, got)
	}
}

func TestSwitchToEnglish(t *testing.T) {
	if err := SetLanguage("en"); err != nil {
		t.Fatal(err)
	}
	defer SetLanguage(DefaultLanguage)

	if got := T("Yes"); got != "Yes" {
		t.Errorf(`T("Yes") after en = %q, want "Yes" (English fallback)`, got)
	}
	if got := CurrentLanguage(); got != "en" {
		t.Errorf("CurrentLanguage() = %q, want en", got)
	}
}

func TestTUnknownKeyReturnsKey(t *testing.T) {
	if got := T("definitely-not-a-key-xyz"); got != "definitely-not-a-key-xyz" {
		t.Errorf("unknown key returned %q, want the key itself", got)
	}
}

func TestTfFormatsTranslatedString(t *testing.T) {
	SetLanguage(DefaultLanguage)
	got := Tf("%d-%d of %d rows", 1, 50, 1000)
	if got != "1-50 / 共 1000 行" {
		t.Errorf("Tf = %q, want %q", got, "1-50 / 共 1000 行")
	}
}

func TestTfFormatsEnglishWhenMissing(t *testing.T) {
	SetLanguage("en")
	got := Tf("%d-%d of %d rows", 1, 50, 1000)
	if got != "1-50 of 1000 rows" {
		t.Errorf("Tf en = %q, want %q", got, "1-50 of 1000 rows")
	}
}

func TestLoadCustomFileOverrides(t *testing.T) {
	SetLanguage(DefaultLanguage)
	defer SetLanguage(DefaultLanguage)

	dir := t.TempDir()
	path := filepath.Join(dir, "custom.json")
	content := `{"Yes": "行", "New Key": "自定义值"}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := LoadCustomFile(path); err != nil {
		t.Fatal(err)
	}
	if got := T("Yes"); got != "行" {
		t.Errorf(`T("Yes") after custom = %q, want "行"`, got)
	}
	if got := T("New Key"); got != "自定义值" {
		t.Errorf(`T("New Key") = %q, want "自定义值"`, got)
	}
	// Existing built-in keys not overridden by the custom file stay intact.
	if got := T("Save"); got != "保存" {
		t.Errorf(`T("Save") = %q, want "保存"`, got)
	}
}

func TestLoadCustomFileMissing(t *testing.T) {
	if err := LoadCustomFile(filepath.Join(t.TempDir(), "nope.json")); err == nil {
		t.Fatal("expected error for missing custom language file")
	}
}

func TestCommandAndGroupPrefixKeys(t *testing.T) {
	SetLanguage(DefaultLanguage)
	if got := T("cmd.Quit"); got != "退出" {
		t.Errorf(`T("cmd.Quit") = %q, want "退出"`, got)
	}
	if got := T("keymap.group.table"); got != "表格" {
		t.Errorf(`T("keymap.group.table") = %q, want "表格"`, got)
	}
}

// TestKeybindDescriptionKeys guards the help modal / status bar strings:
// help_modal translates bind.Description at render time, so those keys are
// not visible to static extraction and must be covered explicitly.
func TestKeybindDescriptionKeys(t *testing.T) {
	SetLanguage(DefaultLanguage)
	cases := map[string]string{
		"Focus table":                         "聚焦表格",
		"Quit":                                "退出",
		"Search":                              "搜索",
		"Check for updates / upgrade LazySQL": "检查更新 / 升级 LazySQL",
		"Export to CSV":                       "导出为 CSV",
		"Toggle value menu to put values like NULL, EMPTY or DEFAULT": "切换值菜单以填入 NULL、EMPTY 或 DEFAULT",
	}
	for key, want := range cases {
		if got := T(key); got != want {
			t.Errorf(`T(%q) = %q, want %q`, key, got, want)
		}
	}
}
