# go-slugify

[![Go dev](https://pkg.go.dev/badge/github.com/rgglez/go-slugify)](https://pkg.go.dev/gitub.com/rgglez/go-slugify)
[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/rgglez/go-slugify/blob/master/LICENSE)
[![Build Test](https://github.com/rgglez/go-slugify/actions/workflows/build.yml/badge.svg)](https://github.com/rgglez/go-slugify/actions/workflows/build-test.yml)
[![Cross Build](https://github.com/rgglez/go-slugify/actions/workflows/cross-build.yml/badge.svg)](https://github.com/rgglez/go-slugify/actions/workflows/cross-build.yml)
[![Unit Test](https://github.com/rgglez/go-slugify/actions/workflows/unit-test.yml/badge.svg)](https://github.com/rgglez/go-slugify/actions/workflows/unit-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/go-slugify)](https://goreportcard.com/report/github.com/rgglez/go-slugify)
![GitHub stars](https://img.shields.io/github/stars/rgglez/go-slugify?style=social)
![GitHub forks](https://img.shields.io/github/forks/rgglez/go-slugify?style=social)

**go-slugify** is a Go port of [github.com/sindresorhus/slugify](github.com/sindresorhus/slugify).

It slugifies a string. Useful for URLs, filenames, and IDs.

It handles most major languages, including German (umlauts), Vietnamese, Arabic, Russian, and more.

## Install

```bash
go get github.com/rgglez/go-slugify
```

## Usage

```go
slugify.Slugify('I ♥ Dogs');
//=> 'i-love-dogs'

slugify.Slugify('  Déjà Vu!  ');
//=> 'deja-vu'

slugify.Slugify('fooBar 123 $#%');
//=> 'foo-bar-123'

slugify.Slugify('я люблю единорогов');
//=> 'ya-lyublyu-edinorogov'
```

## API

### `Slugify(s string, opts ...*Options) (string, error)`

Converts `s` into a URL-friendly slug.

- Omit `opts` (or pass `nil`) to use defaults.
- Returns `*SlugifyError` when options are invalid (for example, if `Separator` is also present in `PreserveCharacters`).

Example:

```go
slug, err := slugify.Slugify("  Déjà Vu!  ")
// slug == "deja-vu"
```

### `type Options struct`

```go
type Options struct {
	Separator                 *string
	Lowercase                 *bool
	Decamelize                *bool
	CustomReplacements        [][2]string
	PreserveLeadingUnderscore bool
	PreserveTrailingDash      bool
	PreserveCharacters        []string
	Transliterate             *bool
	Locale                    string
}
```

Fields:

- `Separator` (default: `"-"`): string used to join slug parts.
- `Lowercase` (default: `true`): lowercase output when enabled.
- `Decamelize` (default: `true`): splits camelCase and PascalCase into word boundaries.
- `CustomReplacements`: extra replacements as key/value pairs (`[][2]string`), applied before slug cleanup.
- `PreserveLeadingUnderscore` (default: `false`): keeps a leading `_` if present.
- `PreserveTrailingDash` (default: `false`): keeps a trailing `-` if present.
- `PreserveCharacters`: characters that should not be removed by slug cleanup.
- `Transliterate` (default: `true`): transliterates non-ASCII text before slug generation.
- `Locale`: locale hint for transliteration rules (for example `"sv"`).

Example:

```go
sep := "_"
lower := false
slug, err := slugify.Slugify("fooBar ÄÖÜ", &slugify.Options{
	Separator:  &sep,
	Lowercase:  &lower,
	Locale:     "sv",
	Transliterate: nil, // defaults to true
})
// slug == "foo_Bar_AOU"
```

### `BuiltinReplacements`

```go
var BuiltinReplacements = [][2]string{
	{"&", " and "},
	{"🦄", " unicorn "},
	{"♥", " love "},
}
```

Builtin replacements are applied during transliteration and can be overridden by `CustomReplacements`.

### `type SlugifyError`

```go
type SlugifyError struct{ /* ... */ }
```

Returned as `error` by `Slugify` when options are invalid.

### Counter API

#### `NewSlugifyWithCounter() *SlugifyWithCounter`

Creates a counter-based slugifier that tracks duplicates.

#### `(*SlugifyWithCounter).Slugify(s string, opts *Options) (string, error)`

Returns a unique slug:

- first occurrence: `my-slug`
- repeated occurrence: `my-slug-2`, `my-slug-3`, ...

#### `(*SlugifyWithCounter).Reset()`

Clears internal occurrence state.

Example:

```go
c := slugify.NewSlugifyWithCounter()

a, _ := c.Slugify("My post", nil) // "my-post"
b, _ := c.Slugify("My post", nil) // "my-post-2"
c.Reset()
d, _ := c.Slugify("My post", nil) // "my-post"

_, _, _ = a, b, d
```

## Reporting incompatibilities

If you find any incompatibility between the original Javascript library and
this one, please open an issue.

## License

Copyright 2026 Rodolfo González González.

Released under [MIT](https://opensource.org/license/mit). Please read the [LICENSE](LICENSE) file.
