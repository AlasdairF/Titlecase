
/*

TITLECASE

This is a production-quality package made for cleaning and formatting book titles, but it can be used for titlecasing anything.

* Supports multiple languages: English, French, German, Italian, Spanish, Portuguese & Generic
* Supports contractions
* Supports initials
* Supports academic honors (M.D., Ph.D, etc.)
* Supports common abbreviations (USA, USSR, YMCA, etc.)
* Supports Roman numerals, without mistaking words for roman numerals
* Supports hyphenation and slashes
* Repairs grammatical errors in English
* Redetermines whitespace
* Converts or strips inappropriate punctuation
* Decodes all HTML entities
* Fully UTF8 compliant
* Written for speed and efficiency - no regular expressions, minimal looping

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
 Language_Generic 		= 0
 Language_English 		= 1
 Language_French  		= 2
 Language_German  		= 3
 Language_Italian 		= 4
 Language_Spanish 		= 5
 Language_Portuguese	= 6
)

type honorStruct struct {
 binsearch.Key_runes
 format [][]rune
}

var romanExceptions, makecaps, englishSmall, frenchSmall, germanSmall, italianSmall, spanishSmall, portugueseSmall, titlesabv, titles, multilast binsearch.Key_runes
var honor honorStruct

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
	
	// Initate exceptions for titlesabv
	titlesabv.Key = [][]rune {
	 []rune("mr"), []rune("ms"), []rune("miss"), []rune("mrs"), []rune("dr"), []rune("prof"), []rune("rev"), []rune("esq"), []rune("hon"), []rune("jr"), []rune("messrs"), []rune("mmes"), []rune("msgr"), []rune("rt"),
	 []rune("sr"), []rune("st"), []rune("lt"), []rune("col"), []rune("gen"), []rune("pseud"), []rune("maj"), []rune("brig"), []rune("capt"), []rune("sgt"), []rune("cpl"), []rune("pvt"), []rune("pfc"), []rune("cmdr"),
	 []rune("adm"), []rune("lieut"), []rune("pte"),
	}
	titlesabv.Build()
	
	// Initate exceptions for titles
	titles.Key = [][]rune {
	 []rune("Sir"), []rune("Lord"), []rune("Baron"), []rune("Count"), []rune("Viscount"), []rune("Duke"), []rune("Marquess"), []rune("Earl"), []rune("Laird"), []rune("Master"), []rune("Bishop"), []rune("Father"), []rune("Sister"),
	 []rune("Pope"), []rune("Rabbi"), []rune("General"), []rune("Major"), []rune("Private"), []rune("Captain"), []rune("Sergent"), []rune("Commander"), []rune("Admiral"), []rune("Lieutenant"), []rune("Marquise"), []rune("Duca"),
	 []rune("Abbot"), []rune("Reverend"), []rune("Deacon"), []rune("Archbishop"), []rune("Cardinal"), []rune("Chancellor"), []rune("Chaplain"), []rune("Vicar"), []rune("Doctor"), []rune("Guru"), []rune("principe"), []rune("marchese"),
	 []rune("Prince"), []rune("King"), []rune("Queen"), []rune("Princess"), []rune("Emperor"), []rune("Caesar"), []rune("tsar"), []rune("Czar"), []rune("Csar"), []rune("Tzar"), []rune("Kaiser"), []rune("Sultan"), []rune("Conte"),
	 []rune("Dauphin"), []rune("Infante"), []rune("Margrave"), []rune("Marquis"), []rune("Freiherr"), []rune("Seigneur"), []rune("Nobile"), []rune("Baronet"), []rune("Dominus"), []rune("Vidame"), []rune("Vavasour"), []rune("Contessa"),
	 []rune("Vidame"), []rune("Kurfürst"), []rune("Prinz"), []rune("Viceroy"), []rune("Markgraf"), []rune("Graf"), []rune("Vizegraf"), []rune("Compte"), []rune("Comptesse"), []rune("Báró"), []rune("Barón"), []rune("Barone"),
	 []rune("Chevalier"), []rune("Ritter"), []rune("Cavaliere"), []rune("Nobiluomo"), []rune("Duque"), []rune("Príncipe"), []rune("Marquês"), []rune("Conde"), []rune("Visconde"), []rune("Barão"), []rune("Baronete"), []rune("Professor"),
	 []rune("Duchess"), []rune("Countess"), []rune("Baroness"), []rune("Dame"), []rune("Duc"), []rune("Viceroi"), []rune("Fürst"), []rune("Baronetto"), []rune("Principessa"), []rune("Visconte"), []rune("Princesse"), []rune("Roi"),
	 []rune("Reine"), []rune("Kaiserin"), []rune("König"), []rune("Königin"), []rune("Re"), []rune("Regina"), []rune("Rei"), []rune("Rainha"), []rune("Pape"), []rune("Papa"), []rune("Papst"), []rune("Monsieur"), []rune("Madame"),
	 []rune("Herr"), []rune("Père"), []rune("Padre"), []rune("Vater"), []rune("Saint"), []rune("Heilige"), []rune("San"), []rune("Arciduca"), []rune("Commodore"), []rune("Regent"), []rune("Lady"),
	}
	titles.Build()
	
	// Initate exceptions for mutli-part last names
	multilast.Key = [][]rune {
	 []rune("de"), []rune("da"), []rune("di"), []rune("von"), []rune("van"), []rune("le"), []rune("la"), []rune("du"), []rune("des"), []rune("del"), []rune("della"), []rune("der"),
	}
	multilast.Build()
	
	// Initate exceptions for honor
	honor.Key = [][]rune {
	 []rune("a.a"), []rune("a.a.s"), []rune("a.a.t"), []rune("a.o.t"), []rune("a.s"), []rune("b.a"), []rune("b.a.b.a"), []rune("b.a.com"), []rune("b.a.e"), []rune("b.a.ed"), []rune("b.arch"), []rune("b.a.s"), []rune("b.b.a"), 
	 []rune("b.b.e"), []rune("b.c.e"), []rune("b.che.e"), []rune("b.e.e"), []rune("b.f.a"), []rune("b.g.s"), []rune("b.i.arch"), []rune("b.in.dsn"), []rune("b.i.s"), []rune("b.i.s.e"), []rune("b.l.a"), []rune("b.m"), []rune("b.m.e"),
	 []rune("b.m.ed"), []rune("b.mtl.e"), []rune("b.p.f.e"), []rune("b.p.h.s"), []rune("b.s"), []rune("b.s.a.e"), []rune("b.s.b.a"), []rune("b.s.b.m.e"), []rune("b.s.c.b.a"), []rune("b.s.c.e"), []rune("b.s.che.e"), []rune("b.s.chem"),
	 []rune("b.s.c.s"), []rune("b.s.e"), []rune("b.s.ed"), []rune("b.s.e.e"), []rune("b.s.e.t"), []rune("b.s.geo"), []rune("b.s.h.e.s"), []rune("b.s.m.e"), []rune("b.s.met"), []rune("b.s.micr"), []rune("b.s.mt.e"), []rune("b.s.n"),
	 []rune("b.s.s.w"), []rune("b.s.w"), []rune("b.sw.e"), []rune("b.t.e"), []rune("b.t.m.t"), []rune("b.w.e"), []rune("ll.m"), []rune("m.a"), []rune("m.acc"), []rune("m.acct"), []rune("m.a.e"), []rune("m.a.ed"), []rune("m.agric"),
	 []rune("m.a.m"), []rune("m.aqua"), []rune("m.a.t"), []rune("m.b.a"), []rune("m.b.c"), []rune("m.c.d"), []rune("m.c.e"), []rune("m.che.e"), []rune("m.com"), []rune("m.comm.pl"), []rune("m.ed"), []rune("m.e.e"), []rune("m.eng"),
	 []rune("m.f.a"), []rune("m.his.st"), []rune("m.h.s"), []rune("m.i.d.c"), []rune("m.in.dsn"), []rune("m.i.s.e"), []rune("m.l.a"), []rune("m.l.arch"), []rune("m.l.i.s"), []rune("m.m"), []rune("m.m.e"), []rune("m.mtl.e"),
	 []rune("m.n.a"), []rune("m.n.r"), []rune("m.p.a"), []rune("m.p.h"), []rune("m.prob.s"), []rune("m.pr.s"), []rune("m.p.s"), []rune("m.p.t.c"), []rune("m.r.c"), []rune("m.r.e.d"), []rune("m.s"), []rune("m.s.a.e"),
	 []rune("m.s.b.m.e"), []rune("m.s.b.m.s"), []rune("m.s.c"), []rune("m.s.c.e"), []rune("m.s.che.e"), []rune("m.s.chem"), []rune("m.s.c.j"), []rune("m.s.c.s"), []rune("m.s.e"), []rune("m.s.ed"), []rune("m.s.e.e"),
	 []rune("m.s.e.s.m"), []rune("m.s.f.s"), []rune("m.s.h.a"), []rune("m.s.h.e.s"), []rune("m.s.h.i"), []rune("m.s.i.e"), []rune("m.s.i.l.a"), []rune("m.s.i.s"), []rune("m.s.j.p.s"), []rune("m.s.m.e"), []rune("m.s.met"),
	 []rune("m.s.m.sci"), []rune("m.s.mt.e"), []rune("m.s.n"), []rune("m.s.o.r"), []rune("m.s.o.t"), []rune("m.s.p.a.s"), []rune("m.s.p.h"), []rune("m.s.s.e"), []rune("m.s.w"), []rune("m.sw.e"), []rune("m.t.a"), []rune("m.tx"),
	 []rune("m.u.r.p"), []rune("ed.s"), []rune("au.d"), []rune("d.b.a"), []rune("d.m.a"), []rune("d.m.d"), []rune("d.n.p"), []rune("d.p.t"), []rune("dr.p.h"), []rune("d.sc"), []rune("d.v.m"), []rune("ed.d"), []rune("j.d"),
	 []rune("m.d"), []rune("o.d"), []rune("pharm.d"), []rune("ph.d"), []rune("e.g"), []rune("i.e"), []rune("lt.col"), []rune("d.d"),
	}
	honor.format = [][]rune {
	 []rune("A.A"), []rune("A.A.S"), []rune("A.A.T"), []rune("A.O.T"), []rune("A.S"), []rune("B.A"), []rune("B.A.B.A"), []rune("B.A.Com"), []rune("B.A.E"), []rune("B.A.Ed"), []rune("B.Arch"), []rune("B.A.S"), []rune("B.B.A"), 
	 []rune("B.B.E"), []rune("B.C.E"), []rune("B.Che.E"), []rune("B.E.E"), []rune("B.F.A"), []rune("B.G.S"), []rune("B.I.Arch"), []rune("B.In.Dsn"), []rune("B.I.S"), []rune("B.I.S.E"), []rune("B.L.A"), []rune("B.M"), []rune("B.M.E"),
	 []rune("B.M.Ed"), []rune("B.Mtl.E"), []rune("B.P.F.E"), []rune("B.P.H.S"), []rune("B.S"), []rune("B.S.A.E"), []rune("B.S.B.A"), []rune("B.S.B.M.E"), []rune("B.S.C.B.A"), []rune("B.S.C.E"), []rune("B.S.Che.E"), []rune("B.S.Chem"),
	 []rune("B.S.C.S"), []rune("B.S.E"), []rune("B.S.Ed"), []rune("B.S.E.E"), []rune("B.S.E.T"), []rune("B.S.Geo"), []rune("B.S.H.E.S"), []rune("B.S.M.E"), []rune("B.S.Met"), []rune("B.S.Micr"), []rune("B.S.Mt.E"), []rune("B.S.N"),
	 []rune("B.S.S.W"), []rune("B.S.W"), []rune("B.Sw.E"), []rune("B.T.E"), []rune("B.T.M.T"), []rune("B.W.E"), []rune("LL.M"), []rune("M.A"), []rune("M.Acc"), []rune("M.Acct"), []rune("M.A.E"), []rune("M.A.Ed"), []rune("M.Agric"),
	 []rune("M.A.M"), []rune("M.Aqua"), []rune("M.A.T"), []rune("M.B.A"), []rune("M.B.C"), []rune("M.C.D"), []rune("M.C.E"), []rune("M.Che.E"), []rune("M.Com"), []rune("M.Comm.Pl"), []rune("M.Ed"), []rune("M.E.E"), []rune("M.Eng"),
	 []rune("M.F.A"), []rune("M.His.St"), []rune("M.H.S"), []rune("M.I.D.C"), []rune("M.In.Dsn"), []rune("M.I.S.E"), []rune("M.L.A"), []rune("M.L.Arch"), []rune("M.L.I.S"), []rune("M.M"), []rune("M.M.E"), []rune("M.Mtl.E"),
	 []rune("M.N.A"), []rune("M.N.R"), []rune("M.P.A"), []rune("M.P.H"), []rune("M.Prob.S"), []rune("M.Pr.S"), []rune("M.P.S"), []rune("M.P.T.C"), []rune("M.R.C"), []rune("M.R.E.D"), []rune("M.S"), []rune("M.S.A.E"),
	 []rune("M.S.B.M.E"), []rune("M.S.B.M.S"), []rune("M.S.C"), []rune("M.S.C.E"), []rune("M.S.Che.E"), []rune("M.S.Chem"), []rune("M.S.C.J"), []rune("M.S.C.S"), []rune("M.S.E"), []rune("M.S.Ed"), []rune("M.S.E.E"),
	 []rune("M.S.E.S.M"), []rune("M.S.F.S"), []rune("M.S.H.A"), []rune("M.S.H.E.S"), []rune("M.S.H.I"), []rune("M.S.I.E"), []rune("M.S.I.L.A"), []rune("M.S.I.S"), []rune("M.S.J.P.S"), []rune("M.S.M.E"), []rune("M.S.Met"),
	 []rune("M.S.M.Sci"), []rune("M.S.Mt.E"), []rune("M.S.N"), []rune("M.S.O.R"), []rune("M.S.O.T"), []rune("M.S.P.A.S"), []rune("M.S.P.H"), []rune("M.S.S.E"), []rune("M.S.W"), []rune("M.Sw.E"), []rune("M.T.A"), []rune("M.Tx"),
	 []rune("M.U.R.P"), []rune("Ed.S"), []rune("Au.D"), []rune("D.B.A"), []rune("D.M.A"), []rune("D.M.D"), []rune("D.N.P"), []rune("D.P.T"), []rune("Dr.P.H"), []rune("D.Sc"), []rune("D.V.M"), []rune("Ed.D"), []rune("J.D"),
	 []rune("M.D"), []rune("O.D"), []rune("Pharm.D"), []rune("Ph.D"), []rune("E.g"), []rune("I.e"), []rune("Lt.Col"), []rune("D.D"),
	}
	temp := make([][]rune, len(honor.format))
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
	
	// Initiate exceptions for Portuguese small words
	portugueseSmall.Key = [][]rune {
	 []rune("à"), []rune("às"), []rune("ao"), []rune("da"), []rune("das"), []rune("de"), []rune("do"), []rune("e"), []rune("em"), []rune("na"), []rune("no"), []rune("o"), []rune("para"), []rune("pelo"), []rune("pelos"),
	 []rune("por"), []rune("se"), []rune("um"), []rune("uma"), []rune("pelas"), []rune("pela"),
	}
	portugueseSmall.Build()
	
	// Initiate exceptions for Spanish small words
	spanishSmall.Key = [][]rune {
	 []rune("a"), []rune("al"), []rune("de"), []rune("del"), []rune("e"), []rune("é"), []rune("el"), []rune("en"), []rune("la"), []rune("las"), []rune("los"), []rune("o"), []rune("ó"), []rune("para"), []rune("por"),
	 []rune("si"), []rune("un"), []rune("una"), []rune("y"),
	}
	spanishSmall.Build()
	
}

func equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range a {
		if c != b[i] {
			return false
		}
	}
	return true
}

// Structs
type wordStruct struct {
 content []rune
 isStart bool
 isEnd bool
 isHonor bool
 isTitle bool
 isRoman bool
 contraction uint8
 spaceAfter uint8 // 0=nothing, 1=space, 2=hypen, 3=slash, 4=end
 puncBefore []rune
 puncAfter []rune
}

type AuthorStruct struct {
 Last string
 First string
 Middle string
 Title string
 Suffix string
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
				noise = append(noise, i4 + 1)
			case 39, '’':
				if int(i4) < (i2 - i) - 2 {
					contraction = i4
				}
		}
		content[i4] = unicode.ToLower(w[i3])
		i4++
	}
	// Check if any noise occurred
	if len(noise) > 0 {
		if id, ok := honor.Find(content); ok { // if it's an honor then save the way it should be displayed
			isHonor = true
			content = honor.format[id]
		} else { // if it's not then split the word
			backup := r.runes
			saver := make([]rune, len(content))
			copy(saver, content)
			r.runes = puncBefore
			r.runes = append(r.runes, saver[0:noise[0]]...)
			r.len = len(r.runes)
			words = r.add(words, 1)
			r.runes = saver[noise[0]:]
			r.runes = append(r.runes, puncAfter...)
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
	words = append(words, wordStruct{content, false, isEnd, isHonor, false, false, contraction, spaceType, puncBefore, puncAfter})
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

// Removes 2 individual bytes from a slice of bytes
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

// Replaces an individual byte
func replaceRune(s []rune, from rune, to rune) {
	for i, v := range s {
		if v == from {
			s[i] = to
		}
	}
}

func English(str string) string {
	str, _ = format(str, Language_English, false)
	return str
}

func French(str string) string {
	str, _ = format(str, Language_French, false)
	return str
}

func German(str string) string {
	str, _ = format(str, Language_German, false)
	return str
}

func Italian(str string) string {
	str, _ = format(str, Language_Italian, false)
	return str
}

func Spanish(str string) string {
	str, _ = format(str, Language_Spanish, false)
	return str
}

func Portuguese(str string) string {
	str, _ = format(str, Language_Portuguese, false)
	return str
}

func Generic(str string) string {
	str, _ = format(str, Language_Generic, false)
	return str
}

func Author(str string, language uint8) (string, *AuthorStruct) {
	return format(str, language, true)
}

func format(str string, language uint8, formatAuthor bool) (string, *AuthorStruct) {

	if len(str) == 0 {
		return ``, nil
	}
	var small binsearch.Key_runes
	switch language {
		case Language_English: small = englishSmall
		case Language_French: small = frenchSmall
		case Language_German: small = germanSmall
		case Language_Italian: small = italianSmall
		case Language_Spanish: small = spanishSmall
		case Language_Portuguese: small = portugueseSmall
	}
	
	// Preprocessing
	str = html.UnescapeString(str)
	b := []byte(str)
	b = bytes.Replace(b, []byte("--"), []byte("—"), -1) // Correct hyphens to em dashes
	b = bytes.Replace(b, []byte("—"), []byte(" — "), -1) // Separate out em dashes
	b = bytes.Replace(b, []byte(" - "), []byte(" — "), -1) // Correct hyphens to em dashes
	b = bytes.Replace(b, []byte("[microform]"), []byte(""), -1)
	b = bytes.Trim(b, ` ;:.,`)
	if len(b) == 0 {
		return ``, nil
	}
	if b[0] == '(' {
		b = removeBytes(b, '(', ')')
	}
	if b[0] == '[' {
		b = removeBytes(b, '[', ']')
	}

	n := len(b)
	if n == 0 {
		return ``, nil
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
		return ``, nil
	}
	words[0].isStart = true
	for i=1; i<l; i++ {
		if words[i-1].isEnd {
			if words[i-1].puncAfter[0] != ',' {
				words[i].isStart = true
			}
		}
	}
	words[l-1].isEnd = true
	
	// On authors, delete the first word if it is "by" or "the"
	if formatAuthor {
		if equal(words[0].content, []rune("by")) || equal(words[0].content, []rune("the")) {
			words[0].content = make([]rune, 0)
		}
	}
	
	// Loop through all and apply rules
	var ws *wordStruct
	var content []rune
	var ok bool
	var ln int
	for i=0; i<l; i++ {
		ws = &words[i]
		content = ws.content
		ln = len(content)
		
		if ln == 0 {
			continue
		}
		
		if ws.isHonor {
			continue
		}
		
		// Uppercase roman numerals
		if isRoman(content) {
			ws.isRoman = true
			upperRune(content, -1) // -1 means uppercase all
			continue
		}
		
		// Titles
		if _, ok = titlesabv.Find(content); ok {
			upperRune(content, 0)
			ws.isTitle = true
			// Ensure title is followed by a period
			if len(ws.puncAfter) == 0 {
				ws.puncAfter = []rune(".")
			} else {
				ws.puncAfter[0] = '.'
			}
			continue
		}
		
		// Check for McStuff
		if ln > 3 {
			if content[0] == 'm' && content[1] == 'c' {
				upperRune(content, 0)
				upperRune(content, 2)
				replaceRune(ws.puncAfter, '.', ';')
				continue
			}
		}
		
		if language == Language_English {
			// Special for English: repair grammatical error on a -> an
			if ln == 1 {
				if content[0] == 'a' {
					if i < len(words) - 1 {
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
		} else {
			// Check for contractions
			if ws.contraction > 0 && language != Language_German {
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
			replaceRune(ws.puncAfter, '.', ';')
			continue
		}
		
		// Beginning and ending words need to be capitalized regardless of what they are
		if ws.isStart || ws.isEnd {
			upperRune(content, 0)
			if ln > 1 {
				replaceRune(ws.puncAfter, '.', ';')
			}
			continue
		}
		
		// Check for small words to keep lowercase, using binary search
		if _, ok = small.Find(content); ok {
			// Exception if it's 1 letter with following punctuation or the next word or previous word are also 1 letter
			if ln > 1 {
				replaceRune(ws.puncAfter, '.', ';')
				continue
			}
			if len(words[i-1].content) > 1 && len(words[i+1].content) > 1 {
				continue
			}
			upperRune(content, 0)
			continue
		}
		
		// Uppercase the first rune if none of the previous rules applied
		replaceRune(ws.puncAfter, '.', ';')
		upperRune(content, 0)
	}
	
	// Rebuild byte stream from words
	var buf bytes.Buffer
	for i=0; i<l; i++ {
		ws = &words[i]
		if len(ws.content) == 0 {
			continue
		}
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
	
	// If it's only a title (not author) then end here
	if !formatAuthor {
		return buf.String(), nil
	}
		
	author := new(AuthorStruct)
	
	// Find author's title
	var going bool
	for i=0; i<l; i++ {
		ws = &words[i]
		content = ws.content
		ln = len(content)
		if ln == 0 {
			continue
		}
		going = false
		
		if ws.isTitle {
			if len(author.Title) == 0 {
				author.Title = string(content) + `.`
			} else {
				switch words[i-1].spaceAfter {
					case 1: author.Title += ` ` + string(content) + `.`
					case 2: author.Title += `-` + string(content) + `.`
					case 3: author.Title += `/` + string(content) + `.`
				}
			}
			ws.content = content[0:0]
			going = true
		} else {
			if _, ok = titles.Find(content); ok {
				if len(author.Title) == 0 {
					author.Title = string(content)
				} else {
					switch words[i-1].spaceAfter {
						case 1: author.Title += ` ` + string(content)
						case 2: author.Title += `-` + string(content)
						case 3: author.Title += `/` + string(content)
					}
				}
				ws.content = content[0:0]
				going = true
			}
		}
		
		if !going {
			break
		}
	}
	
	// Find author's suffix & check for comma in puncAfter
	var comma int
	for i=0; i<l; i++ {
		ws = &words[i]
		content = ws.content
		ln = len(content)
		if ln == 0 {
			continue
		}
		
		if ws.isHonor {
			if len(author.Suffix) == 0 {
				author.Suffix = string(content) + `.`
			} else {
				if going {
					switch words[i-1].spaceAfter {
						case 1: author.Suffix += ` ` + string(content) + `.`
						case 2: author.Suffix += `-` + string(content) + `.`
						case 3: author.Suffix += `/` + string(content) + `.`
					}
				} else {
					author.Suffix += ` ` + string(content) + `.`
				}
			}
			ws.content = content[0:0]
			going = true
		} else {
			if comma == 0 && ws.spaceAfter == 1 {
				for _, r = range ws.puncAfter {
					if r == ',' {
						comma = i + 1
					}
				}
			}
			going = false
		}
	}
	
	// Rebuild what's left
	var first bytes.Buffer
	var middle bytes.Buffer
	var last bytes.Buffer
	
	// Get first and last name if there is a comma
	if comma > 0 {
		for i=0; i<comma; i++ {
			ws = &words[i]
			if len(ws.content) == 0 {
				continue
			}
			for _, r = range ws.puncBefore {
				last.WriteRune(r)
			}
			for _, r = range ws.content {
				last.WriteRune(r)
			}
			for _, r = range ws.puncAfter {
				if r != ',' {
					last.WriteRune(r)
				}
			}
			if i < comma - 1 {
				switch ws.spaceAfter {
					case 1: last.WriteByte(' ')
					case 2: last.WriteByte('-')
					case 3: last.WriteByte('/')
				}
			}
		}
		for ; i<l; i++ {
			ws = &words[i]
			if len(ws.content) == 0 {
				continue
			}
			for _, r = range ws.puncBefore {
				first.WriteRune(r)
			}
			for _, r = range ws.content {
				first.WriteRune(r)
			}
			for _, r = range ws.puncAfter {
				first.WriteRune(r)
			}
			switch ws.spaceAfter {
				case 1: i++; break
				case 2: first.WriteByte('-')
				case 3: first.WriteByte('/')
			}
		}
		for ; i<l; i++ {
			ws = &words[i]
			if len(ws.content) == 0 {
				continue
			}
			for _, r = range ws.puncBefore {
				middle.WriteRune(r)
			}
			for _, r = range ws.content {
				middle.WriteRune(r)
			}
			for _, r = range ws.puncAfter {
				middle.WriteRune(r)
			}
			switch ws.spaceAfter {
				case 1: middle.WriteByte(' ')
				case 2: middle.WriteByte('-')
				case 3: middle.WriteByte('/')
			}
		}
	} else {
		// Get first and last name if there is no comma
		going = true
		for i=l-1; i>=0; i-- {
			ws = &words[i]
			if len(ws.content) == 0 {
				continue
			}
			if going && ws.isRoman {
				i--
				continue
			}
			fmt.Println(i, string(ws.content))
			going = false
			if ws.spaceAfter > 1 {
				continue
			}
			if _, ok = multilast.Find(ws.content); ok {
				continue
			}
			if i < l - 1 {
				break
			}
		}
		lastpos := i + 1
		fmt.Println(`lastpos`, lastpos)
		for i=lastpos; i<l; i++ {
			ws = &words[i]
			if len(ws.content) == 0 {
				continue
			}
			for _, r = range ws.puncBefore {
				last.WriteRune(r)
			}
			for _, r = range ws.content {
				last.WriteRune(r)
			}
			for _, r = range ws.puncAfter {
				last.WriteRune(r)
			}
			switch ws.spaceAfter {
				case 1: last.WriteByte(' ')
				case 2: last.WriteByte('-')
				case 3: last.WriteByte('/')
			}
		}
		for i=0; i<lastpos; i++ {
			ws = &words[i]
			if len(ws.content) == 0 {
				continue
			}
			for _, r = range ws.puncBefore {
				first.WriteRune(r)
			}
			for _, r = range ws.content {
				first.WriteRune(r)
			}
			for _, r = range ws.puncAfter {
				first.WriteRune(r)
			}
			switch ws.spaceAfter {
				case 1: i++; break
				case 2: first.WriteByte('-')
				case 3: first.WriteByte('/')
			}
		}
		for ; i<lastpos; i++ {
			ws = &words[i]
			if len(ws.content) == 0 {
				continue
			}
			for _, r = range ws.puncBefore {
				middle.WriteRune(r)
			}
			for _, r = range ws.content {
				middle.WriteRune(r)
			}
			for _, r = range ws.puncAfter {
				middle.WriteRune(r)
			}
			switch ws.spaceAfter {
				case 1: middle.WriteByte(' ')
				case 2: middle.WriteByte('-')
				case 3: middle.WriteByte('/')
			}
		}
	}
	
	author.First = first.String()
	author.Middle = middle.String()
	author.Last = last.String()
	
	return buf.String(), author
}

