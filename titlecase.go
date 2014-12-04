
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
)

const (
 Generic 	= 0
 English 	= 1
 French  	= 2
 German  	= 3
 Italian 	= 4
 Spanish 	= 5
 Portuguese = 6
 Dutch 		= 7
 Latin 		= 8
)

type honor struct {
 binsearch.Key_runes
 format [][]rune
}

var romanExceptions, makecaps, englishSmall, frenchSmall, germanSmall, italianSmall, spanishSmall, dutchSmall, portugueseSmall, latinSmall binsearch.Key_runes

func init() {

	// Initate exceptions for Roman numerals
	romanExceptions.Key = [][]rune {
	 []rune("ci"), []rune("cid"), []rune("cill"), []rune("civic"), []rune("civil"), []rune("clim"), []rune("cm"), []rune("di"), []rune("did"), []rune("didi"), []rune("dill"), []rune("dilli"),
	 []rune("dim"), []rune("divi"), []rune("dividivi"), []rune("dix"), []rune("dixi"), []rune("dixil"), []rune("dm"), []rune("id"), []rune("ill"), []rune("im"), []rune("imid"), []rune("imidic"),
	 []rune("immix"), []rune("ld"), []rune("li"), []rune("lid"), []rune("lil"), []rune("lili"), []rune("lill"), []rune("lilli"), []rune("lim"), []rune("liv"), []rune("livi"), []rune("livid"),
	 []rune("livvi"), []rune("lm"), []rune("lviv"), []rune("mic"), []rune("mid"), []rune("midi"), []rune("mil"), []rune("mild"), []rune("mill"), []rune("milli"), []rune("mim"),
	 []rune("mimi"), []rune("mimic"), []rune("mix"), []rune("mv"), []rune("vi"), []rune("vic"), []rune("vici"), []rune("vid"), []rune("vild"), []rune("vill"), []rune("villi"), []rune("vim"),
	 []rune("viv"), []rune("vivi"), []rune("vivid"), []rune("vivl"), //[]rune("md"),
	}
	romanExceptions.Build()
	
	// Initate exceptions for ALLCAPS
	makecaps.Key = [][]rune {
	 []rune("abc"), []rune("usa"), []rune("ussr"), []rune("usaf"), []rune("uscg"), []rune("usmc"), []rune("usn"), []rune("ymca"), []rune("raf"), []rune("uk"),
	}
	makecaps.Build()
	
	// Initate exceptions for honor
	honor.Key = [][]rune {
	 []rune("m.d"), []rune("ph.d"),
	}
	honor.format = [][]rune {
	 []rune("M.D"), []rune("Ph.D"),
	}
	temp := make([]int, len(honor.format))
	newindexes := honor.Build()
	for indx_new, indx_old := range newindexes {
		temp[indx_new] = honor.format[indx_old]
	}
	honor.format = temp
	
	// Initiate exceptions for English small words
	englishSmall.Key = [][]rune {
	 []rune("a"), []rune("an"), []rune("and"), []rune("as"), []rune("at"), []rune("but"), []rune("by"), []rune("for"), []rune("if"), []rune("in"), []rune("of"), []rune("on"), []rune("or"), []rune("the"), []rune("to"),
	}
	englishSmall.Build()
	
	// Initiate exceptions for French small words: d', l'
	frenchSmall.Key = [][]rune {
	 []rune("à"), []rune("au"), []rune("aux"), []rune("ce"), []rune("cette"), []rune("dans"), []rune("de"), []rune("des"), []rune("du"), []rune("en"), []rune("la"), []rune("le"), []rune("les"), []rune("ou"), []rune("par"),
	 []rune("pour"), []rune("sur"), []rune("un"),[]rune("une"),
	}
	frenchSmall.Build()
	
	// Initiate exceptions for German small words
	germanSmall.Key = [][]rune {
	 []rune("als"), []rune("am"), []rune("an"), []rune("auf"), []rune("aus"), []rune("bei"), []rune("bis"), []rune("das"), []rune("dem"), []rune("den"), []rune("der"), []rune("des"), []rune("die"), []rune("ein"), []rune("eine"),
	 []rune("für"), []rune("im"), []rune("in"), []rune("ins"), []rune("mit"), []rune("nach"), []rune("oder"), []rune("og"), []rune("und"), []rune("van"), []rune("vom"), []rune("von"), []rune("wie"), []rune("zu"), []rune("zum"), []rune("zur"),
	}
	germanSmall.Build()
	
	// Initiate exceptions for Italian small words
	italianSmall.Key = [][]rune {
	 []rune("a"), []rune("al"), []rune("con"), []rune("da"), []rune("dai"), []rune("dal"), []rune("dei"), []rune("del"), []rune("della"), []rune("di"), []rune("e"), []rune("ed"), []rune("i"), []rune("il"), []rune("in"),
	 []rune("la"), []rune("le"), []rune("lo"), []rune("nella"), []rune("o"), []rune("per"), []rune("se"), []rune("su"), []rune("un"), []rune("una"), []rune("uno"),
	}
	italianSmall.Build()
	
	// Initiate exceptions for Italian small words
	portugueseSmall.Key = [][]rune {
	 []rune("à"), []rune("às"), []rune("ao"), []rune("da"), []rune("das"), []rune("de"), []rune("do"), []rune("e"), []rune("em"), []rune("na"), []rune("no"), []rune("o"), []rune("para"), []rune("pelo"), []rune("pelos"),
	 []rune("por"), []rune("se"), []rune("um"), []rune("uma"),
	}
	portugueseSmall.Build()
	
	// Initiate exceptions for Italian small words
	spanishSmall.Key = [][]rune {
	 []rune("a"), []rune("al"), []rune("de"), []rune("del"), []rune("e"), []rune("é"), []rune("el"), []rune("en"), []rune("la"), []rune("las"), []rune("los"), []rune("o"), []rune("ó"), []rune("para"), []rune("por"),
	 []rune("si"), []rune("un"), []rune("una"), []rune("y"),
	}
	spanishSmall.Build()
	
}

