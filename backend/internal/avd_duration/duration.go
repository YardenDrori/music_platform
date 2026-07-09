// Package avdduration provides audio file duration extraction via a CGO bridge to FFmpeg's libavformat.
package avdduration

/*
#cgo pkg-config: libavformat libavcodec libavutil
#include "bridge.h"
*/
import "C"

//go:generate go tool cgo -exportheader=cgo_export.h duration.go
//go:generate mv _obj/_cgo_export.h ./
//go:generate rm -rf _obj _cgo_2.o _cgo_4.o

//to use run "go generate ./..."

import (
	"errors"
	"io"
	"log/slog"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var fileRegistry = sync.Map{}

// if this actually wraps and the old ids havent been freed then fuck me i guess
var nextID uint64 = 0

//export avdread
func avdread(fileIDPtr unsafe.Pointer, bufferPointer *C.uint8_t, bufferSize C.int) C.int {
	//read the value from fileIdPointer and convert it to a go uint64
	fileID := uint64(*(*C.uint64_t)(fileIDPtr))

	size := int(bufferSize)

	//have a slice pointing over the same memory C allocated for the buffer
	bufPtr := (*byte)(unsafe.Pointer(bufferPointer))
	buffer := unsafe.Slice(bufPtr, size)

	rawStream, ok := fileRegistry.Load(fileID)
	if !ok {
		slog.Error(
			"reading stream: file id for requested stream not found in registry",
			"fileID",
			fileID,
		)
		return C.int(C.avd_averr_bug())
	}
	stream, ok := rawStream.(io.ReadSeekCloser)
	if !ok {
		slog.Error(
			"reading stream: item in registry isnt an io.ReadSeekCloser... somehow, bro how even?! anyway it was a ",
			"type",
			reflect.TypeOf(rawStream).String(),
			"fileID",
			fileID,
		)
		return C.int(C.avd_averr_bug())
	}

	bytesRead, err := stream.Read(buffer)
	if err != nil {
		//Read can return bytes read and an error in the same call even if some bytes were still read
		if bytesRead > 0 {
			return C.int(bytesRead)
		}

		if errors.Is(err, io.EOF) {
			return C.int(C.avd_averr_eof()) //'E' 'O' 'F' ' ' for ffmpeg
		}

		slog.Error("reading stream", "error", err)
		return C.int(C.avd_averr_eio()) //posix I/O error code (5)
	}

	return C.int(bytesRead)
}

//export avdseek
func avdseek(fileIDPtr unsafe.Pointer, offset C.int64_t, whence C.int) C.int64_t {
	if whence == C.avd_avseek_file_size() {
		return C.int64_t(C.avd_averr_enosys())
	}
	//read the value from fileIdPointer and convert it to a go uint64
	fileID := uint64(*(*C.uint64_t)(fileIDPtr))

	rawStream, ok := fileRegistry.Load(fileID)
	if !ok {
		slog.Error(
			"seeking stream: file id for requested stream not found in registry",
			"fileID",
			fileID,
		)
		return C.int64_t(C.avd_averr_bug())
	}
	stream, ok := rawStream.(io.ReadSeekCloser)
	if !ok {
		slog.Error(
			"seeking stream: item in registry isnt an io.ReadSeekCloser... somehow, bro how even?! anyway it was a ",
			"type",
			reflect.TypeOf(rawStream).String(),
			"fileID",
			fileID,
		)
		return C.int64_t(C.avd_averr_bug())
	}

	offsetFromStart, err := stream.Seek(int64(offset), int(whence))
	if err != nil {
		slog.Error("seeking stream", "error", err)
		return C.int64_t(C.avd_averr_eio())
	}

	return C.int64_t(offsetFromStart)
}
