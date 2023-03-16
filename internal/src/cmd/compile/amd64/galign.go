// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package amd64

import (
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/obj/x86"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ssagen"
)

var leaptr = x86.ALEAQ

func Init(arch *ssagen.ArchInfo) {
	arch.LinkArch = &x86.Linkamd64
	arch.REGSP = x86.REGSP
	arch.MAXWIDTH = 1 << 50

	arch.ZeroRange = zerorange
	arch.Ginsnop = ginsnop

	arch.SSAMarkMoves = ssaMarkMoves
	arch.SSAGenValue = ssaGenValue
	arch.SSAGenBlock = ssaGenBlock
	arch.LoadRegResult = loadRegResult
	arch.SpillArgReg = spillArgReg
}
