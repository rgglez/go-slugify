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
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// charTable maps runes to their ASCII transliteration.
// Only characters that cannot be handled by NFD+strip are listed here.
var charTable = map[rune]string{
	// German umlauts (multi-char, default table)
	'Ä': "Ae", 'ä': "ae",
	'Ö': "Oe", 'ö': "oe",
	'Ü': "Ue", 'ü': "ue",
	'ß': "ss",

	// Nordic / Latin Extended
	'Æ': "Ae", 'æ': "ae",
	'Ø': "O", 'ø': "o",
	'Ð': "D", 'ð': "d",
	'Þ': "Th", 'þ': "th",

	// D with stroke (Vietnamese etc.)
	'Đ': "D", 'đ': "d",

	// Turkish dotless i (no NFD decomposition)
	'ı': "i",

	// Polish L with stroke
	'Ł': "L", 'ł': "l",

	// Cyrillic — Russian
	'А': "A", 'а': "a",
	'Б': "B", 'б': "b",
	'В': "V", 'в': "v",
	'Г': "G", 'г': "g",
	'Д': "D", 'д': "d",
	'Е': "E", 'е': "e",
	'Ё': "Yo", 'ё': "yo",
	'Ж': "Zh", 'ж': "zh",
	'З': "Z", 'з': "z",
	'И': "I", 'и': "i",
	'Й': "J", 'й': "j",
	'К': "K", 'к': "k",
	'Л': "L", 'л': "l",
	'М': "M", 'м': "m",
	'Н': "N", 'н': "n",
	'О': "O", 'о': "o",
	'П': "P", 'п': "p",
	'Р': "R", 'р': "r",
	'С': "S", 'с': "s",
	'Т': "T", 'т': "t",
	'У': "U", 'у': "u",
	'Ф': "F", 'ф': "f",
	'Х': "Kh", 'х': "kh",
	'Ц': "Ts", 'ц': "ts",
	'Ч': "Ch", 'ч': "ch",
	'Ш': "Sh", 'ш': "sh",
	'Щ': "Shch", 'щ': "shch",
	'Ъ': "", 'ъ': "",
	'Ы': "Y", 'ы': "y",
	'Ь': "", 'ь': "",
	'Э': "E", 'э': "e",
	'Ю': "Yu", 'ю': "yu",
	'Я': "Ya", 'я': "ya",

	// Ukrainian Cyrillic
	'І': "I", 'і': "i",
	'Ї': "Yi", 'ї': "yi",
	'Є': "Ye", 'є': "ye",
	'Ґ': "G", 'ґ': "g",

	// Arabic
	'ء': "",   // ء hamza
	'آ': "a",  // آ alef madda
	'أ': "a",  // أ alef hamza above
	'ؤ': "w",  // ؤ waw hamza
	'إ': "i",  // إ alef hamza below
	'ئ': "y",  // ئ ya hamza
	'ا': "a",  // ا alef
	'ب': "b",  // ب ba
	'ة': "t",  // ة ta marbuta
	'ت': "t",  // ت ta
	'ث': "th", // ث tha
	'ج': "j",  // ج jim
	'ح': "h",  // ح ha
	'خ': "kh", // خ kha
	'د': "d",  // د dal
	'ذ': "dh", // ذ dhal
	'ر': "r",  // ر ra
	'ز': "z",  // ز zay
	'س': "s",  // س sin
	'ش': "sh", // ش shin
	'ص': "s",  // ص sad
	'ض': "d",  // ض dad
	'ط': "t",  // ط ta
	'ظ': "z",  // ظ za
	'ع': "",   // ع ayn
	'غ': "gh", // غ ghayn
	'ف': "f",  // ف fa
	'ق': "q",  // ق qaf
	'ك': "k",  // ك kaf
	'ل': "l",  // ل lam
	'م': "m",  // م mim
	'ن': "n",  // ن nun
	'ه': "h",  // ه ha
	'و': "w",  // و waw
	'ى': "a",  // ى alef maqsura
	'ي': "y",  // ي ya

	// Arabic vowel marks (diacritics — strip)
	'ً': "", 'ٌ': "", 'ٍ': "",
	'َ': "", 'ُ': "", 'ِ': "",
	'ّ': "", 'ْ': "",

	// Persian extras
	'پ': "p",  // پ pa
	'چ': "ch", // چ che
	'ژ': "zh", // ژ zhe
	'ک': "k",  // ک kaf
	'گ': "g",  // گ gaf
	'ی': "y",  // ی ya (Persian)
	'ھ': "h",  // ھ he (Urdu)

	// Urdu extras
	'ٹ': "t", // ٹ tte
	'ڈ': "d", // ڈ ddal
	'ڑ': "r", // ڑ rra
	'ں': "n", // ں nun ghunna
	'ہ': "h", // ہ he

	// Pashto extras
	'څ': "c",  // څ
	'ړ': "r",  // ړ
	'ګ': "g",  // ګ
	'ځ': "z",  // ځ
	'ډ': "d",  // ډ
	'ږ': "zh", // ږ
	'ټ': "t",  // ټ
	'ۍ': "i",  // ۍ

	// Armenian capital letters (U+0531–U+0556)
	'Ա': "A",  // Ա
	'Բ': "B",  // Բ
	'Գ': "G",  // Գ
	'Դ': "D",  // Դ
	'Ե': "Ye", // Ե
	'Զ': "Z",  // Զ
	'Է': "E",  // Է
	'Ը': "Ye", // Ը
	'Թ': "T",  // Թ
	'Ժ': "Zh", // Ժ
	'Ի': "I",  // Ի
	'Լ': "L",  // Լ
	'Խ': "Kh", // Խ
	'Ծ': "Ts", // Ծ
	'Կ': "K",  // Կ
	'Հ': "H",  // Հ
	'Ձ': "Dz", // Ձ
	'Ղ': "Gh", // Ղ
	'Ճ': "Ch", // Ճ
	'Մ': "M",  // Մ
	'Յ': "Y",  // Յ
	'Ն': "N",  // Ն
	'Շ': "Sh", // Շ
	'Ո': "Vo", // Ո
	'Չ': "Ch", // Չ
	'Պ': "P",  // Պ
	'Ջ': "J",  // Ջ
	'Ռ': "R",  // Ռ
	'Ս': "S",  // Ս
	'Վ': "V",  // Վ
	'Տ': "T",  // Տ
	'Ր': "R",  // Ր
	'Ց': "Ts", // Ց
	'Ւ': "W",  // Փ
	'Փ': "P",  // Փ — reuse U+0553
	'Ք': "K",  // Ք

	// Armenian small letters (U+0561–U+0587)
	'ա': "a",   // ա
	'բ': "b",   // բ
	'գ': "g",   // գ
	'դ': "d",   // դ
	'ե': "ye",  // ե
	'զ': "z",   // զ
	'է': "e",   // է
	'ը': "ye",  // ը
	'թ': "t",   // թ
	'ժ': "zh",  // ժ
	'ի': "i",   // ի
	'լ': "l",   // լ
	'խ': "kh",  // խ
	'ծ': "ts",  // ծ
	'կ': "k",   // կ
	'հ': "h",   // հ
	'ձ': "dz",  // ձ
	'ղ': "gh",  // ղ
	'ճ': "ch",  // ճ
	'մ': "m",   // մ
	'յ': "y",   // յ
	'ն': "n",   // ն
	'շ': "sh",  // շ
	'ո': "vo",  // ո
	'չ': "ch",  // չ
	'պ': "p",   // պ
	'ջ': "j",   // ջ
	'ռ': "r",   // ռ
	'ս': "s",   // ս
	'վ': "v",   // վ
	'տ': "t",   // տ
	'ր': "r",   // ր
	'ց': "ts",  // ց
	'ւ': "",    // ւ (silent in modern Armenian)
	'փ': "p",   // փ
	'ք': "k",   // ք
	'և': "yev", // և
}

