// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86

import (
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/obj/x86"
	"fmt"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/base"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ssagen"
	"github.com/ryan-berger/golang_tv/internal/src/internal/buildcfg"
	"os"
)

func Init(arch *ssagen.ArchInfo) {
	arch.LinkArch = &x86.Link386
	arch.REGSP = x86.REGSP
	arch.SSAGenValue = ssaGenValue
	arch.SSAGenBlock = ssaGenBlock
	arch.MAXWIDTH = (1 << 32) - 1
	switch v := buildcfg.GO386; v {
	case "sse2":
	case "softfloat":
		arch.SoftFloat = true
	case "387":
		fmt.Fprintf(os.Stderr, "unsupported setting GO386=387. Consider using GO386=softfloat instead.\n")
		base.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "unsupported setting GO386=%s\n", v)
		base.Exit(1)

	}

	arch.ZeroRange = zerorange
	arch.Ginsnop = ginsnop

	arch.SSAMarkMoves = ssaMarkMoves
}
