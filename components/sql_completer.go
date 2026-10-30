package components

import (
	"sort"
	"strings"
	"unicode"

	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/jorgerojas26/lazysql/internal/i18n"
)

// CompletionItem represents a single autocomplete suggestion.
type CompletionItem struct {
	Text        string // what gets inserted
	Description string // displayed help text
}

// Autocompleter manages SQL keywords and schema-aware completions.
type Autocompleter struct {
	keywords []CompletionItem
	tables   []CompletionItem
	columns  map[string][]CompletionItem // table name -> columns
}

// NewAutocompleter creates an autocompleter with built-in SQL keywords.
func NewAutocompleter() *Autocompleter {
	return &Autocompleter{
		keywords: builtinKeywords(),
		tables:   nil,
		columns:  make(map[string][]CompletionItem),
	}
}

// SetTables updates the list of known table names.
func (a *Autocompleter) SetTables(tables []string) {
	a.tables = make([]CompletionItem, len(tables))
	for i, t := range tables {
		a.tables[i] = CompletionItem{Text: t, Description: i18n.T("table")}
	}
}

// SetColumns sets the columns for a given table.
func (a *Autocompleter) SetColumns(table string, columns []string) {
	items := make([]CompletionItem, len(columns))
	for i, c := range columns {
		items[i] = CompletionItem{Text: c, Description: i18n.T("column")}
	}
	a.columns[strings.ToLower(table)] = items
}

// HasColumns reports whether the loader has delivered a result for a table.
// An empty slice is still a loaded result and must not trigger repeated
// database requests.
func (a *Autocompleter) HasColumns(table string) bool {
	_, ok := a.columns[strings.ToLower(table)]
	return ok
}

// GetCompletions returns completion items matching the given prefix using
// fuzzy search (same ranking as the tree: exact > prefix > substring > fuzzy).
// If tableHint is non-empty, it prioritizes columns from that table.
// When prefix is empty but tableHint is set (e.g. user typed "table."), all
// columns for that table are returned.
func (a *Autocompleter) GetCompletions(prefix string, tableHint string) []CompletionItem {
	// When prefix is empty but a table is specified, show all columns for that table.
	if prefix == "" && tableHint != "" {
		if cols, ok := a.columns[strings.ToLower(tableHint)]; ok {
			return cols
		}
		return nil
	}
	if prefix == "" {
		return nil
	}

	type scoredCandidate struct {
		item  CompletionItem
		score int // lower = better match (0=exact, 1-99=prefix, 100+=substr/fuzzy)
		order int // priority group (0=column, 1=table, 2=keyword)
	}

	var candidates []scoredCandidate
	seen := make(map[string]bool)
	lowerPrefix := strings.ToLower(prefix)

	// Collect candidates with a score and dedup
	tryAdd := func(items []CompletionItem, order int) {
		for _, item := range items {
			key := strings.ToLower(item.Text)
			if seen[key] {
				continue
			}
			// Use the same prioritization as the tree search
			rank := fuzzy.RankMatch(lowerPrefix, key)
			if rank < 0 {
				continue
			}
			score := prioritizeResult(lowerPrefix, key, rank)
			candidates = append(candidates, scoredCandidate{item, score, order})
			seen[key] = true
		}
	}

	// 1. Columns from the hinted table (highest priority)
	if tableHint != "" {
		if cols, ok := a.columns[strings.ToLower(tableHint)]; ok {
			tryAdd(cols, 0)
		}
	}

	// 2. Table names
	tryAdd(a.tables, 1)

	// 3. Keywords (lower priority)
	tryAdd(a.keywords, 2)

	// Sort by score ascending (lower = better), then by priority group
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score < candidates[j].score
		}
		return candidates[i].order < candidates[j].order
	})

	// Limit to top 20 results
	maxResults := 20
	if len(candidates) < maxResults {
		maxResults = len(candidates)
	}

	results := make([]CompletionItem, maxResults)
	for i := 0; i < maxResults; i++ {
		results[i] = candidates[i].item
	}

	return results
}

// GetAllCompletions returns all known keywords, tables, and column names
// matching the prefix, in order of relevance.
func (a *Autocompleter) GetAllCompletions(prefix string) []CompletionItem {
	return a.GetCompletions(prefix, "")
}

// extractPrefix extracts the word being typed at the cursor position.
// text is the full text, cursorPos is the byte offset of the cursor.
func extractPrefix(text string, cursorPos int) string {
	prefix, _ := extractCompletionContext(text, cursorPos)
	return prefix
}

