package configs

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// WriteDefault prints a neofetch-compatible config for the built-in defaults.
func WriteDefault(w io.Writer) {
	cfg := DefaultConfig()
	fmt.Fprintln(w, "# gf config. The print_info block chooses the lines, in order.")
	fmt.Fprintln(w, "print_info() {")
	for _, line := range cfg.InfoList {
		switch len(line) {
		case 0:
			continue
		case 1:
			fmt.Fprintf(w, "    info %s\n", line[0])
		default:
			fmt.Fprintf(w, "    info \"%s\" %s\n", line[0], line[len(line)-1])
		}
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	writeFields(w, reflect.ValueOf(cfg))
}

func writeFields(w io.Writer, v reflect.Value) {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Anonymous && value.Kind() == reflect.Struct {
			writeFields(w, value)
			continue
		}
		tag := field.Tag.Get("gf")
		if tag == "" || tag == "info_list" || value.Kind() != reflect.String {
			continue
		}
		text := value.String()
		if text == "" || strings.ContainsAny(text, " ()") {
			fmt.Fprintf(w, "%s=%s\n", tag, quoteIfNeeded(text))
			continue
		}
		fmt.Fprintf(w, "%s=%q\n", tag, text)
	}
}

func quoteIfNeeded(text string) string {
	if text == "" {
		return `""`
	}
	if strings.HasPrefix(text, "(") {
		return text
	}
	return fmt.Sprintf("%q", text)
}
