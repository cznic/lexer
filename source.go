// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: 
//		jnml, labs.nic.cz
//		Miek Gieben, SIDN, miek@miek.nl


package lexer

import (
	"go/token"
	"io"
	"os"
)

// EOFReader implements a RuneReader allways returning 0 (EOF) 
type EOFReader int

func (r EOFReader) ReadRune() (rune int, size int, err os.Error) {
	return 0, 0, os.EOF
}

type source struct {
	reader   io.RuneReader
	position token.Position
}

func newSource(fname string, r io.RuneReader) (s source) {
	s.reader = r
	s.position.Filename = fname
	s.position.Line = 1
	s.position.Column = 1
	return
}

// Source provides a stack of rune streams with position information.
type Source struct {
	stack []source
	tos   source
}

// NewSource returns a new Source from a RuneReader having fname.
// The RuneReader can be nil. Then an EOFReader is supplied and
// the real RuneReader(s) can be Included anytime afterwards.
func NewSource(fname string, r io.RuneReader) *Source {
	s := &Source{}
	if r == nil {
		r = EOFReader(0)
	}
	s.tos = newSource(fname, r)
	return s
}

// Include includes a RuneReader having fname. Recursive including is not checked.
func (s *Source) Include(fname string, r io.RuneReader) {
	s.stack = append(s.stack, s.tos)
	s.tos = newSource(fname, r)
}

// Position return the position of the next Read.
func (s *Source) Position() token.Position {
	return s.tos.position
}

// Read returns the next Source ScannerRune.
func (s *Source) Read() (r ScannerRune) {
	for {
		r.Position = s.Position()
		r.Rune, r.Size, r.Err = s.tos.reader.ReadRune()
		if r.Err == nil || r.Err != os.EOF {
			p := &s.tos.position
			p.Offset += r.Size
			if r.Rune != '\n' {
				p.Column++
			} else {
				p.Line++
				p.Column = 1
			}
			return
		}

		// err == os.EOF, try parent source
		if sp := len(s.stack) - 1; sp >= 0 {
			s.tos = s.stack[sp]
			s.stack = s.stack[:sp]
		} else {
			r.Rune, r.Size = 0, 0
			return
		}
	}
	panic("unreachable")
}

// ScannerRune is a struct holding info about a rune and it's origin
type ScannerRune struct {
	Position token.Position // Starting position of Rune
	Rune     int            // Rune value
	Size     int            // Rune size
	Err      os.Error       // os.EOF or nil. Any other value invalidates all other fields of a ScannerRune.
}

// ScannerSource is a Source with one ScannerRune look behind and an on demand one ScannerRune lookahead.
type ScannerSource struct {
	source  *Source
	prev    ScannerRune
	current ScannerRune
	next    ScannerRune
	runes   []int
}

// Accept checks if rune matches Current. If true then does Move.
func (s *ScannerSource) Accept(rune int) bool {
	if rune == s.Current() {
		s.Move()
		return true
	}

	return false
}

// NewScannerSource returns a new ScannerSource from a RuneReader having fname.
// The RuneReader can be nil. Then an EOFReader is supplied and
// the real RuneReader(s) can be Included anytime afterwards.
func NewScannerSource(fname string, r io.RuneReader) *ScannerSource {
	s := &ScannerSource{}
	s.source = NewSource(fname, r)
	s.Move()
	return s
}

// Collect returns all runes seen by the ScannerSource since last Collect or CollectString.
// Either Collect or CollectString can be called but only one of them as both clears the collector.
func (s *ScannerSource) Collect() (runes []int) {
	runes, s.runes = s.runes, nil
	return
}

// CollectString returns all runes seen by the ScannerSource since last CollectString or Collect as a string.
// Either Collect or CollectString can be called but only one of them as both clears the collector.
func (s *ScannerSource) CollectString() string {
	return string(s.Collect())
}

// CurrentRune returns the current ScannerSource rune. At EOF it's zero.
func (s *ScannerSource) Current() int {
	return s.current.Rune
}

// Current returns the current ScannerSource ScannerRune.
func (s *ScannerSource) CurrentRune() ScannerRune {
	return s.current
}

// Include includes a RuneReader having fname. Recursive including is not checked.
// Include discards the one rune lookahead data if there are any.
// Lookahead data exists iff Next() has been called and Move() has not yet been called afterwards.
func (s *ScannerSource) Include(fname string, r io.RuneReader) {
	s.invalidateNext()
	s.runes = nil
	s.source.Include(fname, r)
	s.Move()
}

func (s *ScannerSource) invalidateNext() {
	s.next.Position.Line = 0
}

func (s *ScannerSource) lookahead() {
	if !s.next.Position.IsValid() {
		s.read(&s.next)
	}
}

// Move moves ScannerSource one rune ahead.
func (s *ScannerSource) Move() {
	if rune := s.Current(); rune != 0 { // collect
		s.runes = append(s.runes, rune)
	}
	s.prev = s.current
	if s.next.Position.IsValid() {
		s.current = s.next
		s.invalidateNext()
		return
	} else {
		s.read(&s.current)
	}
}

// Next returns ScannerSource next (lookahead) ScannerRune. It's Rune is zero if next is EOF.
func (s *ScannerSource) NextRune() ScannerRune {
	s.lookahead()
	return s.next
}

// NextRune returns ScannerSource next (lookahead) rune. It is zero if next is EOF
func (s *ScannerSource) Next() int {
	s.lookahead()
	return s.next.Rune
}

// Position returns the current ScannerSource position, i.e. after a Move() it returns the position after CurrentRune.
func (s *ScannerSource) Position() token.Position {
	return s.source.Position()
}

// Prev returns then previous (look behind) ScanerRune. Before first Move() it's Rune is zero and Position.IsValid == false
func (s *ScannerSource) PrevRune() ScannerRune {
	return s.prev
}

// PrevRune returns the previous (look behind) ScannerRune rune. Before first Move() it's zero.
func (s *ScannerSource) Prev() int {
	return s.prev.Rune
}

func (s *ScannerSource) read(dest *ScannerRune) {
	*dest = s.source.Read()
}
