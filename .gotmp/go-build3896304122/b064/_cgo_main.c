#include <stddef.h>
int main() { return 0; }
void crosscall2(void(*fn)(void*) __attribute__((unused)), void *a __attribute__((unused)), int c __attribute__((unused)), size_t ctxt __attribute__((unused))) { }
size_t _cgo_wait_runtime_init_done(void) { return 0; }
void _cgo_release_context(size_t ctxt __attribute__((unused))) { }
char* _cgo_topofstack(void) { return (char*)0; }
void _cgo_allocate(void *a __attribute__((unused)), int c __attribute__((unused))) { }
void _cgo_panic(void *a __attribute__((unused)), int c __attribute__((unused))) { }
void _cgo_reginit(void) { }
#line 1 "cgo-generated-wrappers"
extern void authorizerTrampoline();
extern void callbackTrampoline();
extern void commitHookTrampoline();
extern void compareTrampoline();
extern void doneTrampoline();
extern void rollbackHookTrampoline();
extern void stepTrampoline();
extern void updateHookTrampoline();
void _cgoexp_0470baca2faa_callbackTrampoline(void* p){}
void _cgoexp_0470baca2faa_stepTrampoline(void* p){}
void _cgoexp_0470baca2faa_doneTrampoline(void* p){}
void _cgoexp_0470baca2faa_compareTrampoline(void* p){}
void _cgoexp_0470baca2faa_commitHookTrampoline(void* p){}
void _cgoexp_0470baca2faa_rollbackHookTrampoline(void* p){}
void _cgoexp_0470baca2faa_updateHookTrampoline(void* p){}
void _cgoexp_0470baca2faa_authorizerTrampoline(void* p){}
void _cgoexp_0470baca2faa_preUpdateHookTrampoline(void* p){}