// Structs
type wordStruct struct {
 content []rune
 isStart bool
 isEnd bool
 isHonor bool
 contraction uint8
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
	var i, i2, i3 int
	var i4 uint8
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
	var isHonor bool
	var contraction uint8
	i4 = 0
	content := make([]rune, i2 - i)
	noise := make([]uint8, 0)
	for i3=i; i3<i2; i3++ {
		rn = w[i3]
		switch rn {
			case '.', ',', ';', ':', '!', '?', '&': // if any of these occur in the middle of a word (surrounded by letters) then split into two words
				noise = append(noise, i4)
			case 39, '’':
				if i4 < (i2 - i) - 2 {
					contraction = i4
				}
		}
		content[i4] = unicode.ToLower(w[i3])
		i4++
	}
	// Check if any noise occurred
	if len(noise) > 0 {
		if id, ok := honors.Find(content); ok { // if it's an honor then save the way it should be displayed
			isHonor = true
			content = honors.format[id]
		} else { // if it's not then split the word
			backup := r.runes
			saver := make([]rune, len(content))
			copy(saver, content)
			r.runes = saver[0:noise[0]]
			r.len = len(r.runes)
			words = r.add(words, 1)
			r.runes = saver[noise[0]+1:]
			r.len = len(r.runes)
			words = r.add(words, 1)
			r.len = 0
			r.runes = backup
			return words
		}
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
	words = append(words, wordStruct{content, false, isEnd, isHonor, contraction, spaceType, puncBefore, puncAfter})
	return words
}

func isRoman(word []rune) bool {
	var r rune
	for _, r = range word {
		switch r {
			case 'i', 'v', 'x', 'm', 'c', 'd', 'l':
				continue
			default:
				return false
		}
	}
	if _, ok := romanExceptions.Find(word); ok {
		return false
	} else {
		return true
	}
}

func isContraction(wb []rune) bool {
	switch len(wb) {
		case 1:
			switch wb[0] {
				case 'b': fallthrough
				case 's': fallthrough
				case 'd': fallthrough
				case 'n': fallthrough
				case 'l': fallthrough
				case 'm': fallthrough
				case 't': fallthrough
				case 'v': fallthrough
				case 'j': return true
			}
		case 2:
			if (wb[0]=='u' && wb[1]=='n') || (wb[0]=='q' && wb[1]=='u') || (wb[0]=='g' && wb[1]=='l') {
				return true
			}
		case 3:
			if (wb[0]=='a' && wb[1]=='l' && wb[2]=='l') || (wb[0]=='a' && wb[1]=='g' && wb[2]=='l') {
				return true
			}
		case 4:
			if (wb[3]!='l') {
				return false
			}
			if (wb[2]!='l' && wb[2]!='g') {
				return false
			}
			switch wb[1] {
				case 'a': fallthrough
				case 'e': fallthrough
				case 'u': fallthrough
				case 'o':
					switch wb[0] {
						case 'd': fallthrough
						case 'n': fallthrough
						case 's': fallthrough
						case 'c': fallthrough
						case 'p': return true
					}
				
			}
	}
	return false
}

func upperRune(word []rune, which int) {
	if which == -1 {
		for i, r := range word {
			if r == 39 || r == '’' { // stop uppercasing when an apostrophe is reached
				return
			}
			word[i] = unicode.ToTitle(r)
		}
		return
	}
	word[which] = unicode.ToTitle(word[which])
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
		case French: small = frenchSmall
		case German: small = germanSmall
		case Italian: small = italianSmall
		case Spanish: small = spanishSmall
		case Portuguese: small = portugueseSmall
	}
	
	// Preprocessing
	str = html.UnescapeString(str)
	b := []byte(str)
	b = bytes.Replace(b, []byte("--"), []byte("—"), -1) // Correct hyphens to em dashes
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
					continue
				}
			case '/':
				if word.len > 0 {
					words = word.add(words, 3)
					continue
				}
			case '[', '{': r = '('
			case ']', '}': r = ')'
		}
		word.write(r)
	}
	if word.len > 0 {
		words = word.add(words, 4)
	}
	word = nil
	
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
		
		if ws.isHonor {
			continue
		}
		
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
		
		if language == English {
			// Special for English: repair grammatical error on a -> an
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
		} else {
			// Check for contractions
			if ws.contraction > 0 && language != German {
				if isContraction(content[0:ws.contraction]) {
					if ws.isStart {
						upperRune(content, 0)
					}
					content = ws.content[ws.contraction+1:]
				}
			}
		}
		
		if _, ok = makecaps.Find(content); ok {
			upperRune(content, -1)
			continue
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

