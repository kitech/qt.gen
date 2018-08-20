#ifndef _EXTRA_EXPORT_H_
#define _EXTRA_EXPORT_H_

// clang reports 'extern "C" ...' as an unexposed decl, which we definitely
// need to recurse into.
// example:
// extern "C"
//  const char *qVersion(void) noexcept;

// some not exported useful functions in any headers
bool qRegisterResourceData(int, const unsigned char *, const unsigned char *, const unsigned char *);
bool qUnregisterResourceData(int, const unsigned char *, const unsigned char *, const unsigned char *);

#endif

