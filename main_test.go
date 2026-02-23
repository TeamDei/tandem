package main

import (
	"reflect"
	"testing"
)

func TestParseWords(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "empty file",
			raw:  "",
			want: nil,
		},
		{
			name: "whitespace only",
			raw:  "   \n\n\t\n",
			want: nil,
		},
		{
			name: "single line",
			raw:  "hello world",
			want: []string{"hello", "world"},
		},
		{
			name: "two lines produce sentinel between them",
			raw:  "hello world\nfoo bar",
			want: []string{"hello", "world", "\n", "foo", "bar"},
		},
		{
			// Each blank line contributes one sentinel; two non-empty lines flanking
			// one blank line get one sentinel between them.
			name: "blank line between produces one sentinel",
			raw:  "hello world\n\nfoo bar",
			want: []string{"hello", "world", "\n", "\n", "foo", "bar"},
		},
		{
			// Three blank lines → three sentinels in between.
			name: "multiple consecutive blank lines each produce a sentinel",
			raw:  "hello\n\n\n\nworld",
			want: []string{"hello", "\n", "\n", "\n", "\n", "world"},
		},
		{
			name: "CRLF line endings normalised",
			raw:  "hello world\r\nfoo bar",
			want: []string{"hello", "world", "\n", "foo", "bar"},
		},
		{
			// A trailing newline means the last "line" is blank; because prevHadWords
			// is true we emit one sentinel at the end.
			name: "trailing newline appends a sentinel",
			raw:  "\nhello world\n",
			want: []string{"hello", "world", "\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseWords([]byte(tt.raw))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseWords(%q)\n  got  %v\n  want %v", tt.raw, got, tt.want)
			}
		})
	}
}
