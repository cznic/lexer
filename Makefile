# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=github.com/cznic/lexer

GOFILES=\
	is.go\
	lexer.go\
	nfa.go\
	ranges.go\
	regex.go\
	scanner.go\
	source.go\
	vm.go\

include $(GOROOT)/src/Make.pkg
