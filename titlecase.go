
/*

TITLECASE

* Supports multiple languages: English, French, German, Italian, Spanish & Generic
* Written for speed - no regular expressions, minimal looping
* Decodes all HTML entities
* Supports roman numerals
* Repairs grammatical errors
* Fully UTF8 compliant

This is a comprehensive titlecase formatter made by Forgotten Books for formatting book titles.
No authority "standards" are adhered to because none of them cover all eventualities. Instead the rules were taken from first hand experience in the problem of titlecasing.

*/

package titlecase

import (
 "html"
 "bytes"
 "unicode"
 "unicode/utf8"
 "github.com/AlasdairF/BinSearch"
 "fmt"
)

const (
 Generic = 0
 English = 1
 French  = 2
 German  = 3
 Italian = 4
 Spanish = 5
)

var romanExceptions binsearch.Key_runes
var englishSmall binsearch.Key_runes

func init() {

	// Initate exceptions for Roman numerals
	romanExceptions.Key = [][]rune {
	 []rune("ci"), []rune("cid"), []rune("cill"), []rune("civic"), []rune("civil"), []rune("clim"), []rune("cm"), []rune("di"), []rune("did"), []rune("didi"), []rune("dill"), []rune("dilli"),
	 []rune("dim"), []rune("divi"), []rune("dividivi"), []rune("dix"), []rune("dixi"), []rune("dixil"), []rune("dm"), []rune("id"), []rune("ill"), []rune("im"), []rune("imid"), []rune("imidic"),
	 []rune("immix"), []rune("ld"), []rune("li"), []rune("lid"), []rune("lil"), []rune("lili"), []rune("lill"), []rune("lilli"), []rune("lim"), []rune("liv"), []rune("livi"), []rune("livid"),
	 []rune("livvi"), []rune("lm"), []rune("lviv"), []rune("md"), []rune("mic"), []rune("mid"), []rune("midi"), []rune("mil"), []rune("mild"), []rune("mill"), []rune("milli"), []rune("mim"),
	 []rune("mimi"), []rune("mimic"), []rune("mix"), []rune("mv"), []rune("vi"), []rune("vic"), []rune("vici"), []rune("vid"), []rune("vild"), []rune("vill"), []rune("villi"), []rune("vim"),
	 []rune("viv"), []rune("vivi"), []rune("vivid"), []rune("vivl"),
	}
	romanExceptions.Build()
	
	// Initiate exceptions for English small words
	englishSmall.Key = [][]rune {
	 []rune("a"), []rune("an"), []rune("and"), []rune("as"), []rune("at"), []rune("but"), []rune("by"), []rune("for"), []rune("if"), []rune("in"), []rune("of"), []rune("on"), []rune("or"), []rune("the"), []rune("to"),
	}
	englishSmall.Build()
	
}

// Structs
type wordStruct struct {
 content []rune
 isStart bool
 isEnd bool
 spaceAfter uint8 // 0=nothing, 1=space, 2=hypen, 3=slash, 4=end
 puncBefore []rune
 puncAfter []rune
}

type runebuf struct {
 runes []rune
 len int
}
func (r *runebuf) write(rn rune) {
	r.runes[r.len] = rn
	r.len++
}
func newRuneBuf() *runebuf {
	r := new(runebuf)
	r.runes = make([]rune, 256)
	return r
}
func (r *runebuf) add(words []wordStruct, spaceType uint8) []wordStruct {
	l := r.len
	w := r.runes[0:l]
	puncBefore := make([]rune, 0)
	var i, i2, i3, i4 int
	// Get punctuation before word
	for i=0; i<l; i++ {
		if unicode.IsPunct(w[i]) {
			puncBefore = append(puncBefore, w[i])
		} else {
			break
		}
	}
	// Get punctuation after the word
	for i2=l-1; i2>i; i2-- {
		if !unicode.IsPunct(w[i2]) {
			break
		}
	}
	i2++
	puncAfter := make([]rune, l - i2)
	for i3=i2; i3<l; i3++ {
		puncAfter[i4] = w[i3]
		i4++
	}
	// Get word
	var rn rune
	i4 = 0
	content := make([]rune, (i2-i)+1)
	for i3=i; i3<i2; i3++ {
		rn = w[i3]
		switch rn {
			case '.', ',', ';', ':', '!', '?', '&': // if any of these occur in the middle of a word (surrounded by letters) then split into two words
				r.runes = w[i:i3+1]
				r.len = (i3 - i)+1
				words = r.add(words, 1)
				r.runes = w[i3+1:]
				r.len = len(r.runes)
				words = r.add(words, 1)
				r.len = 0
				return words
		}
		content[i4] = unicode.ToLower(w[i3])
		i4++
	}
	// Determine if this is an ending, that means any punctuation except an aprostrophe
	var isEnd bool
	if len(puncAfter) > 0 {
		isEnd = true
		if len(puncAfter) == 1 {
			switch puncAfter[0] {
				case '‘', '’', '`', 39: isEnd = false
			}
		}
	}
	// Reset buffer
	r.len = 0
	words = append(words, wordStruct{content, false, isEnd, spaceType, puncBefore, puncAfter})
	return words
}