// extractCompletionContext splits the current word segment at the cursor into
// a column prefix and an optional table name. For "table.col|" it returns
// ("col", "table"). For "prefix|" alone it returns ("prefix", "").
func extractCompletionContext(text string, cursorPos int) (prefix, tableName string) {
	if cursorPos <= 0 || cursorPos > len(text) {
		return "", ""
	}

	// Walk backward from cursor to find the start of the current token.
	// Unlike extractPrefix, we do NOT stop at '.' — we want table.col intact.
	start := cursorPos - 1
	for start >= 0 {
		ch := rune(text[start])
		if unicode.IsSpace(ch) || ch == ';' || ch == ',' || ch == '(' || ch == ')' {
			break
		}
		start--
	}
	start++

	if start >= cursorPos {
		return "", ""
	}

	segment := text[start:cursorPos]

	// Check for table.col pattern (last dot in the segment)
	if dotIdx := strings.LastIndex(segment, "."); dotIdx >= 0 {
		tableName = segment[:dotIdx]
		prefix = segment[dotIdx+1:]
		return prefix, tableName
	}

	// Regular prefix without a dot
	return segment, ""
}

// resolveAliases scans the query text for table alias definitions.
// Uses the lexer-based SQL context scanner instead of the old string-search
// heuristics, so it correctly handles subqueries, CTEs, and comma lists.
// Returns aliases visible at cursorPos (aliases from outer scopes are visible,
// but aliases from deeper subqueries are not).
func resolveAliases(sql string, cursorPos int) map[string]string {
	ctx := scanSQLContext(sql)
	depth := cursorDepth(sql, cursorPos)
	result := make(map[string]string)

	// Table ref aliases, filtered by subquery depth.
	for _, ref := range ctx.tableRefs {
		if ref.alias != "" && ref.depth <= depth {
			result[strings.ToLower(ref.alias)] = ref.name
		}
	}
	// CTE names are always visible at depth 0.
	for cteName := range ctx.CTEs {
		if 0 <= depth {
			result[cteName] = cteName
		}
	}

	return result
}

// extractTableHint uses the lexer-based context scanner to find the most
// recently mentioned table at the cursor's subquery depth.  This correctly
// handles subquery scoping — a cursor inside a subquery won't see tables
// from the outer query as the "hint".
func extractTableHint(text string, cursorPos int) string {
	if cursorPos > len(text) {
		cursorPos = len(text)
	}

	ctx := scanSQLContext(text)
	depth := cursorDepth(text, cursorPos)

	// Find the table reference closest to (but before) the cursor at the
	// same nesting depth.
	var best tableRef
	bestPos := -1
	for _, ref := range ctx.tableRefs {
		if ref.depth == depth && ref.pos < cursorPos && ref.pos > bestPos {
			best = ref
			bestPos = ref.pos
		}
	}

	if best.name != "" {
		// If the cursor follows a comma (e.g. "FROM users, |"), there is
		// no active table hint — the user is starting a new table name.
		trimmed := strings.TrimRight(text[:cursorPos], " \t\n\r")
		if len(trimmed) > 0 && trimmed[len(trimmed)-1] == ',' {
			return ""
		}
		return best.name
	}
	return ""
}

