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
			_, err := w.Write([]byte(e))
			if err != nil {
				return fmt.Errorf("failed to encode data type: string : %w", err)
			}
		}
	}

	return nil
}
