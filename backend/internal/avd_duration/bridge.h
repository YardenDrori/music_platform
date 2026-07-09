#ifndef AVD_BRIDGE_H
#define AVD_BRIDGE_H

#include <stdint.h>

// typedef int (*avdread)(void *, uint8_t *, int);
// typedef int64_t (*avdseek)(void *, int64_t, int);

int avd_averr_eof();
int avd_averr_eio();
int avd_averr_bug();
int avd_avseek_file_size();
int avd_averr_enosys();

double avdGetFileRuntime(uint64_t fileId);

#endif
