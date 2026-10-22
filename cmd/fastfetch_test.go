package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ffgan/gf/internal/utils"
)

func TestSeparatorAndVersion(t *testing.T) {
	h := &hardware{}
	h.Stdout = utils.ON
	h.length = 4
	h.SeparatorString = "-"
	if got := fastfetchDetectors["separator"](h); got != "----" {
		t.Fatalf("separator = %q", got)
	}
	h.CustomText = "hello"
	if fastfetchDetectors["custom"](h) != "hello" {
		t.Fatal(fastfetchDetectors["custom"](h))
	}
	if !strings.HasPrefix(fastfetchDetectors["version"](h), "gf ") {
		t.Fatal(fastfetchDetectors["version"](h))
	}
}

func TestJSONOmitsBlanks(t *testing.T) {
	h := &hardware{}
	h.Stdout = utils.ON
	h.InfoList = [][]string{{"break"}, {"Version", "version"}}
	var buf bytes.Buffer
	h.writeJSON(&buf)
	text := buf.String()
	if !strings.Contains(text, `"key": "version"`) || strings.Contains(text, "line_break") {
		t.Fatalf("json = %s", text)
	}
}