// localeOverrides maps locale codes to rune-level overrides that replace charTable entries.
var localeOverrides = map[string]map[rune]string{
	"sv": {
		'Ä': "A", 'ä': "a",
		'Ö': "O", 'ö': "o",
		'Å': "A", 'å': "a",
		'Ü': "U", 'ü': "u",
	},
	"de": {}, // same as default
}

// mergeReplacements merges builtins with user replacements; user overrides builtins for same key.
func mergeReplacements(builtins, user [][2]string) [][2]string {
	result := make([][2]string, len(builtins))
	copy(result, builtins)
	keyIdx := make(map[string]int, len(builtins))
	for i, r := range result {
		keyIdx[r[0]] = i
	}
	for _, r := range user {
		if i, ok := keyIdx[r[0]]; ok {
			result[i] = r
		} else {
			keyIdx[r[0]] = len(result)
			result = append(result, r)
		}
	}
	return result
}

// transliterateStr converts Unicode string to ASCII using the char table and NFD fallback.
func transliterateStr(s, locale string, replacements [][2]string) string {
	for _, r := range replacements {
		s = strings.ReplaceAll(s, r[0], r[1])
	}

	locOvr := localeOverrides[locale]

	var buf strings.Builder
	buf.Grow(len(s))
	for _, r := range s {
		if locOvr != nil {
			if repl, ok := locOvr[r]; ok {
				buf.WriteString(repl)
				continue
			}
		}
		if repl, ok := charTable[r]; ok {
			buf.WriteString(repl)
			continue
		}
		// NFD decomposition: strip combining (non-spacing) marks
		nfd := norm.NFD.String(string(r))
		for _, nr := range nfd {
			if unicode.Is(unicode.Mn, nr) {
				continue
			}
			buf.WriteRune(nr)
		}
	}
	return buf.String()
}
