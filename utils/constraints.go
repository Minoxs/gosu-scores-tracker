package utils

type (
	ConSigned interface {
		int | int8 | int16 | int32 | int64
	}
	ConUnsigned interface {
		uint | uint8 | uint16 | uint32 | uint64
	}
	ConInteger interface {
		ConSigned | ConUnsigned
	}
)
