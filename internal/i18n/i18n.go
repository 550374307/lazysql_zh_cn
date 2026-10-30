// Package i18n provides lightweight, dependency-free internationalization
// support for LazySQL.
//
// Translation keys are the English source strings themselves. When a key has
// no translation, T returns the key unchanged, so the UI degrades gracefully
// to English. Two built-in languages ship with the binary (embedded via
// go:embed): "en" (no-op) and "zh" (Simplified Chinese, the default).
//
// Users can supply a custom language file (any valid JSON mapping English
// source strings to their translation). The custom file is merged on top of
// the built-in language and takes precedence. This allows arbitrary languages
// or terminology overrides without recompiling.
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "embed"
)

//go:embed locales/en.json
var enLocale []byte

//go:embed locales/zh.json
var zhLocale []byte

// DefaultLanguage is the language used when no configuration overrides it.
const DefaultLanguage = "zh"

var (
	mu           sync.RWMutex
	language     = DefaultLanguage
	translations = map[string]string{}
)

// init loads the default language (Simplified Chinese) so that any code path
// running before LoadConfig — such as flag usage text — already sees the
// configured default instead of untranslated English keys.
func init() {
	translations = parseLocale(zhLocale)
}

// SetLanguage selects the built-in language pack. Supported built-in codes
// are "zh" (Simplified Chinese) and "en" (English). Any other code resets the
// built-in translations to English (an empty map) so a custom language file
// can take over completely.
func SetLanguage(lang string) error {
	mu.Lock()
	defer mu.Unlock()

	lang = strings.ToLower(strings.TrimSpace(lang))
	if lang == "" {
		lang = DefaultLanguage
	}

	switch lang {
	case "zh", "zh-cn", "zh_CN", "zh-CN":
		translations = parseLocale(zhLocale)
	case "en", "en-us", "en_US", "en-US", "":
		translations = parseLocale(enLocale)
	default:
		// Unknown built-in code: fall back to English and let a custom
		// language file provide the actual translations.
		translations = parseLocale(enLocale)
	}

	language = lang
	return nil
}

// CurrentLanguage returns the active language code.
func CurrentLanguage() string {
	mu.RLock()
	defer mu.RUnlock()
	return language
}

// parseLocale decodes an embedded locale JSON payload into a map. Malformed
// embedded data returns an empty map; the file is maintained in-repo so this
// should never happen.
func parseLocale(data []byte) map[string]string {
	locale := map[string]string{}
	if len(data) > 0 {
		_ = json.Unmarshal(data, &locale)
	}
	return locale
}

// LoadCustomFile merges a user-supplied language file on top of the current
// translations. Entries in the custom file override built-in ones; keys that
// do not exist in the custom file keep their current translation. The file
// must be valid JSON mapping English source strings to translated strings.
// Environment variables in the path (${env:VAR}) are expanded for
// consistency with the TOML config loader.
func LoadCustomFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	expanded := os.Expand(path, func(key string) string {
		if envKey, found := strings.CutPrefix(key, "env:"); found {
			return os.Getenv(envKey)
		}
		return "${" + key + "}"
	})

	data, err := os.ReadFile(expanded)
	if err != nil {
		return fmt.Errorf("reading language file %q: %w", expanded, err)
	}

	custom := map[string]string{}
	if err := json.Unmarshal(data, &custom); err != nil {
		return fmt.Errorf("parsing language file %q: %w", expanded, err)
	}

	mu.Lock()
	defer mu.Unlock()
	for key, value := range custom {
		translations[key] = value
	}
	return nil
}

// T translates key into the active language. When no translation exists the
// key itself is returned, which keeps English as the built-in fallback.
func T(key string) string {
	mu.RLock()
	value, ok := translations[key]
	mu.RUnlock()
	if !ok {
		return key
	}
	return value
}

// Tf translates key and formats it with fmt.Sprintf-style arguments.
// Arguments are only applied when a translation exists; otherwise the raw
// English source string is formatted, preserving its own verbs.
func Tf(key string, args ...any) string {
	mu.RLock()
	value, ok := translations[key]
	mu.RUnlock()
	if !ok {
		return fmt.Sprintf(key, args...)
	}
	return fmt.Sprintf(value, args...)
}

// LoadDir is a helper for tests and embedding tools: it returns the absolute
// path of a locale file inside the i18n package directory. It is not used by
// the runtime.
func localePath(name string) string {
	return filepath.Join("locales", name)
}
