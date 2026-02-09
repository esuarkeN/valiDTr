
#line 1 "cgo-builtin-prolog"
#include <stddef.h>

/* Define intgo when compiling with GCC.  */
typedef ptrdiff_t intgo;

#define GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; intgo n; } _GoString_;
typedef struct { char *p; intgo n; intgo c; } _GoBytes_;
_GoString_ GoString(char *p);
_GoString_ GoStringN(char *p, int l);
_GoBytes_ GoBytes(void *p, int n);
char *CString(_GoString_);
void *CBytes(_GoBytes_);
void *_CMalloc(size_t);

__attribute__ ((unused))
static size_t _GoStringLen(_GoString_ s) { return (size_t)s.n; }

__attribute__ ((unused))
static const char *_GoStringPtr(_GoString_ s) { return s.p; }

#line 8 "C:\\Users\\nekkr\\go\\pkg\\mod\\github.com\\mattn\\go-sqlite3@v1.14.33\\sqlite3_context.go"


#ifndef USE_LIBSQLITE3
#include "sqlite3-binding.h"
#else
#include <sqlite3.h>
#endif
#include <stdlib.h>
// These wrappers are necessary because SQLITE_TRANSIENT
// is a pointer constant, and cgo doesn't translate them correctly.

static inline void my_result_text(sqlite3_context *ctx, char *p, int np) {
	sqlite3_result_text(ctx, p, np, SQLITE_TRANSIENT);
}

static inline void my_result_blob(sqlite3_context *ctx, void *p, int np) {
	sqlite3_result_blob(ctx, p, np, SQLITE_TRANSIENT);
}

#line 1 "cgo-generated-wrapper"


#line 1 "cgo-gcc-prolog"
/*
  If x and y are not equal, the type will be invalid
  (have a negative array count) and an inscrutable error will come
  out of the compiler and hopefully mention "name".
*/
#define __cgo_compile_assert_eq(x, y, name) typedef char name[(x-y)*(x-y)*-2UL+1UL];

/* Check at compile time that the sizes we use match our expectations. */
#define __cgo_size_assert(t, n) __cgo_compile_assert_eq(sizeof(t), (size_t)n, _cgo_sizeof_##t##_is_not_##n)

__cgo_size_assert(char, 1)
__cgo_size_assert(short, 2)
__cgo_size_assert(int, 4)
typedef long long __cgo_long_long;
__cgo_size_assert(__cgo_long_long, 8)
__cgo_size_assert(float, 4)
__cgo_size_assert(double, 8)

extern char* _cgo_topofstack(void);

/*
  We use packed structs, but they are always aligned.
  The pragmas and address-of-packed-member are only recognized as warning
  groups in clang 4.0+, so ignore unknown pragmas first.
*/
#pragma GCC diagnostic ignored "-Wunknown-pragmas"
#pragma GCC diagnostic ignored "-Wpragmas"
#pragma GCC diagnostic ignored "-Waddress-of-packed-member"
#pragma GCC diagnostic ignored "-Wunknown-warning-option"
#pragma GCC diagnostic ignored "-Wunaligned-access"

#include <errno.h>
#include <string.h>


#define CGO_NO_SANITIZE_THREAD
#define _cgo_tsan_acquire()
#define _cgo_tsan_release()


#define _cgo_msan_write(addr, sz)

CGO_NO_SANITIZE_THREAD
void
_cgo_0470baca2faa_Cfunc_my_result_blob(void *v)
{
	struct {
		sqlite3_context* p0;
		void* p1;
		int p2;
		char __pad20[4];
	} __attribute__((__packed__, __gcc_struct__)) *_cgo_a = v;
	_cgo_tsan_acquire();
	my_result_blob(_cgo_a->p0, _cgo_a->p1, _cgo_a->p2);
	_cgo_tsan_release();
}

CGO_NO_SANITIZE_THREAD
void
_cgo_0470baca2faa_Cfunc_my_result_text(void *v)
{
	struct {
		sqlite3_context* p0;
		char* p1;
		int p2;
		char __pad20[4];
	} __attribute__((__packed__, __gcc_struct__)) *_cgo_a = v;
	_cgo_tsan_acquire();
	my_result_text(_cgo_a->p0, _cgo_a->p1, _cgo_a->p2);
	_cgo_tsan_release();
}

CGO_NO_SANITIZE_THREAD
void
_cgo_0470baca2faa_Cfunc_sqlite3_result_error_toobig(void *v)
{
	struct {
		sqlite3_context* p0;
	} __attribute__((__packed__, __gcc_struct__)) *_cgo_a = v;
	_cgo_tsan_acquire();
	sqlite3_result_error_toobig(_cgo_a->p0);
	_cgo_tsan_release();
}

CGO_NO_SANITIZE_THREAD
void
_cgo_0470baca2faa_Cfunc_sqlite3_result_int(void *v)
{
	struct {
		sqlite3_context* p0;
		int p1;
		char __pad12[4];
	} __attribute__((__packed__, __gcc_struct__)) *_cgo_a = v;
	_cgo_tsan_acquire();
	sqlite3_result_int(_cgo_a->p0, _cgo_a->p1);
	_cgo_tsan_release();
}

CGO_NO_SANITIZE_THREAD
void
_cgo_0470baca2faa_Cfunc_sqlite3_result_zeroblob(void *v)
{
	struct {
		sqlite3_context* p0;
		int p1;
		char __pad12[4];
	} __attribute__((__packed__, __gcc_struct__)) *_cgo_a = v;
	_cgo_tsan_acquire();
	sqlite3_result_zeroblob(_cgo_a->p0, _cgo_a->p1);
	_cgo_tsan_release();
}