func isRoman(r []rune) bool {
	for _, rn := range r {
		switch rn {
			case 'm', 'c', 'd', 'x', 'l', 'v', 'i': continue
			default: return false
		}
	}
	if _, ok := romanExceptions.Find(r); ok {
		return false
	} else {
		return true
	}
}

func upperRune(r []rune, which int) {
	if which == -1 {
		for i, rn := range r {
			r[i] = unicode.ToTitle(rn)
		}
		return
	}
	r[which] = unicode.ToTitle(r[which])
}

// Removes an individual byte from a slice of bytes
func removeBytes(s []byte, a byte, b byte) []byte {
	var on int
	for _, v := range s {
		if v != a && v != b {
			s[on] = v
			on++
		}
	}
	return s[0:on]
}

func Format(str string, language uint8) string {

	var small binsearch.Key_runes
	switch language {
		case English: small = englishSmall
		//case French: small = frenchSmall
		//case German: small = germanSmall
		//case Italian: small = italianSmall
		//case Spanish: small = spanishSmall
	}
	
	// Preprocessing
	str = html.UnescapeString(str)
	b := []byte(str)
	b = bytes.Replace(b, []byte("—"), []byte(" — "), -1) // Separate out em dashes
	b = bytes.Replace(b, []byte(" - "), []byte(" — "), -1) // Correct hyphens to em dashes
	b = bytes.Replace(b, []byte("[microform]"), []byte(""), -1)
	b = bytes.Trim(b, ` ;:.,`)
	if b[0] == '(' {
		b = removeBytes(b, '(', ')')
	}
	if b[0] == '[' {
		b = removeBytes(b, '[', ']')
	}

	n := len(b)
	if n == 0 {
		return ``
	}
	
	// Load all into struct
	var r rune
	var i, w int
	words := make([]wordStruct, 0, 4)
	word := newRuneBuf()
    for i=0; i<n; i+=w {
        r, w = utf8.DecodeRune(b[i:])
		// Parse spacers
		if r <= 32 { // space
			if word.len > 0 {
				words = word.add(words, 1)
			}
			continue
		}
		switch r {
			case '-':
				if word.len > 0 {
					words = word.add(words, 2)
				}
				continue
			case '/':
				if word.len > 0 {
					words = word.add(words, 3)
				}
				continue
			case '[', '{': r = '('
			case ']', '}': r = ')'
		}
		fmt.Print(word.len, ` `)
		word.write(r)
	}
	if word.len > 0 {
		words = word.add(words, 4)
	}
	word = nil
	fmt.Println()
	
	// Determine isStart from isEnd
	l := len(words)
	if l == 0 {
		return ``
	}
	words[0].isStart = true
	for i=1; i<l; i++ {
		if words[i-1].isEnd {
			words[i].isStart = true
		}
	}
	words[l-1].isEnd = true
	
	// Loop through all and apply rules
	var ws *wordStruct
	var content []rune
	var ok bool
	var ln int
	for i=0; i<l; i++ {
		ws = &words[i]
		content = ws.content
		ln = len(content)
		
		// Uppercase roman numerals
		if isRoman(content) {
			upperRune(content, -1) // -1 means uppercase all
			continue
		}
		
		// Check for McStuff
		if ln > 3 {
			if content[0] == 'm' && content[1] == 'c' {
				upperRune(content, 0)
				upperRune(content, 2)
				continue
			}
		}
		
		// Repair grammatical error on a -> an
		if language == English {
			if ln == 1 {
				if content[0] == 'a' {
					tmp := words[i+1].content
					if len(tmp) > 1 {
						switch tmp[0] {
							case 'a', 'e', 'i', 'o', 'u':
								ws.content = []rune("an")
								content = ws.content
								ln = 2
						}
					}
				}
			}
		}
		
		// Beginning and ending words need to be capitalized regardless of what they are
		if ws.isStart || ws.isEnd {
			upperRune(content, 0)
			continue
		}
		
		// Check for small words to keep lowercase, using binary search
		if _, ok = small.Find(content); ok {
			// Exception if it's 1 letter with following punctuation or the next word is also 1 letter
			if ln > 1 {
				continue
			}
			if len(words[i+1].content) > 1 {
				continue
			}
		}
		
		// Uppercase the first rune if none of the previous rules applied
		upperRune(content, 0)
	}
	
	// Rebuild byte stream from words
	var buf bytes.Buffer
	for i=0; i<l; i++ {
		ws = &words[i]
		for _, r = range ws.puncBefore {
			buf.WriteRune(r)
		}
		for _, r = range ws.content {
			buf.WriteRune(r)
		}
		for _, r = range ws.puncAfter {
			buf.WriteRune(r)
		}
		switch ws.spaceAfter {
			case 1: buf.WriteByte(' ')
			case 2: buf.WriteByte('-')
			case 3: buf.WriteByte('/')
		}
	}
	
	return buf.String()
}

