package gui

import "testing"

func TestWindowsTerminalFontInheritance(t *testing.T) {
	raw := []byte(`{
		// user comment
		"profiles": {
			"defaults": {
				"font": { "face": "Cascadia Mono", "size": 12, }
			},
			"list": [
				{
					"guid": "{61c54bbd-c2c6-5271-96e7-009a87ff44bf}",
					"font": { "face": "JetBrains Mono" }
				}
			]
		}
	}`)
	font, matched, ok := parseWTDocument(raw, "{61C54BBD-C2C6-5271-96E7-009A87FF44BF}")
	if !ok || !matched {
		t.Fatalf("parsed ok=%v matched=%v", ok, matched)
	}
	if got := finishTermFont(font); got != "JetBrains Mono (12pt)" {
		t.Fatalf("profile font = %q", got)
	}

	font, matched, ok = parseWTDocument(raw, "{missing}")
	if !ok || matched {
		t.Fatalf("fallback ok=%v matched=%v", ok, matched)
	}
	if got := finishTermFont(font); got != "Cascadia Mono (12pt)" {
		t.Fatalf("default font = %q", got)
	}

	empty := resolvedFont{}
	applyWTFragment([]byte(`{"profiles":[{"updates":"{abc}","font":{"face":"Iosevka","size":11}}]}`), "{abc}", &empty)
	if got := finishTermFont(empty); got != "Iosevka (11pt)" {
		t.Fatalf("fragment = %q", got)
	}
	// A profile font outranks a fragment.
	applyWTFragment([]byte(`{"profiles":[{"guid":"{abc}","font":{"face":"Consolas"}}]}`), "{abc}", &font)
	if font.face != "Cascadia Mono" {
		t.Fatalf("fragment overrode defaults: %q", font.face)
	}
}

func TestLegacyWTFontAndMintty(t *testing.T) {
	raw := []byte(`{"profiles":{"list":[{"guid":"{abc}","fontFace":"Sarasa Gothic SC","fontSize":14}]}}`)
	font, matched, ok := parseWTDocument(raw, "abc")
	if !ok || !matched {
		t.Fatal("legacy profile was not parsed")
	}
	if got := finishTermFont(font); got != "Sarasa Gothic SC (14pt)" {
		t.Fatalf("legacy = %q", got)
	}
	name, size := parseMinttyConfig("# comment\nFont=DejaVu Sans Mono\nFontHeight=10\n")
	if name != "DejaVu Sans Mono" || size != "10" {
		t.Fatalf("mintty = %q %q", name, size)
	}
	n, s := parseConEmuFont(`<value name="FontName" type="string" data="Consolas"/><value name="FontSize" type="ulong" data="16"/>`)
	if n != "Consolas" || s != "16" {
		t.Fatalf("conemu = %q %q", n, s)
	}
}
