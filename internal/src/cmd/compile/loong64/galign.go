// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loong64

import (
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/obj/loong64"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ssa"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ssagen"
)

func Init(arch *ssagen.ArchInfo) {
	arch.LinkArch = &loong64.Linkloong64
	arch.REGSP = loong64.REGSP
	arch.MAXWIDTH = 1 << 50
	arch.ZeroRange = zerorange
	arch.Ginsnop = ginsnop

	arch.SSAMarkMoves = func(s *ssagen.State, b *ssa.Block) {}
	arch.SSAGenValue = ssaGenValue
	arch.SSAGenBlock = ssaGenBlock
}
