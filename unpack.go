package binary

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"encoding/binary"

	"github.com/fatih/structs"
)

var Debug bool

// Unpack takes an arbitrary number of destination structs to decode the byte slice
func Unpack(source []byte, endianness binary.ByteOrder, destinations ...interface{}) (int, error) {
	bitString := func(x uint32) string {
		bits := strconv.FormatUint(uint64(x), 2)
		for len(bits) < 32 {
			bits = "0" + bits
		}
		return bits
	}

	unpack := func(r io.Reader, destination interface{}) (int, error) {
		switch destination.(type) {
		case *uint64:
			return 8, binary.Read(r, endianness, destination)
		case *uint32:
			return 4, binary.Read(r, endianness, destination)
		case *uint16:
			return 2, binary.Read(r, endianness, destination)
		case *uint8:
			return 1, binary.Read(r, endianness, destination)
		default:
			bits, err := Size(destination)
			if err != nil {
				return 0, err
			}
			if (bits%8) != 0 || bits < 8 {
				return 0, errors.New("We can only unpack bits within full bytes (multiple of 8)")
			}
			switch bits {
			case 32:
				// struct with bits tags
				var reading uint32
				err = binary.Read(r, endianness, &reading)
				if err != nil {
					return 0, err
				}

				if Debug {
					fmt.Println(bitString(reading), fmt.Sprintf("read value (%x)", reading))
				}

				dest := structs.New(destination)
				names := structs.Names(destination)
				bitOffset := 0
				for _, name := range names {
					bitSize, _ := strconv.Atoi(dest.Field(name).Tag("bits"))
					bitMask := uint32((1 << uint32(bitSize)) - 1) // bitSize=3, 1<<3 (1000) - 1 = 111

					// in little endian we read the bits from the right (lsb)
					bitShift := uint32(bitOffset)
					// and in big endian from the left (msb)
					if endianness.String() == "BigEndian" {
						bitShift = uint32(32-bitOffset) - uint32(bitSize)
					}

					bitValue := (reading & (bitMask << bitShift)) >> bitShift

					dest.Field(name).Set(bitValue)

					bitOffset += bitSize
					bitMaskString := bitString(bitMask << bitShift)
					if Debug {
						fmt.Println(bitMaskString, name, "bits", bitSize, "mask", strconv.FormatUint(uint64(bitMask), 2), "bitshift", bitShift)
					}
				}
				return 4, nil
			}
			return 0, errors.New(fmt.Sprintf("Unsupported type/number of bits to decode: %d", bits))
		}
	}

	var offset int
	reader := bytes.NewReader(source)

	for idx, _ := range destinations {
		n, err := unpack(reader, destinations[idx])
		if err != nil {
			return 0, err
		}
		offset += n
	}
	return offset, nil
}

// Size returns size of struct in bits, or error
func Size(destination interface{}) (int, error) {
	dest := structs.New(destination)
	names := structs.Names(destination)
	bits := 0
	for _, name := range names {
		b, err := strconv.Atoi(dest.Field(name).Tag("bits"))
		if err != nil {
			return 0, err
		}
		bits += b
	}
	return bits, nil
}
