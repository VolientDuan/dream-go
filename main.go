package main

// #include <stdlib.h>
// #include <sys/types.h>
// static void callLogger(void *func, const char *msg)
// {
// 	((void(*)(const char *))func)(msg);
// }
import "C"

import (
	"unsafe"

	"dream/httpService"
)

var loggerFunc unsafe.Pointer

//export setLogger
func setLogger(loggerFn uintptr) {
	loggerFunc = unsafe.Pointer(loggerFn)
}

//export startFileServer
func startFileServer(p int, path string) {
	httpService.StartFileServer(p, path)
}

//export startManhuaServer
func startManhuaServer(p int) {
	httpService.StartManhuaService(p)
}

func main() {
	// startManhuaServer(8081)
}
