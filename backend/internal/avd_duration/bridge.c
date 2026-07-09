#include "bridge.h"
#include "_cgo_export.h"
#include <asm-generic/errno-base.h>
#include <asm-generic/errno.h>
#include <libavformat/avio.h>
#include <libavutil/error.h>
#include <libavutil/mem.h>
#include <stdint.h>


int avd_averr_eof() { return AVERROR_EOF; }
int avd_averr_bug() { return AVERROR_BUG; }
int avd_avseek_file_size() { return AVSEEK_SIZE; }
int avd_averr_eio() { return AVERROR(EIO); }
int avd_averr_enosys() { return AVERROR(ENOSYS); }
