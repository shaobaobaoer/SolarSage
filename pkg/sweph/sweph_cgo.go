package sweph

// Link against the pre-built Swiss Ephemeris static library

/*
#cgo LDFLAGS: ${SRCDIR}/../../third_party/swisseph/libswe.a -lm
*/
import "C"
