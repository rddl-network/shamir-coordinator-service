package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/syndtr/goleveldb/leveldb"
)

func isPrintable(b []byte) bool {
	for _, r := range string(b) {
		if !unicode.IsPrint(r) && r != '\n' && r != '\t' {
			return false
		}
	}
	return true
}

func formatKey(key []byte) string {
	if isPrintable(key) {
		return string(key)
	}
	// Try prefix + binary suffix (e.g. Request/SendTokens/\x00\x00\x00\x00\x00\x00\x00\x01)
	for i := len(key) - 1; i >= 0; i-- {
		if isPrintable(key[:i]) {
			suffix := key[i:]
			var id uint64
			if len(suffix) == 8 {
				for _, b := range suffix {
					id = id<<8 | uint64(b)
				}
				return fmt.Sprintf("%s[id=%d]", string(key[:i]), id)
			}
		}
	}
	return fmt.Sprintf("%x", key)
}

func formatValue(val []byte) string {
	if !isPrintable(val) {
		return fmt.Sprintf("(binary) %x", val)
	}
	// Try to pretty-print JSON
	var js interface{}
	if err := json.Unmarshal(val, &js); err == nil {
		pretty, err := json.MarshalIndent(js, "", "  ")
		if err == nil {
			return string(pretty)
		}
	}
	return string(val)
}

func main() {
	dbPath := "./data"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening LevelDB: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var sb strings.Builder
	sb.WriteString("# LevelDB Dump\n\n")
	sb.WriteString("**Path:** `" + dbPath + "`\n\n")
	sb.WriteString("---\n\n")

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	count := 0
	for iter.Next() {
		count++
		key := formatKey(iter.Key())
		val := formatValue(iter.Value())

		sb.WriteString(fmt.Sprintf("## Entry %d\n\n", count))
		sb.WriteString(fmt.Sprintf("**Key:** `%s`\n\n", key))
		sb.WriteString("**Value:**\n\n```json\n")
		sb.WriteString(val)
		sb.WriteString("\n```\n\n")
		sb.WriteString("---\n\n")
	}

	if err := iter.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "Iterator error: %v\n", err)
		os.Exit(1)
	}

	if count == 0 {
		sb.WriteString("_(empty database)_\n")
	} else {
		sb.WriteString(fmt.Sprintf("_Total entries: %d_\n", count))
	}

	outputPath := "./data/db-dump.md"
	if err := os.WriteFile(outputPath, []byte(sb.String()), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %d entries to %s\n", count, outputPath)
}
