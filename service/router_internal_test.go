package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeRecipientAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		recipient   string
		expected    string
		wasPrefixed bool
	}{
		{
			name:        "plain liquid address",
			recipient:   "lq1qqvzq56gesrms7nm4yyqzvckvazhj7c5ytxwvpfku4yf7nchpth8uzyxrw0j3734fe6psh5wc0kqhlklqyydrjvyru899pszjz",
			expected:    "lq1qqvzq56gesrms7nm4yyqzvckvazhj7c5ytxwvpfku4yf7nchpth8uzyxrw0j3734fe6psh5wc0kqhlklqyydrjvyru899pszjz",
			wasPrefixed: false,
		},
		{
			name:        "prefixed address",
			recipient:   "liquidnetwork:lq1qqvzq56gesrms7nm4yyqzvckvazhj7c5ytxwvpfku4yf7nchpth8uzyxrw0j3734fe6psh5wc0kqhlklqyydrjvyru899pszjz",
			expected:    "lq1qqvzq56gesrms7nm4yyqzvckvazhj7c5ytxwvpfku4yf7nchpth8uzyxrw0j3734fe6psh5wc0kqhlklqyydrjvyru899pszjz",
			wasPrefixed: true,
		},
		{
			name:        "case-insensitive prefix",
			recipient:   "LiquidNetwork:lq1qqvzq56gesrms7nm4yyqzvckvazhj7c5ytxwvpfku4yf7nchpth8uzyxrw0j3734fe6psh5wc0kqhlklqyydrjvyru899pszjz",
			expected:    "lq1qqvzq56gesrms7nm4yyqzvckvazhj7c5ytxwvpfku4yf7nchpth8uzyxrw0j3734fe6psh5wc0kqhlklqyydrjvyru899pszjz",
			wasPrefixed: true,
		},
		{
			name:        "short non-prefixed string",
			recipient:   "liquid",
			expected:    "liquid",
			wasPrefixed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			normalized, wasPrefixed := normalizeRecipientAddress(tt.recipient)
			assert.Equal(t, tt.expected, normalized)
			assert.Equal(t, tt.wasPrefixed, wasPrefixed)
		})
	}
}
