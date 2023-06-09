// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loong64

import (
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/obj"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/obj/loong64"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ir"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/objw"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/types"
)

func zerorange(pp *objw.Progs, p *obj.Prog, off, cnt int64, _ *uint32) *obj.Prog {
	if cnt == 0 {
		return p
	}
	if cnt < int64(4*types.PtrSize) {
		for i := int64(0); i < cnt; i += int64(types.PtrSize) {
			p = pp.Append(p, loong64.AMOVV, obj.TYPE_REG, loong64.REGZERO, 0, obj.TYPE_MEM, loong64.REGSP, 8+off+i)
		}
	} else if cnt <= int64(128*types.PtrSize) {
		p = pp.Append(p, loong64.AADDV, obj.TYPE_CONST, 0, 8+off-8, obj.TYPE_REG, loong64.REGRT1, 0)
		p.Reg = loong64.REGSP
		p = pp.Append(p, obj.ADUFFZERO, obj.TYPE_NONE, 0, 0, obj.TYPE_MEM, 0, 0)
		p.To.Name = obj.NAME_EXTERN
		p.To.Sym = ir.Syms.Duffzero
		p.To.Offset = 8 * (128 - cnt/int64(types.PtrSize))
	} else {
		//	ADDV	$(8+frame+lo-8), SP, r1
		//	ADDV	$cnt, r1, r2
		// loop:
		//	MOVV	R0, (Widthptr)r1
		//	ADDV	$Widthptr, r1
		//	BNE		r1, r2, loop
		p = pp.Append(p, loong64.AADDV, obj.TYPE_CONST, 0, 8+off-8, obj.TYPE_REG, loong64.REGRT1, 0)
		p.Reg = loong64.REGSP
		p = pp.Append(p, loong64.AADDV, obj.TYPE_CONST, 0, cnt, obj.TYPE_REG, loong64.REGRT2, 0)
		p.Reg = loong64.REGRT1
		p = pp.Append(p, loong64.AMOVV, obj.TYPE_REG, loong64.REGZERO, 0, obj.TYPE_MEM, loong64.REGRT1, int64(types.PtrSize))
		p1 := p
		p = pp.Append(p, loong64.AADDV, obj.TYPE_CONST, 0, int64(types.PtrSize), obj.TYPE_REG, loong64.REGRT1, 0)
		p = pp.Append(p, loong64.ABNE, obj.TYPE_REG, loong64.REGRT1, 0, obj.TYPE_BRANCH, 0, 0)
		p.Reg = loong64.REGRT2
		p.To.SetTarget(p1)
	}

	return p
}

func ginsnop(pp *objw.Progs) *obj.Prog {
	p := pp.Prog(loong64.ANOR)
	p.From.Type = obj.TYPE_REG
	p.From.Reg = loong64.REG_R0
	p.To.Type = obj.TYPE_REG
	p.To.Reg = loong64.REG_R0
	return p
}
