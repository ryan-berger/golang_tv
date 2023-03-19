#include "llvm-c/Core.h"
#include "llvm-c/DebugInfo.h"

#ifdef __cplusplus
extern "C" {
#endif

void validate(LLVMModuleRef m, LLVMValueRef src, LLVMValueRef tgt);

#ifdef __cplusplus
}
#endif