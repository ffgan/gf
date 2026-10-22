package gui

import (
	"encoding/json"
	"strconv"
	"strings"
)

type resolvedFont struct {
	face    string
	size    float64
	hasSize bool
	facePri int
	sizePri int
}

func formatTermFont(name, size, unit string) string {
	name = strings.TrimSpace(name)
	size = strings.TrimSpace(size)
	if name == "" {
		return ""
	}
	if size == "" {
		return name
	}
	if unit == "" {
		unit = "pt"
	}
	return name + " (" + size + unit + ")"
}

func formatFontSize(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}

func stripJSONC(s string) string {
	s = strings.TrimPrefix(s, "\uFEFF")
	var b strings.Builder
	b.Grow(len(s))
	inStr := false
	esc := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr {
			b.WriteByte(c)
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
				continue
			}
			if c == '"' {
				inStr = false
			}
			continue
		}
		if c == '"' {
			inStr = true
			b.WriteByte(c)
			continue
		}
		if c == '/' && i+1 < len(s) && (s[i+1] == '/' || s[i+1] == '*') {
			if s[i+1] == '/' {
				i += 2
				for i < len(s) && s[i] != '\n' {
					i++
				}
				if i < len(s) {
					b.WriteByte('\n')
				}
				continue
			}
			i += 2
			for i+1 < len(s) && !(s[i] == '*' && s[i+1] == '/') {
				i++
			}
			if i+1 < len(s) {
				i++
			}
			continue
		}
		b.WriteByte(c)
	}
	return stripTrailingCommas(b.String())
}

func stripTrailingCommas(s string) string {
	var b strings.Builder
	inStr := false
	esc := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inStr {
			b.WriteByte(c)
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
				continue
			}
			if c == '"' {
				inStr = false
			}
			continue
		}
		if c == '"' {
			inStr = true
			b.WriteByte(c)
			continue
		}
		if c == ',' {
			j := i + 1
			for j < len(s) && (s[j] == ' ' || s[j] == '\n' || s[j] == '\r' || s[j] == '\t') {
				j++
			}
			if j < len(s) && (s[j] == '}' || s[j] == ']') {
				continue
			}
		}
		b.WriteByte(c)
	}
	return b.String()
}

func parseWTDocument(data []byte, profileID string) (font resolvedFont, matched, ok bool) {
	var root map[string]any
	if json.Unmarshal([]byte(stripJSONC(string(data))), &root) != nil {
		return font, false, false
	}
	profiles, _ := root["profiles"].(map[string]any)
	if profiles == nil {
		return font, false, false
	}
	if profileID != "" {
		if p := findWTProfile(asProfileList(profiles["list"]), profileID, false); p != nil {
			applyProfileFont(p, 4, &font)
			matched = true
		}
	}
	if defaults, isMap := profiles["defaults"].(map[string]any); isMap {
		applyProfileFont(defaults, 3, &font)
	}
	return font, matched, true
}

func applyWTFragment(data []byte, profileID string, font *resolvedFont) {
	var root map[string]any
	if json.Unmarshal([]byte(stripJSONC(string(data))), &root) != nil {
		return
	}
	list := fragmentProfiles(root["profiles"])
	if p := findWTProfile(list, profileID, true); p != nil {
		applyProfileFont(p, 2, font)
		return
	}
	if p := findWTProfile(list, profileID, false); p != nil {
		applyProfileFont(p, 1, font)
	}
}

func fragmentProfiles(v any) []map[string]any {
	switch p := v.(type) {
	case []any:
		return asProfileList(p)
	case map[string]any:
		return asProfileList(p["list"])
	default:
		return nil
	}
}

func asProfileList(v any) []map[string]any {
	raw, _ := v.([]any)
	if len(raw) == 0 {
		return nil
	}
	list := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			list = append(list, m)
		}
	}
	return list
}

func findWTProfile(list []map[string]any, id string, updates bool) map[string]any {
	key := "guid"
	if updates {
		key = "updates"
	}
	for _, p := range list {
		s, _ := p[key].(string)
		if s != "" && sameProfileID(s, id) {
			return p
		}
	}
	return nil
}

func sameProfileID(a, b string) bool {
	return strings.EqualFold(strings.Trim(a, "{} \t"), strings.Trim(b, "{} \t"))
}

func applyProfileFont(profile map[string]any, pri int, font *resolvedFont) {
	if obj, ok := profile["font"].(map[string]any); ok {
		if face, ok := obj["face"].(string); ok {
			considerFont(font, face, 0, false, pri)
		}
		if n, ok := jsonFloat(obj["size"]); ok {
			considerFont(font, "", n, true, pri)
		}
	}
	if face, ok := profile["fontFace"].(string); ok {
		considerFont(font, face, 0, false, pri)
	}
	if n, ok := jsonFloat(profile["fontSize"]); ok {
		considerFont(font, "", n, true, pri)
	}
}

func considerFont(font *resolvedFont, face string, size float64, hasSize bool, pri int) {
	face = strings.TrimSpace(face)
	if face != "" && pri > font.facePri {
		font.face = face
		font.facePri = pri
	}
	if hasSize && pri > font.sizePri {
		font.size = size
		font.hasSize = true
		font.sizePri = pri
	}
}

func jsonFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func parseMinttyConfig(text string) (font, size string) {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "font":
			if font == "" {
				font = strings.TrimSpace(val)
			}
		case "fontheight":
			if size == "" {
				size = strings.TrimSpace(val)
			}
		}
	}
	return font, size
}

func parseConEmuFont(text string) (name, size string) {
	return xmlDataAttr(text, "FontName"), xmlDataAttr(text, "FontSize")
}

func xmlDataAttr(text, key string) string {
	needle := `name="` + key + `"`
	i := strings.Index(text, needle)
	if i < 0 {
		return ""
	}
	rest := text[i:]
	j := strings.Index(rest, `data="`)
	if j < 0 || j > 240 {
		return ""
	}
	rest = rest[j+len(`data="`):]
	k := strings.Index(rest, `"`)
	if k < 0 {
		return ""
	}
	return rest[:k]
}

func finishTermFont(font resolvedFont) string {
	if font.face == "" {
		font.face = "Cascadia Mono"
	}
	if !font.hasSize {
		font.size = 12
	}
	return formatTermFont(font.face, formatFontSize(font.size), "pt")
}
