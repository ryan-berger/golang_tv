// Code generated by stringer -i cpu.go -o anames.go -p riscv; DO NOT EDIT.

package riscv

import "github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/obj"

var Anames = []string{
	obj.A_ARCHSPECIFIC: "ADDI",
	"SLTI",
	"SLTIU",
	"ANDI",
	"ORI",
	"XORI",
	"SLLI",
	"SRLI",
	"SRAI",
	"LUI",
	"AUIPC",
	"ADD",
	"SLT",
	"SLTU",
	"AND",
	"OR",
	"XOR",
	"SLL",
	"SRL",
	"SUB",
	"SRA",
	"SLLIRV32",
	"SRLIRV32",
	"SRAIRV32",
	"JAL",
	"JALR",
	"BEQ",
	"BNE",
	"BLT",
	"BLTU",
	"BGE",
	"BGEU",
	"LW",
	"LWU",
	"LH",
	"LHU",
	"LB",
	"LBU",
	"SW",
	"SH",
	"SB",
	"FENCE",
	"FENCEI",
	"FENCETSO",
	"ADDIW",
	"SLLIW",
	"SRLIW",
	"SRAIW",
	"ADDW",
	"SLLW",
	"SRLW",
	"SUBW",
	"SRAW",
	"LD",
	"SD",
	"MUL",
	"MULH",
	"MULHU",
	"MULHSU",
	"MULW",
	"DIV",
	"DIVU",
	"REM",
	"REMU",
	"DIVW",
	"DIVUW",
	"REMW",
	"REMUW",
	"LRD",
	"SCD",
	"LRW",
	"SCW",
	"AMOSWAPD",
	"AMOADDD",
	"AMOANDD",
	"AMOORD",
	"AMOXORD",
	"AMOMAXD",
	"AMOMAXUD",
	"AMOMIND",
	"AMOMINUD",
	"AMOSWAPW",
	"AMOADDW",
	"AMOANDW",
	"AMOORW",
	"AMOXORW",
	"AMOMAXW",
	"AMOMAXUW",
	"AMOMINW",
	"AMOMINUW",
	"RDCYCLE",
	"RDCYCLEH",
	"RDTIME",
	"RDTIMEH",
	"RDINSTRET",
	"RDINSTRETH",
	"FRCSR",
	"FSCSR",
	"FRRM",
	"FSRM",
	"FRFLAGS",
	"FSFLAGS",
	"FSRMI",
	"FSFLAGSI",
	"FLW",
	"FSW",
	"FADDS",
	"FSUBS",
	"FMULS",
	"FDIVS",
	"FMINS",
	"FMAXS",
	"FSQRTS",
	"FMADDS",
	"FMSUBS",
	"FNMADDS",
	"FNMSUBS",
	"FCVTWS",
	"FCVTLS",
	"FCVTSW",
	"FCVTSL",
	"FCVTWUS",
	"FCVTLUS",
	"FCVTSWU",
	"FCVTSLU",
	"FSGNJS",
	"FSGNJNS",
	"FSGNJXS",
	"FMVXS",
	"FMVSX",
	"FMVXW",
	"FMVWX",
	"FEQS",
	"FLTS",
	"FLES",
	"FCLASSS",
	"FLD",
	"FSD",
	"FADDD",
	"FSUBD",
	"FMULD",
	"FDIVD",
	"FMIND",
	"FMAXD",
	"FSQRTD",
	"FMADDD",
	"FMSUBD",
	"FNMADDD",
	"FNMSUBD",
	"FCVTWD",
	"FCVTLD",
	"FCVTDW",
	"FCVTDL",
	"FCVTWUD",
	"FCVTLUD",
	"FCVTDWU",
	"FCVTDLU",
	"FCVTSD",
	"FCVTDS",
	"FSGNJD",
	"FSGNJND",
	"FSGNJXD",
	"FMVXD",
	"FMVDX",
	"FEQD",
	"FLTD",
	"FLED",
	"FCLASSD",
	"FLQ",
	"FSQ",
	"FADDQ",
	"FSUBQ",
	"FMULQ",
	"FDIVQ",
	"FMINQ",
	"FMAXQ",
	"FSQRTQ",
	"FMADDQ",
	"FMSUBQ",
	"FNMADDQ",
	"FNMSUBQ",
	"FCVTWQ",
	"FCVTLQ",
	"FCVTSQ",
	"FCVTDQ",
	"FCVTQW",
	"FCVTQL",
	"FCVTQS",
	"FCVTQD",
	"FCVTWUQ",
	"FCVTLUQ",
	"FCVTQWU",
	"FCVTQLU",
	"FSGNJQ",
	"FSGNJNQ",
	"FSGNJXQ",
	"FMVXQ",
	"FMVQX",
	"FEQQ",
	"FLEQ",
	"FLTQ",
	"FCLASSQ",
	"CSRRW",
	"CSRRS",
	"CSRRC",
	"CSRRWI",
	"CSRRSI",
	"CSRRCI",
	"ECALL",
	"SCALL",
	"EBREAK",
	"SBREAK",
	"MRET",
	"SRET",
	"URET",
	"DRET",
	"WFI",
	"SFENCEVMA",
	"HFENCEGVMA",
	"HFENCEVVMA",
	"WORD",
	"BEQZ",
	"BGEZ",
	"BGT",
	"BGTU",
	"BGTZ",
	"BLE",
	"BLEU",
	"BLEZ",
	"BLTZ",
	"BNEZ",
	"FABSD",
	"FABSS",
	"FNEGD",
	"FNEGS",
	"FNED",
	"FNES",
	"MOV",
	"MOVB",
	"MOVBU",
	"MOVF",
	"MOVD",
	"MOVH",
	"MOVHU",
	"MOVW",
	"MOVWU",
	"NEG",
	"NEGW",
	"NOT",
	"SEQZ",
	"SNEZ",
	"LAST",
}
