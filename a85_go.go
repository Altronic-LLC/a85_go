package a85_go
// Altronic a85 encoding/decoding library for Golang
// Similar to z85 but no need to pre-pad before encoding/decoding.
// Support for runs of \x00 and \xff to be encoded as a single character to be added later

import (
	"fmt"
)

var (
	encoder = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.-:+=^!/*?&<>()[]{}@%$#"
	decoder = []byte{
		0x00, 0x44, 0x00, 0x54, 0x53, 0x52, 0x48, 0x00,
		0x4B, 0x4C, 0x46, 0x41, 0x00, 0x3F, 0x3E, 0x45,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x40, 0x00, 0x49, 0x42, 0x4A, 0x47,
		0x51, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A,
		0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32,
		0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A,
		0x3B, 0x3C, 0x3D, 0x4D, 0x00, 0x4E, 0x43, 0x00,
		0x00, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
		0x21, 0x22, 0x23, 0x4F, 0x00, 0x50, 0x00, 0x00}
	sizeDecoder = len(decoder)
)

//Encode a byte array to a string. No need to pad the string. Padding will be done automatically
func Encode(slice []byte, sliceLength int) []byte {
	rv := make([]byte, (sliceLength*5/4)+4)
	dst := rv
	digits := 0
	cur := 0
	dstPtr := 0
	for lp := 0; lp < sliceLength; lp++ {
		b := slice[lp]
		cur *= 1 << 8
		cur += int(b)
		digits++
		if digits%4 != 0 {
			continue
		}

		for bytes := 4; bytes >= 0; bytes-- {
			r := cur % 85
			dst[dstPtr+bytes] = encoder[r]
			cur = (cur - r) / 85
		}
		dstPtr += 5
		cur = 0
		digits = 0

	}


	if digits > 0 {
		padding := 4 - digits
		fmt.Println("padding: ", padding, "cur: ", cur)
		for lp := padding; lp > 0; lp-- {
			cur *= 1 << 8
		}

		for bytes := 4; bytes >= 0; bytes-- {
			r := cur % 85
			dst[dstPtr+bytes] = encoder[r]
			cur = (cur - r) / 85
		}
		dstPtr += (5 - padding)
		fmt.Println("dstPtr: ", dstPtr)
	}
	dst[dstPtr] = '\x00'
	return dst[0:dstPtr]
}

// Decode string to a slice
// String length does not need to be a multiple of 5 if it was created with the corresponding
//   encode function
func Decode(str string) []byte {
	size := len(str)
	rv := make([]byte, (size*4/5)+10)
	dst := rv
	digits := 0
	cur := 0
	dstPtr := 0
	for lp := 0; lp < size; lp++ {
		cur *= 85
		c := str[lp]
		decByte := (c - 32) & 127
		cur += int(decoder[decByte])
		digits++

		if digits%5 != 0 {
			continue
		}

		dst[dstPtr] = byte(cur >> 24)
		dstPtr += 1
		dst[dstPtr] = byte(cur >> 16 & 255)
		dstPtr += 1
		dst[dstPtr] = byte(cur >> 8 & 255)
		dstPtr += 1
		dst[dstPtr] = byte(cur & 255)
		dstPtr += 1
		digits = 0
		cur = 0
	}


	if digits != 0 {
		padding := 5 - digits
		fmt.Println("padding: ", padding, "cur: ", cur)
		for i := 0; i < padding; i++ {
			cur *= 85
			cur += 85 - 1
		}

		len := padding - 1
		for i := 3; i > len; i-- {
			dst[dstPtr] = byte((cur >> (i * 8)) & 0xFF)
			dstPtr += 1
		}

	}

	return dst[0:dstPtr]
}
