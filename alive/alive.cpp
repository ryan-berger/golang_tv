#include "alive.h"

#include <iostream>

#include "llvm/Analysis/TargetLibraryInfo.h"
#include "llvm/IR/Function.h"
#include "llvm/IR/Module.h"

#include "cache/cache.h"
#include "llvm_util/compare.h"
#include "llvm_util/llvm2alive.h"
#include "llvm_util/llvm_optimizer.h"
#include "llvm_util/utils.h"
#include "smt/smt.h"
#include "tools/transform.h"
#include "util/version.h"
#include "util/config.h"

void validate(LLVMModuleRef m, LLVMValueRef src, LLVMValueRef tgt) {
    auto srcFn = llvm::unwrap<llvm::Function>(src);
    auto tgtFn = llvm::unwrap<llvm::Function>(tgt);

    auto mod = llvm::unwrap(m);

    auto &DL = mod->getDataLayout();
    auto targetTriple = llvm::Triple(mod->getTargetTriple());
    auto TLI = llvm::TargetLibraryInfoWrapperPass(targetTriple);


    // go doesn't have undef/poison??
    util::config::disable_undef_input = true;
    util::config::disable_poison_input = true;

    llvm_util::initializer llvm_util_init(std::cout, DL);
    smt::smt_initializer smt_init;
    llvm_util::Verifier verifier(TLI, smt_init, std::cout);

    verifier.compareFunctions(*srcFn, *tgtFn);
}
