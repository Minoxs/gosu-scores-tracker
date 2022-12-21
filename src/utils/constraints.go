package utils

type (
	ASignedInteger interface {
		int | int8 | int16 | int32 | int64
	}
	AUnsignedInteger interface {
		uint | uint8 | uint16 | uint32 | uint64
	}
	AInteger interface {
		ASignedInteger | AUnsignedInteger
	}
)
