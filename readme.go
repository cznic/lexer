// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: 
//		jnml, labs.nic.cz
//		Miek Gieben, SIDN, miek@miek.nl

/*

Syntax supported by ParseRE (ATM a very basic subset of RE2, docs bellow are a mod of: http://code.google.com/p/re2/wiki/Syntax, original docs license unclear)

Single characters:
	.            any character, excluding newline
	[xyz]        character class
	[^xyz]       negated character class
	\p{Greek}    Unicode character class
	\P{Greek}    negated Unicode character class

Composites:
	xy           x followed by y
	x|y          x or y

Repetitions:
	x*           zero or more x
	x+           one or more x
	x?           zero or one x

Grouping:
	(re)         group

Empty strings:
	^            at beginning of text or line
	$            at end of text or line
	\A           at beginning of text
	\z           at end of text

Escape sequences:
	\a           bell (≡ \007)
	\b           backspace (≡ \010) 
	\f           form feed (≡ \014)
	\n           newline (≡ \012)
	\r           carriage return (≡ \015)
	\t           horizontal tab (≡ \011)
	\v           vertical tab character (≡ \013)
	\M           M is one of metachars \.+*?()|[]^$
	\xhh         rune \u00hh, h is a hex digit

Character class elements:
	x            single Unicode character
	A-Z          Unicode character range (inclusive)

Unicode character class names--general category:
	Cc           control
	Cf           format
	Co           private use
	Cs           surrogate
	letter       Lu, Ll, Lt, Lm, or Lo
	Ll           lowercase letter
	Lm           modifier letter
	Lo           other letter
	Lt           titlecase letter
	Lu           uppercase letter
	Mc           spacing mark
	Me           enclosing mark
	Mn           non-spacing mark
	Nd           decimal number
	Nl           letter number
	No           other number
	Pc           connector punctuation
	Pd           dash punctuation
	Pe           close punctuation
	Pf           final punctuation
	Pi           initial punctuation
	Po           other punctuation
	Ps           open punctuation
	Sc           currency symbol
	Sk           modifier symbol
	Sm           math symbol
	So           other symbol
	Zl           line separator
	Zp           paragraph separator
	Zs           space separator

Unicode character class names--scripts:
	Arabic                 Arabic
	Armenian               Armenian
	Avestan                Avestan
	Balinese               Balinese
	Bamum                  Bamum
	Bengali                Bengali
	Bopomofo               Bopomofo
	Braille                Braille
	Buginese               Buginese
	Buhid                  Buhid
	Canadian_Aboriginal    Canadian Aboriginal
	Carian                 Carian
	Common                 Common
	Coptic                 Coptic
	Cuneiform              Cuneiform
	Cypriot                Cypriot
	Cyrillic               Cyrillic
	Deseret                Deseret
	Devanagari             Devanagari
	Egyptian_Hieroglyphs   Egyptian Hieroglyphs
	Ethiopic               Ethiopic
	Georgian               Georgian
	Glagolitic             Glagolitic
	Gothic                 Gothic
	Greek                  Greek
	Gujarati               Gujarati
	Gurmukhi               Gurmukhi
	Hangul                 Hangul
	Han                    Han
	Hanunoo                Hanunoo
	Hebrew                 Hebrew
	Hiragana               Hiragana
	Cham                   Cham
	Cherokee               Cherokee
	Imperial_Aramaic       Imperial Aramaic
	Inherited              Inherited
	Inscriptional_Pahlavi  Inscriptional Pahlavi
	Inscriptional_Parthian Inscriptional Parthian
	Javanese               Javanese
	Kaithi                 Kaithi
	Kannada                Kannada
	Katakana               Katakana
	Kayah_Li               Kayah Li
	Kharoshthi             Kharoshthi
	Khmer                  Khmer
	Lao                    Lao
	Latin                  Latin
	Lepcha                 Lepcha
	Limbu                  Limbu
	Linear_B               Linear B
	Lisu                   Lisu
	Lycian                 Lycian
	Lydian                 Lydian
	Malayalam              Malayalam
	Meetei_Mayek           Meetei Mayek
	Mongolian              Mongolian
	Myanmar                Myanmar
	New_Tai_Lue            New Tai Lue
	Nko                    Nko
	Ogham                  Ogham
	Old_Italic             Old Italic
	Old_Persian            Old Persian
	Old_South_Arabian      Old South Arabian
	Old_Turkic             Old Turkic
	Ol_Chiki               Ol Chiki
	Oriya                  Oriya
	Osmanya                Osmanya
	Phags_Pa               Phags Pa
	Phoenician             Phoenician
	Rejang                 Rejang
	Runic                  Runic
	Samaritan              Samaritan
	Saurashtra             Saurashtra
	Shavian                Shavian
	Sinhala                Sinhala
	Sundanese              Sundanese
	Syloti_Nagri           Syloti Nagri
	Syriac                 Syriac
	Tagalog                Tagalog
	Tagbanwa               Tagbanwa
	Tai_Le                 Tai Le
	Tai_Tham               Tai Tham
	Tai_Viet               Tai Viet
	Tamil                  Tamil
	Telugu                 Telugu
	Thaana                 Thaana
	Thai                   Thai
	Tibetan                Tibetan
	Tifinagh               Tifinagh
	Ugaritic               Ugaritic
	Vai                    Vai
	Yi                     Yi


*/
package readme
