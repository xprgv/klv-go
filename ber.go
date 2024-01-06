package klv

import (
	"errors"
	"io"
)

// BER - basic encoding rules
// Algorithm:
// If the first byte in the length field does not have the high bit set (0x80),
// then that single byte represents an integer between 0 and 127 and indicates
// the number of Value bytes that immediately follows. If the high bit is set,
// then the lower seven bits indicate how many bytes follow that themselves make up
// a length field.

// For example if the first byte of a BER length field is binary 10000010,
// that would indicate that the next two bytes make up an integer that then indicates
// how many Value bytes follow. Therefore a total of three bytes were taken up to specify a length.

var (
	ErrEmptyInput           = errors.New("empty input")
	ErrIncorrectInputLength = errors.New("incorrect input length")
)

func BerEncodeChunk(length int) []byte {
	if length < 128 {
		return []byte{byte(length)}
	}

	var (
		berLength uint8 = 0
		i         uint8 = 0

		delimeter = 256
		ber       = []byte{}
	)

	for {
		berLength++

		if length < delimeter {
			ber = append(ber, 0b10000000|byte(berLength))

			for i = 0; i < berLength; i++ {
				shift := 8 * (berLength - i - 1)
				ber = append(ber, byte(length>>shift))
			}

			return ber
		}

		delimeter *= 256
	}
}

func BerEncode(w io.Writer, length uint64) error {
	if length < 128 {
		_, err := w.Write([]byte{byte(length)})
		return err
	}

	var (
		berLength uint8 = 0
		i         uint8 = 0

		delimeter uint64 = 256
		ber       []byte = []byte{}
	)

	for {
		berLength++

		if length < delimeter {
			ber = append(ber, 0b10000000|byte(berLength))

			for i = 0; i < berLength; i++ {
				shift := 8 * (berLength - i - 1)
				ber = append(ber, byte(length>>shift))
			}

			_, err := w.Write(ber)
			return err
		}

		delimeter *= 256
	}
}

func BerDecodeChunk(input []byte) (uint64, error) {
	inputLength := len(input)
	if inputLength == 0 {
		return 0, ErrEmptyInput
	}

	if input[0] < 128 {
		return uint64(input[0]), nil
	}

	berCodeLength := uint8(input[0] & 127)

	if berCodeLength+1 != uint8(inputLength) {
		return 0, ErrIncorrectInputLength
	}

	var (
		i     uint8  = 0
		value uint64 = 0
	)

	for i = 1; i < uint8(inputLength); i++ {
		shift := (berCodeLength - i) * 8
		shifted := uint64(input[i]) << shift
		value += shifted
	}

	return value, nil
}

func BerDecode(r io.Reader) (uint64, []byte, error) {
	buf := make([]byte, 1)

	if _, err := r.Read(buf); err != nil {
		return 0, []byte{}, err
	}

	if buf[0]>>7 == 0b00000000 {
		return uint64(buf[0]), buf, nil
	}

	var (
		berCode   []byte = []byte{buf[0]}
		berLength uint8  = uint8(buf[0] & 127)
		i         uint8  = 0
		value     uint64 = 0
	)

	for i = 0; i < berLength; i++ {
		if _, err := r.Read(buf); err != nil {
			return 0, berCode, err
		}

		berCode = append(berCode, buf[0])
		shift := (berLength - i - 1) * 8
		shifted := uint64(buf[0]) << shift
		value += shifted
	}

	return value, berCode, nil
}
