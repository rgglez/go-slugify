/*

Copyright 2026 Rodolfo González González

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the “Software”), to deal in
the Software without restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the
Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

package slugify

import (
	"testing"
)

func sp(s string) *string { return &s }
func bp(b bool) *bool     { return &b }

func must(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func slug(s string, opts *Options) string {
	return must(Slugify(s, opts))
}

func TestMain(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Foo Bar", "foo-bar"},
		{"foo bar baz", "foo-bar-baz"},
		{"foo bar ", "foo-bar"},
		{"       foo bar", "foo-bar"},
		{"[foo] [bar]", "foo-bar"},
		{"Foo ÿ", "foo-y"},
		{"FooBar", "foo-bar"},
		{"fooBar", "foo-bar"},
		{"UNICORNS AND RAINBOWS", "unicorns-and-rainbows"},
		{"Foo & Bar", "foo-and-bar"},
		{"Foo & Bar", "foo-and-bar"},
		{"Hællæ, hva skjera?", "haellae-hva-skjera"},
		{"Foo Bar2", "foo-bar2"},
		{"I ♥ Dogs", "i-love-dogs"},
		{"Déjà Vu!", "deja-vu"},
		{"fooBar 123 $#%", "foo-bar-123"},
		{"foo🦄", "foo-unicorn"},
		{"🦄🦄🦄", "unicorn-unicorn-unicorn"},
		{"foo&bar", "foo-and-bar"},
		{"foo360BAR", "foo360-bar"},
		{"FOO360", "foo-360"},
		{"FOOBar", "foo-bar"},
		{"APIs", "apis"},
		{"APISection", "api-section"},
		{"Util APIs", "util-apis"},
	}
	for _, c := range cases {
		got := slug(c.in, nil)
		if got != c.want {
			t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPossessivesAndContractions(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Conway's Law", "conways-law"},
		{"Conway's", "conways"},
		{"Don't Repeat Yourself", "dont-repeat-yourself"},
		{"my parents' rules", "my-parents-rules"},
		{"it-s-hould-not-modify-t-his", "it-s-hould-not-modify-t-his"},
		{"Sindre’s app", "sindres-app"},
		{"can’t stop", "cant-stop"},
		{"won’t work", "wont-work"},
	}
	for _, c := range cases {
		got := slug(c.in, nil)
		if got != c.want {
			t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCustomSeparator(t *testing.T) {
	cases := []struct {
		in, want string
		opts     *Options
	}{
		{"foo bar", "foo_bar", &Options{Separator: sp("_")}},
		{"aaa bbb", "aaabbb", &Options{Separator: sp("")}},
		{"BAR&baz", "bar_and_baz", &Options{Separator: sp("_")}},
		{"Déjà Vu!", "deja-vu", &Options{Separator: sp("-")}},
		{"UNICORNS AND RAINBOWS!", "unicorns@and@rainbows", &Options{Separator: sp("@")}},
		{"[foo] [bar]", "foo.bar", &Options{Separator: sp(".")}},
		{"a   b   c", "a__b__c", &Options{Separator: sp("__")}},
		{"a____b", "a__b", &Options{Separator: sp("__")}},
		{"__a__b__", "a__b", &Options{Separator: sp("__")}},
		{"foo---bar", "foo---bar", &Options{Separator: sp("---")}},
	}
	for _, c := range cases {
		got := slug(c.in, c.opts)
		if got != c.want {
			t.Errorf("Slugify(%q, sep=%q) = %q, want %q", c.in, *c.opts.Separator, got, c.want)
		}
	}
}

func TestCustomReplacements(t *testing.T) {
	cases := []struct {
		in, want string
		opts     *Options
	}{
		{"foo | bar", "foo-or-bar", &Options{CustomReplacements: [][2]string{{"| ", " or "}}}},
		{"10 | 20 %", "10-or-20-percent", &Options{CustomReplacements: [][2]string{{"|", " or "}, {"%", " percent "}}}},
		{"I ♥ 🦄", "i-amour-licorne", &Options{CustomReplacements: [][2]string{{"♥", " amour "}, {"🦄", " licorne "}}}},
		{"x.y.z", "xyz", &Options{CustomReplacements: [][2]string{{".", ""}}}},
		{"Zürich", "zuerich", &Options{CustomReplacements: [][2]string{{"ä", "ae"}, {"ö", "oe"}, {"ü", "ue"}, {"ß", "ss"}}}},
	}
	for _, c := range cases {
		got := slug(c.in, c.opts)
		if got != c.want {
			t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestLowercaseOption(t *testing.T) {
	cases := []struct {
		in, want string
		opts     *Options
	}{
		{"foo bar", "foo-bar", &Options{Lowercase: bp(false)}},
		{"BAR&baz", "BAR-and-baz", &Options{Lowercase: bp(false)}},
		{"Déjà Vu!", "Deja_Vu", &Options{Separator: sp("_"), Lowercase: bp(false)}},
		{"UNICORNS AND RAINBOWS!", "UNICORNS@AND@RAINBOWS", &Options{Separator: sp("@"), Lowercase: bp(false)}},
		{"[foo] [bar]", "foo.bar", &Options{Separator: sp("."), Lowercase: bp(false)}},
		{"Foo🦄", "Foo-unicorn", &Options{Lowercase: bp(false)}},
	}
	for _, c := range cases {
		got := slug(c.in, c.opts)
		if got != c.want {
			t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDecamelizeOption(t *testing.T) {
	if slug("fooBar", nil) != "foo-bar" {
		t.Error("fooBar should = foo-bar")
	}
	if slug("fooBar", &Options{Decamelize: bp(false)}) != "foobar" {
		t.Error("fooBar (no decamelize) should = foobar")
	}
}

func TestGermanUmlauts(t *testing.T) {
	got := slug("ä ö ü Ä Ö Ü ß", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "ae oe ue Ae Oe Ue ss"
	if got != want {
		t.Errorf("German umlauts = %q, want %q", got, want)
	}
}

func TestVietnamese(t *testing.T) {
	got := slug("ố Ừ Đ", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "o U D"
	if got != want {
		t.Errorf("Vietnamese = %q, want %q", got, want)
	}
}

func TestArabic(t *testing.T) {
	got := slug("ث س و", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "th s w"
	if got != want {
		t.Errorf("Arabic = %q, want %q", got, want)
	}
}

func TestPersian(t *testing.T) {
	got := slug("چ ی پ", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "ch y p"
	if got != want {
		t.Errorf("Persian = %q, want %q", got, want)
	}
}

func TestUrdu(t *testing.T) {
	got := slug("ٹ ڈ ھ", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "t d h"
	if got != want {
		t.Errorf("Urdu = %q, want %q", got, want)
	}
}

func TestPashto(t *testing.T) {
	got := slug("ګ ړ څ", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "g r c"
	if got != want {
		t.Errorf("Pashto = %q, want %q", got, want)
	}
}

func TestRussian(t *testing.T) {
	got := slug("Ж п ю", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "Zh p yu"
	if got != want {
		t.Errorf("Russian = %q, want %q", got, want)
	}
	got2 := slug("я люблю единорогов", nil)
	if got2 != "ya-lyublyu-edinorogov" {
		t.Errorf("Russian sentence = %q, want ya-lyublyu-edinorogov", got2)
	}
}

func TestRomanian(t *testing.T) {
	got := slug("ș Ț", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "s T"
	if got != want {
		t.Errorf("Romanian = %q, want %q", got, want)
	}
}

func TestTurkish(t *testing.T) {
	got := slug("İ ı Ş ş Ç ç Ğ ğ", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "I i S s C c G g"
	if got != want {
		t.Errorf("Turkish = %q, want %q", got, want)
	}
}

func TestArmenian(t *testing.T) {
	got := slug("Ե ր ե ւ ա ն", &Options{Lowercase: bp(false), Separator: sp(" ")})
	want := "Ye r ye a n"
	if got != want {
		t.Errorf("Armenian = %q, want %q", got, want)
	}
}

func TestLeadingUnderscore(t *testing.T) {
	cases := []struct{ in, want string }{
		{"_foo bar", "_foo-bar"},
		{"_foo_bar", "_foo-bar"},
		{"__foo__bar", "_foo-bar"},
		{"____-___foo__bar", "_foo-bar"},
	}
	opts := &Options{PreserveLeadingUnderscore: true}
	for _, c := range cases {
		got := slug(c.in, opts)
		if got != c.want {
			t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTrailingDash(t *testing.T) {
	opts := &Options{PreserveTrailingDash: true}
	cases := []struct{ in, want string }{
		{"foo bar-", "foo-bar-"},
		{"foo-bar--", "foo-bar-"},
		{"foo-bar -", "foo-bar-"},
		{"foo-bar - ", "foo-bar"},
		{"foo-bar ", "foo-bar"},
	}
	for _, c := range cases {
		got := slug(c.in, opts)
		if got != c.want {
			t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCounter(t *testing.T) {
	c := NewSlugifyWithCounter()
	cases := []struct {
		in   string
		opts *Options
		want string
	}{
		{"foo bar", nil, "foo-bar"},
		{"foo bar", nil, "foo-bar-2"},
	}
	for _, tc := range cases {
		got, err := c.Slugify(tc.in, tc.opts)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("counter.Slugify(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}

	c.Reset()

	cases2 := []struct {
		in   string
		opts *Options
		want string
	}{
		{"foo", nil, "foo"},
		{"foo", nil, "foo-2"},
		{"foo 1", nil, "foo-1"},
		{"foo-1", nil, "foo-1-2"},
		{"foo-1", nil, "foo-1-3"},
		{"foo", nil, "foo-3"},
		{"foo", nil, "foo-4"},
		{"foo-1", nil, "foo-1-4"},
		{"foo-2", nil, "foo-2-1"},
		{"foo-2", nil, "foo-2-2"},
		{"foo-2-1", nil, "foo-2-1-1"},
		{"foo-2-1", nil, "foo-2-1-2"},
		{"foo-11", nil, "foo-11-1"},
		{"foo-111", nil, "foo-111-1"},
		{"foo-111-1", nil, "foo-111-1-1"},
		{"fooCamelCase", &Options{Lowercase: bp(false), Decamelize: bp(false)}, "fooCamelCase"},
		{"fooCamelCase", &Options{Decamelize: bp(false)}, "foocamelcase-2"},
		{"_foo", nil, "foo-5"},
		{"_foo", &Options{PreserveLeadingUnderscore: true}, "_foo"},
		{"_foo", &Options{PreserveLeadingUnderscore: true}, "_foo-2"},
	}
	for _, tc := range cases2 {
		got, err := c.Slugify(tc.in, tc.opts)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("counter.Slugify(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}

	c2 := NewSlugifyWithCounter()
	if g, _ := c2.Slugify("foo", nil); g != "foo" {
		t.Errorf("c2 foo = %q", g)
	}
	if g, _ := c2.Slugify("foo", nil); g != "foo-2" {
		t.Errorf("c2 foo = %q", g)
	}
	if g, _ := c2.Slugify("", nil); g != "" {
		t.Errorf("c2 empty = %q", g)
	}
	if g, _ := c2.Slugify("", nil); g != "" {
		t.Errorf("c2 empty = %q", g)
	}
}

func TestPreserveCharacters(t *testing.T) {
	cases := []struct {
		in, want string
		opts     *Options
	}{
		{"foo#bar", "foo-bar", &Options{PreserveCharacters: []string{}}},
		{"foo.bar", "foo-bar", &Options{PreserveCharacters: []string{}}},
		{"foo?bar ", "foo-bar", &Options{PreserveCharacters: []string{"#"}}},
		{"foo#bar", "foo#bar", &Options{PreserveCharacters: []string{"#"}}},
		{"foo_bar#baz", "foo-bar#baz", &Options{PreserveCharacters: []string{"#"}}},
		{"foo.bar#baz-quux", "foo.bar#baz-quux", &Options{PreserveCharacters: []string{".", "#"}}},
		{"foo.bar#baz-quux", "foo.bar.baz-quux", &Options{Separator: sp("."), PreserveCharacters: []string{"-"}}},
	}
	for _, c := range cases {
		got := slug(c.in, c.opts)
		if got != c.want {
			t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}

	// separator in preserveCharacters must error
	if _, err := Slugify("foo", &Options{Separator: sp("-"), PreserveCharacters: []string{"-"}}); err == nil {
		t.Error("expected error for separator in preserveCharacters")
	}
	if _, err := Slugify("foo", &Options{Separator: sp("."), PreserveCharacters: []string{"."}}); err == nil {
		t.Error("expected error for separator in preserveCharacters")
	}
}

func TestLocaleOption(t *testing.T) {
	cases := []struct {
		in, want string
		opts     *Options
	}{
		{"Räksmörgås", "raeksmoergas", nil},
		{"Räksmörgås", "raksmorgas", &Options{Locale: "sv"}},
		{"Räksmörgås", "raeksmoergas", &Options{Locale: "de"}},
		{"Fön", "foen", &Options{Locale: "de"}},
		{"Fön", "fon", &Options{Locale: "sv"}},
		{"TEST", "test", &Options{Locale: "tr"}},
		{"TEST", "test", nil},
	}
	for _, c := range cases {
		got := slug(c.in, c.opts)
		if got != c.want {
			t.Errorf("Slugify(%q, locale) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTransliterateDisabled(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"foo bar", "foo-bar"},
		{"hello world", "hello-world"},
		{"Déjà Vu", "déjà-vu"},
		{"Räksmörgås", "räksmörgås"},
		{"你好世界", "你好世界"},
		{"مرحبا", "مرحبا"},
		{"Hello Déjà Vu", "hello-déjà-vu"},
		{"Déjà Vu", "deja-vu"}, // last entry: default (transliterate=true)
		{"foo & bar", "foo-and-bar"},
	}
	for i, c := range cases {
		var opts *Options
		if i < len(cases)-2 {
			opts = &Options{Transliterate: bp(false)}
		}
		got := slug(c.in, opts)
		if got != c.want {
			t.Errorf("[%d] Slugify(%q) = %q, want %q", i, c.in, got, c.want)
		}
	}

	// custom replacements still work when transliterate=false
	got := slug("foo & bar", &Options{
		Transliterate:      bp(false),
		CustomReplacements: [][2]string{{"&", " and "}},
	})
	if got != "foo-and-bar" {
		t.Errorf("custom repl without transliterate = %q", got)
	}

	// builtin replacements (& -> and) are disabled when transliterate=false
	got = slug("foo & bar", &Options{Transliterate: bp(false)})
	if got != "foo-bar" {
		t.Errorf("no transliterate, no custom = %q, want foo-bar", got)
	}
}