// builtinKeywords returns the list of SQL keywords for autocomplete.
func builtinKeywords() []CompletionItem {
	return []CompletionItem{
		{Text: "SELECT", Description: i18n.T("Retrieve rows from a table")},
		{Text: "FROM", Description: i18n.T("Specify the source table")},
		{Text: "WHERE", Description: i18n.T("Filter results")},
		{Text: "AND", Description: i18n.T("Logical AND")},
		{Text: "OR", Description: i18n.T("Logical OR")},
		{Text: "NOT", Description: i18n.T("Logical NOT")},
		{Text: "IN", Description: i18n.T("Check membership")},
		{Text: "IS", Description: i18n.T("Null comparison")},
		{Text: "NULL", Description: i18n.T("Null value")},
		{Text: "LIKE", Description: i18n.T("Pattern matching")},
		{Text: "BETWEEN", Description: i18n.T("Range check")},
		{Text: "EXISTS", Description: i18n.T("Check existence")},
		{Text: "AS", Description: i18n.T("Alias")},
		{Text: "ON", Description: i18n.T("Join condition")},
		{Text: "JOIN", Description: i18n.T("Join tables")},
		{Text: "INNER", Description: i18n.T("Inner join")},
		{Text: "LEFT", Description: i18n.T("Left join")},
		{Text: "RIGHT", Description: i18n.T("Right join")},
		{Text: "OUTER", Description: i18n.T("Outer join")},
		{Text: "CROSS", Description: i18n.T("Cross join")},
		{Text: "FULL", Description: i18n.T("Full join")},
		{Text: "NATURAL", Description: i18n.T("Natural join")},
		{Text: "USING", Description: i18n.T("Join using columns")},
		{Text: "INSERT", Description: i18n.T("Insert rows")},
		{Text: "INTO", Description: i18n.T("Specify target table")},
		{Text: "VALUES", Description: i18n.T("Row values")},
		{Text: "UPDATE", Description: i18n.T("Update rows")},
		{Text: "SET", Description: i18n.T("Set column values")},
		{Text: "DELETE", Description: i18n.T("Delete rows")},
		{Text: "CREATE", Description: i18n.T("Create object")},
		{Text: "TABLE", Description: i18n.T("Create table")},
		{Text: "INDEX", Description: i18n.T("Create index")},
		{Text: "VIEW", Description: i18n.T("Create view")},
		{Text: "ALTER", Description: i18n.T("Modify object")},
		{Text: "DROP", Description: i18n.T("Delete object")},
		{Text: "ADD", Description: i18n.T("Add column")},
		{Text: "COLUMN", Description: i18n.T("Column keyword")},
		{Text: "CONSTRAINT", Description: i18n.T("Add constraint")},
		{Text: "PRIMARY KEY", Description: i18n.T("Primary key")},
		{Text: "FOREIGN KEY", Description: i18n.T("Foreign key")},
		{Text: "UNIQUE", Description: i18n.T("Unique constraint")},
		{Text: "CHECK", Description: i18n.T("Check constraint")},
		{Text: "DEFAULT", Description: i18n.T("Default value")},
		{Text: "REFERENCES", Description: i18n.T("Foreign key reference")},
		{Text: "ORDER BY", Description: i18n.T("Order results")},
		{Text: "ASC", Description: i18n.T("Ascending order")},
		{Text: "DESC", Description: i18n.T("Descending order")},
		{Text: "GROUP BY", Description: i18n.T("Group results")},
		{Text: "HAVING", Description: i18n.T("Filter groups")},
		{Text: "LIMIT", Description: i18n.T("Limit rows")},
		{Text: "OFFSET", Description: i18n.T("Offset rows")},
		{Text: "UNION", Description: i18n.T("Combine queries")},
		{Text: "ALL", Description: i18n.T("All results")},
		{Text: "DISTINCT", Description: i18n.T("Distinct rows")},
		{Text: "CASE", Description: i18n.T("Case expression")},
		{Text: "WHEN", Description: i18n.T("Case when")},
		{Text: "THEN", Description: i18n.T("Case then")},
		{Text: "ELSE", Description: i18n.T("Case else")},
		{Text: "END", Description: i18n.T("End expression")},
		{Text: "BEGIN", Description: i18n.T("Start transaction")},
		{Text: "COMMIT", Description: i18n.T("Commit transaction")},
		{Text: "ROLLBACK", Description: i18n.T("Rollback transaction")},
		{Text: "EXPLAIN", Description: i18n.T("Explain query plan")},
		{Text: "DESCRIBE", Description: i18n.T("Describe table")},
		{Text: "SHOW", Description: i18n.T("Show objects")},
		{Text: "USE", Description: i18n.T("Use database")},
		{Text: "TRUNCATE", Description: i18n.T("Truncate table")},
		{Text: "REPLACE", Description: i18n.T("Replace rows")},
		{Text: "CALL", Description: i18n.T("Call procedure")},
		{Text: "WITH", Description: i18n.T("Common table expression")},
		{Text: "RECURSIVE", Description: i18n.T("Recursive CTE")},
		{Text: "FETCH", Description: i18n.T("Fetch rows")},
		{Text: "NEXT", Description: i18n.T("Fetch next")},
		{Text: "ROWS", Description: i18n.T("Row count")},
		{Text: "ONLY", Description: i18n.T("Fetch only")},
		{Text: "TOP", Description: i18n.T("Select top")},
		{Text: "COUNT", Description: i18n.T("Count rows")},
		{Text: "SUM", Description: i18n.T("Sum values")},
		{Text: "AVG", Description: i18n.T("Average values")},
		{Text: "MIN", Description: i18n.T("Minimum value")},
		{Text: "MAX", Description: i18n.T("Maximum value")},
		{Text: "COALESCE", Description: i18n.T("First non-null")},
		{Text: "IFNULL", Description: i18n.T("Null fallback")},
		{Text: "CAST", Description: i18n.T("Type cast")},
		{Text: "CONVERT", Description: i18n.T("Type convert")},
		{Text: "CROSS JOIN", Description: i18n.T("Cross join")},
		{Text: "NATURAL JOIN", Description: i18n.T("Natural join")},
		{Text: "INNER JOIN", Description: i18n.T("Inner join")},
		{Text: "LEFT JOIN", Description: i18n.T("Left join")},
		{Text: "RIGHT JOIN", Description: i18n.T("Right join")},
		{Text: "FULL JOIN", Description: i18n.T("Full join")},
		{Text: "OUTER JOIN", Description: i18n.T("Outer join")},
		{Text: "ASC", Description: i18n.T("Ascending")},
		{Text: "DESC", Description: i18n.T("Descending")},
		{Text: "COUNT(*)", Description: i18n.T("Count all rows")},
		{Text: "DISTINCT", Description: i18n.T("Unique values only")},
	}
}
