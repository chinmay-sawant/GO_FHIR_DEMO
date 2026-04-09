package utils //nolint:revive

import "strings"

// RedactPIIFromSQL redacts PII fields from SQL logs.
func RedactPIIFromSQL(sql string) string {
	// Redact email values in JSON (e.g., "value":"charlie.brown@example.com")
	sql = RedactJSONField(sql, `"value":"`, `"`, "email")
	// Redact gender values in JSON (e.g., "gender":"other")
	sql = RedactJSONField(sql, `"gender":"`, `"`, "")
	return sql
}

// RedactJSONField redacts a JSON field value in a SQL string.
func RedactJSONField(sql, prefix, suffix, system string) string {
	// If system is set, only redact when system matches (for email)
	if system != "" {
		return redactMatchedField(sql, `"system":"`+system+`"`, `"value":"`, `"`)
	}
	return redactMatchedField(sql, "", prefix, suffix)
}

func redactMatchedField(sql, matcher, prefix, suffix string) string {
	var builder strings.Builder
	builder.Grow(len(sql))

	searchFrom := 0
	for {
		matchIdx := 0
		if matcher != "" {
			matchIdx = strings.Index(sql[searchFrom:], matcher)
			if matchIdx == -1 {
				break
			}
			matchIdx += searchFrom
		}

		start := strings.Index(sql[searchFrom:], prefix)
		if start == -1 {
			break
		}
		start += searchFrom + len(prefix)
		end := strings.Index(sql[start:], suffix)
		if end == -1 {
			break
		}
		end += start

		if matcher != "" && matchIdx > start {
			searchFrom = matchIdx + len(matcher)
			continue
		}

		builder.WriteString(sql[searchFrom:start])
		builder.WriteString("***REDACTED***")
		searchFrom = end
	}

	builder.WriteString(sql[searchFrom:])
	return builder.String()
}
