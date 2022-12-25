package server

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

type (
	PHeader struct {
		ID      uint16
		Version uint8
	}

	PMessage[BodyType any] struct {
		Header PHeader
		Body   BodyType
	}
)

const (
	messageStart = uint16(0xAAAA)
	messageEnd   = uint16(0xFFFF)

	// Missing start of message
	MissingSOM = "not at the start of a valid message"
	// Missing end of message
	MissingEOM = "there are bytes left on the message"
)

func (h *PHeader) FromStream(buf io.Reader) (err error) {
	var flag uint16
	err = binary.Read(buf, binary.LittleEndian, &flag)
	if err != nil {
		return
	}
	if flag != messageStart {
		return errors.New(MissingSOM)
	}

	return binary.Read(buf, binary.LittleEndian, h)
}

func (h *PHeader) ToStream(buf io.Writer) (err error) {
	err = binary.Write(buf, binary.LittleEndian, messageStart)
	if err != nil {
		return
	}

	return binary.Write(buf, binary.LittleEndian, h)
}

func (m *PMessage[BodyType]) FromStream(buf io.Reader) (err error) {
	// Process header
	err = m.Header.FromStream(buf)
	if err != nil {
		return
	}

	// Process body
	err = unserializeGeneric[BodyType](buf, &m.Body)
	if err != nil {
		return
	}

	// Process footer
	var flag uint16
	err = binary.Read(buf, binary.LittleEndian, &flag)
	if err != nil {
		return
	}
	if flag != messageEnd {
		return errors.New(MissingEOM)
	}

	return
}

func (m *PMessage[BodyType]) ToStream(buf io.Writer) (err error) {
	// Write header
	err = m.Header.ToStream(buf)
	if err != nil {
		return
	}

	// Write body
	err = serializeGeneric[BodyType](buf, m.Body)
	if err != nil {
		return
	}

	// Write footer and return
	return binary.Write(buf, binary.LittleEndian, messageEnd)
}

func serializeGeneric[T any](buf io.Writer, data T) (err error) {
	var face = reflect.ValueOf(data)

	for i := 0; i < face.NumField(); i++ {
		var (
			fld = face.Field(i)
			val = fld.Interface()
		)

		var bSize = binary.Size(val)
		if bSize >= 0 {
			// Fixed size type
			if fld.Kind() == reflect.Slice {
				err = binary.Write(buf, binary.LittleEndian, uint32(bSize))
				if err != nil {
					return
				}
			}

			err = binary.Write(buf, binary.LittleEndian, val)
			if err != nil {
				return
			}
		} else {
			// Annoying type
			switch fld.Kind() {
			case reflect.String:
				var size = fld.Len()
				err = binary.Write(buf, binary.LittleEndian, uint16(size))
				if err != nil {
					return
				}
				for j := 0; j < size; j++ {
					err = binary.Write(buf, binary.LittleEndian, fld.Index(j).Interface())
					if err != nil {
						return
					}
				}
			default:
				return errors.New("invalid type " + fld.Type().String())
			}
		}

	}

	return nil
}

func unserializeGeneric[T any](buf io.Reader, data *T) (err error) {
	var face = reflect.Indirect(reflect.ValueOf(data))

	for i := 0; i < face.NumField(); i++ {
		var (
			fld = face.Field(i)
			val = fld.Interface()
		)

		var bSize = binary.Size(val)
		if bSize >= 0 {
			// Fixed size type
			if fld.Kind() == reflect.Slice {
				var tmp uint32
				err = binary.Read(buf, binary.LittleEndian, &tmp)
				if err != nil {
					return
				}
				var size = int(tmp)
				fld.Set(reflect.MakeSlice(fld.Type(), size, size))
			}

			err = binary.Read(buf, binary.LittleEndian, fld.Addr().Interface())
			if err != nil {
				return
			}
		} else {
			// Annoying type
			switch fld.Kind() {
			case reflect.String:
				var tmp uint16
				err = binary.Read(buf, binary.LittleEndian, &tmp)
				if err != nil {
					return
				}
				var size = int(tmp)
				var strBuffer = make([]byte, size)
				for j := 0; j < size; j++ {
					err = binary.Read(buf, binary.LittleEndian, &strBuffer[j])
					if err != nil {
						return
					}
				}
				fld.SetString(string(strBuffer))
			default:
				return errors.New("invalid type " + fld.Type().String())
			}
		}
	}

	return nil
}
