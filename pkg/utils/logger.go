package utils

import "strings"

// Add this helper function to redact PII fields from SQL logs
func RedactPIIFromSQL(sql string) string {
	// Redact email values in JSON (e.g., "value":"charlie.brown@example.com")
	sql = RedactJSONField(sql, `"value":"`, `"`, "email")
	// Redact gender values in JSON (e.g., "gender":"other")
	sql = RedactJSONField(sql, `"gender":"`, `"`, "")
	return sql
}

// Helper to redact a JSON field value in a SQL string
func RedactJSONField(sql, prefix, suffix, system string) string {
	// If system is set, only redact when system matches (for email)
	if system != "" {
		// Find all occurrences of system":"email" and redact the next "value":"..."
		systemKey := `"system":"` + system + `"`
		idx := 0
		for {
			sysIdx := strings.Index(sql[idx:], systemKey)
			if sysIdx == -1 {
				break
			}
			sysIdx += idx
			valIdx := strings.Index(sql[sysIdx:], `"value":"`)
			if valIdx == -1 {
				break
			}
			valIdx += sysIdx + len(`"value":"`)
			endIdx := strings.Index(sql[valIdx:], `"`)
			if endIdx == -1 {
				break
			}
			endIdx += valIdx
			sql = sql[:valIdx] + "***REDACTED***" + sql[endIdx:]
			idx = valIdx + len("***REDACTED***")
		}
		return sql
	}
	// Otherwise, redact all occurrences of the field
	idx := 0
	for {
		start := strings.Index(sql[idx:], prefix)
		if start == -1 {
			break
		}
		start += idx + len(prefix)
		end := strings.Index(sql[start:], suffix)
		if end == -1 {
			break
		}
		end += start
		sql = sql[:start] + "***REDACTED***" + sql[end:]
		idx = start + len("***REDACTED***")
	}
	return sql
}
