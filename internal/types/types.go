package types

type IntGH interface {
	int | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
}

type SignIntGH interface {
	int | int8 | int16 | int32 | int64
}

type UIntGH interface {
	uint8 | uint16 | uint32 | uint64
}
