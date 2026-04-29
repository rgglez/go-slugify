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

// Package slugify is a Go port of sindresorhus/slugify.
// It converts strings into URL-friendly slugs.
package slugify

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// BuiltinReplacements are applied during transliteration and can be overridden by CustomReplacements.
var BuiltinReplacements = [][2]string{
	{"&", " and "},
	{"🦄", " unicorn "},
	{"♥", " love "},
}

// Options controls slug generation behaviour.
// Nil pointer fields fall back to their documented defaults.
type Options struct {
	Separator                 *string    // default: "-"
	Lowercase                 *bool      // default: true
	Decamelize                *bool      // default: true
	CustomReplacements        [][2]string
	PreserveLeadingUnderscore bool
	PreserveTrailingDash      bool
	PreserveCharacters        []string
	Transliterate             *bool   // default: true
	Locale                    string
}

// resolved holds all options with defaults applied.
type resolved struct {
	separator                 string
	lowercase                 bool
	decamelize                bool
	customReplacements        [][2]string
	preserveLeadingUnderscore bool
	preserveTrailingDash      bool
	preserveCharacters        []string
	transliterate             bool
	locale                    string
}

func resolveOpts(opts *Options) resolved {
	r := resolved{
		separator:     "-",
		lowercase:     true,
		decamelize:    true,
		transliterate: true,
	}
	if opts == nil {
		return r
	}
	if opts.Separator != nil {
		r.separator = *opts.Separator
	}
	if opts.Lowercase != nil {
		r.lowercase = *opts.Lowercase
	}
	if opts.Decamelize != nil {
		r.decamelize = *opts.Decamelize
	}
	if opts.Transliterate != nil {
		r.transliterate = *opts.Transliterate
	}
	r.customReplacements = opts.CustomReplacements
	r.preserveLeadingUnderscore = opts.PreserveLeadingUnderscore
	r.preserveTrailingDash = opts.PreserveTrailingDash
	r.preserveCharacters = opts.PreserveCharacters
	r.locale = opts.Locale
	return r
}

var (
	decamRe1 = regexp.MustCompile(`([A-Z]{2,})(\d+)`)
	decamRe2 = regexp.MustCompile(`([a-z\d]+)([A-Z]{2,})`)
	decamRe3 = regexp.MustCompile(`([a-z\d])([A-Z])`)
	decamRe4 = regexp.MustCompile(`([A-Z]+)([A-Z][a-rt-z\d]+)`)

	contractionRe = regexp.MustCompile(`([a-zA-Z0-9]+)['\x{2019}]([ts])(\s|$)`)
)

func decamelize(s string) string {
	s = decamRe1.ReplaceAllString(s, "$1 $2")
	s = decamRe2.ReplaceAllString(s, "$1 $2")
	s = decamRe3.ReplaceAllString(s, "$1 $2")
	s = decamRe4.ReplaceAllString(s, "$1 $2")
	return s
}

func removeMootSeparators(s, sep string) string {
	if sep == "" {
		return s
	}

	// Collapse repeated separators.
	doubleSep := sep + sep
	for strings.Contains(s, doubleSep) {
		s = strings.ReplaceAll(s, doubleSep, sep)
	}

	// Trim separators at the edges.
	for strings.HasPrefix(s, sep) {
		s = strings.TrimPrefix(s, sep)
	}
	for strings.HasSuffix(s, sep) {
		s = strings.TrimSuffix(s, sep)
	}

	return s
}

func trimTrailingDashDigitGroups(s string) string {
	i := len(s)
	for i > 0 {
		j := i
		for j > 0 {
			r, size := utf8.DecodeLastRuneInString(s[:j])
			if !unicode.IsDigit(r) {
				break
			}
			j -= size
		}

		if j == i {
			break
		}

		k := j
		r, size := utf8.DecodeLastRuneInString(s[:k])
		if r != '-' {
			break
		}
		k -= size
		i = k
	}
	return s[:i]
}

func buildSlugPattern(opts resolved) *regexp.Regexp {
	neg := `a-z0-9`
	if !opts.lowercase {
		neg += `A-Z`
	}
	if !opts.transliterate {
		neg += `\p{L}\p{N}`
	}
	for _, ch := range opts.preserveCharacters {
		neg += regexp.QuoteMeta(ch)
	}
	return regexp.MustCompile(`[^` + neg + `]+`)
}

// Slugify converts s into a URL slug using optional options (defaults when omitted or nil).
// Returns an error if a preserved character conflicts with the separator.
func Slugify(s string, opts ...*Options) (string, error) {
	var opt *Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	o := resolveOpts(opt)

	// Validate: preserveCharacters must not contain the separator
	for _, ch := range o.preserveCharacters {
		if ch == o.separator {
			return "", &SlugifyError{
				msg: "the separator character `" + o.separator +
					"` cannot be included in preserveCharacters",
			}
		}
	}

	shouldPrependUnderscore := o.preserveLeadingUnderscore && strings.HasPrefix(s, "_")
	shouldAppendDash := o.preserveTrailingDash && strings.HasSuffix(s, "-")

	if o.transliterate {
		merged := mergeReplacements(BuiltinReplacements, o.customReplacements)
		s = transliterateStr(s, o.locale, merged)
	} else if len(o.customReplacements) > 0 {
		for _, r := range o.customReplacements {
			s = strings.ReplaceAll(s, r[0], r[1])
		}
	}

	if o.decamelize {
		s = decamelize(s)
	}

	if o.lowercase {
		s = strings.ToLower(s)
	}

	// Remove contractions / possessives: "Conway's" → "Conways"
	s = contractionRe.ReplaceAllString(s, "$1$2$3")

	slugPattern := buildSlugPattern(o)
	s = slugPattern.ReplaceAllLiteralString(s, o.separator)
	s = strings.ReplaceAll(s, `\`, "")

	if o.separator != "" {
		s = removeMootSeparators(s, o.separator)
	}

	if shouldPrependUnderscore {
		s = "_" + s
	}
	if shouldAppendDash {
		s = s + "-"
	}

	return s, nil
}

// SlugifyError is returned when slug options are invalid.
type SlugifyError struct{ msg string }

func (e *SlugifyError) Error() string { return e.msg }

// SlugifyWithCounter generates unique slugs by appending a counter for duplicates.
type SlugifyWithCounter struct {
	occurrences map[string]int
}

// NewSlugifyWithCounter returns a new counter-based slugifier.
func NewSlugifyWithCounter() *SlugifyWithCounter {
	return &SlugifyWithCounter{occurrences: make(map[string]int)}
}

// Reset clears the occurrence counter.
func (c *SlugifyWithCounter) Reset() {
	c.occurrences = make(map[string]int)
}

// Slugify returns a unique slug, appending -N for repeated inputs.
func (c *SlugifyWithCounter) Slugify(s string, opts *Options) (string, error) {
	slug, err := Slugify(s, opts)
	if err != nil {
		return "", err
	}
	if slug == "" {
		return "", nil
	}

	lower := strings.ToLower(slug)
	numberless := trimTrailingDashDigitGroups(lower)

	nlCount := c.occurrences[numberless]
	counter, exists := c.occurrences[lower]
	if exists {
		c.occurrences[lower] = counter + 1
	} else {
		c.occurrences[lower] = 1
	}
	newCounter := c.occurrences[lower]

	if newCounter >= 2 || nlCount > 2 {
		slug = slug + "-" + strconv.Itoa(newCounter)
	}

	return slug, nil
}
