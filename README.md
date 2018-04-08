# binary

A helper to decode binary data into structures, with bit-size accuracy.

In C, you can declare structures like so:

~~~
typedef struct _SEQ_WR_CTL_D1_FORMAT
{
	uint32_t DAT_DLY : 4;
	uint32_t DQS_DLY : 4;
	uint32_t DQS_XTR : 1;
	uint32_t DAT_2Y_DLY : 1;
	uint32_t ADR_2Y_DLY : 1;
	uint32_t CMD_2Y_DLY : 1;
	uint32_t OEN_DLY : 4;
	uint32_t OEN_EXT : 4;
	uint32_t OEN_SEL : 2;
	uint32_t Pad0 : 2;
	uint32_t ODT_DLY : 4;
	uint32_t ODT_EXT : 1;
	uint32_t ADR_DLY : 1;
	uint32_t CMD_DLY : 1;
	uint32_t Pad1 : 1;
} SEQ_WR_CTL_D1_FORMAT;
~~~

The notation with the `:` (colon) operator specifies how many bits each field represents; We can
convert this structure to Go with the help of tags:

~~~
type SEQ_WR_CTL_D1_FORMAT struct {
	DAT_DLY uint32 `bits:"4"`
	DQS_DLY uint32 `bits:"4"`
	DQS_XTR    uint32 `bits:"1"`
	DAT_2Y_DLY uint32 `bits:"1"`
	ADR_2Y_DLY uint32 `bits:"1"`
	CMD_2Y_DLY uint32 `bits:"1"`
	OEN_DLY    uint32 `bits:"4"`
	OEN_EXT uint32 `bits:"4"`
	OEN_SEL uint32 `bits:"2"`
	Pad0    uint32 `bits:"2"`
	ODT_DLY uint32 `bits:"4"`
	ODT_EXT uint32 `bits:"1"`
	ADR_DLY uint32 `bits:"1"`
	CMD_DLY uint32 `bits:"1"`
	Pad1    uint32 `bits:"1"`
}
~~~

In C, you can create a typed pointer that points to a byte array with this representation.
In Go, you'd have to resort to the `unsafe` package for something like that. I use bitwise
operators on scanned uint32 values in order to fill the above struct with appropriate data.

Due to the limitations of the underlying `fatih/structs`, the Unpack function takes a
variadic argument to the fields you want to decode in the order you want to decode
them. In the future, a more generic function might be available, where you would pass
just a single struct with any number of nested structs, and it would linearelly decode
the available data into the structs.

## API

API is subject to change. While reasonable efforts will be made to keep it as it is,
there might be circumstances where APIs will need to be changed.

## Usage

See [godoc page](https://godoc.org/titpetric/binary). Generally, it can be something like:

~~~
	n, err := b.Unpack(strap, binary.LittleEndian,
		&result.SEQ_WR_CTL_D1,
		&result.SEQ_WR_CTL_2,
		&result.SEQ_PMG_TIMING,
		&result.SEQ_RAS_TIMING,
		&result.SEQ_CAS_TIMING,
		&result.SEQ_MISC_TIMING,
		&result.SEQ_MISC_TIMING2,
		&result.SEQ_MISC1,
		&result.SEQ_MISC3,
		&result.SEQ_MISC8,
		&result.ARB_DRAM_TIMING,
		&result.ARB_DRAM_TIMING2,
	)
~~~

The variadic parameters may currently be any uint type, or a struct which reads 32 bits (uint32, 4 bytes).
Support for structs that align to other uint types is possible, just takes a bit of copy pasting. As
I didn't need it, I didn't write it, but feel free to add it on if you want.

## Other

- [x] support uint8, uint16, uint32, uint64 natively (without bitmasks)
- [x] support uint32 struct with bitmasks
- [ ] support uint8, uint16, uint64 with bitmasks
- [ ] support any number of fields with encapsulating struct

## License

Written by [@TitPetric](https://twitter.com/TitPetric) and licensed under the permissive [WTFPL](http://www.wtfpl.net/txt/copying/).
