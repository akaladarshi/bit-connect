package messages

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func generateNonce() uint64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
}

// EncodeData encodes the data into the writer
func EncodeData(w io.Writer, elements ...interface{}) error {
	for _, element := range elements {
		switch e := element.(type) {
		case int32:
			var buf = make([]byte, 4)
			binary.LittleEndian.PutUint32(buf, uint32(e))
			_, err := w.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to encode data type: %T : %w", e, err)
			}
		case int64:
			var buf = make([]byte, 8)
			binary.LittleEndian.PutUint64(buf, uint64(e))
			_, err := w.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to encode data type: %T : %w", e, err)
			}
		case uint32:
			var buf = make([]byte, 4)
			binary.LittleEndian.PutUint32(buf, e)
			_, err := w.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to encode data type: %T : %w", e, err)
			}
		case uint64:
			var buf = make([]byte, 8)
			binary.LittleEndian.PutUint64(buf, e)
			_, err := w.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to encode data type: %T : %w", e, err)
			}
		case [4]byte:
			_, err := w.Write(e[:])
			if err != nil {
				return fmt.Errorf("failed to encode data type: checksum : %w", err)
			}
			return nil
		case [CommandSize]byte:
			_, err := w.Write(e[:])
			if err != nil {
				return fmt.Errorf("failed to encode data type: command : %w", err)
			}
		case [16]byte:
			_, err := w.Write(e[:])
			if err != nil {
				return fmt.Errorf("failed to encode data type: IP : %w", err)
			}
		case string:
			var buf = make([]byte, 1)
			buf[0] = byte(len(e)) // write the length of the string
			_, err := w.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to encode data type: string : %w", err)
			}

			_, err = w.Write([]byte(e))
			if err != nil {
				return fmt.Errorf("failed to encode data type: string : %w", err)
			}
		case bool:
			var buf = make([]byte, 1)
			if e {
				buf[0] = 1
			} else {
				buf[0] = 0
			}

			_, err := w.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to encode data type: %T : %w", e, err)
			}
		default:
			return fmt.Errorf("unknown data type: %T", e)
		}
	}

	return nil
}

func DecodeData(r io.Reader, elements ...interface{}) error {
	for _, element := range elements {
		switch e := element.(type) {
		case *uint32:
			var buf = make([]byte, 4)
			_, err := io.ReadFull(r, buf[:])
			if err != nil {
				return fmt.Errorf("failed to decode data type: %T : %w", e, err)
			}

			*e = binary.LittleEndian.Uint32(buf)
		case *uint64:
			var buf = make([]byte, 8)
			_, err := io.ReadFull(r, buf[:])
			if err != nil {
				return fmt.Errorf("failed to decode data type: %T : %w", e, err)
			}

			*e = binary.LittleEndian.Uint64(buf)
		case *int32:
			var buf = make([]byte, 4)
			_, err := io.ReadFull(r, buf[:])
			if err != nil {
				return fmt.Errorf("failed to decode data type: %T : %w", e, err)
			}

			*e = int32(binary.LittleEndian.Uint32(buf))
		case *int64:
			var buf = make([]byte, 8)
			_, err := io.ReadFull(r, buf[:])
			if err != nil {
				return fmt.Errorf("failed to decode data type: %T : %w", e, err)
			}

			*e = int64(binary.LittleEndian.Uint64(buf))
		case *[CommandSize]byte:
			_, err := io.ReadFull(r, e[:])
			if err != nil {
				return fmt.Errorf("failed to decode data type: %T : %w", e, err)
			}
		case *[4]byte:
			_, err := io.ReadFull(r, e[:])
			if err != nil {
				return fmt.Errorf("failed to decode data type: checksum : %w", err)
			}
		case *[16]byte:
			_, err := io.ReadFull(r, e[:])
			if err != nil {
				return fmt.Errorf("failed to decode data type: IP : %w", err)
			}
		case *string:
			buf := make([]byte, 1)
			_, err := io.ReadFull(r, buf)
			if err != nil {
				return fmt.Errorf("failed to decode data length of type: string : %w", err)
			}

			str := make([]byte, buf[0])
			_, err = io.ReadFull(r, str)
			if err != nil {
				return fmt.Errorf("failed to decode data type: string : %w", err)
			}

			*e = string(str)
		case *bool:
			buf := make([]byte, 1)
			_, err := io.ReadFull(r, buf)
			if err != nil {
				return fmt.Errorf("failed to decode data type: %T : %w", e, err)
			}

			if buf[0] == 1 {
				*e = true
			} else {
				*e = false
			}
		default:
			return fmt.Errorf("unknown data type: %T", e)
		}
	}

	return nil
}
