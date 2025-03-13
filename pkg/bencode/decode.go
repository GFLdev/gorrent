package bencode

import (
	"errors"
	"fmt"
	"strconv"
)

type bReader struct {
	data []byte
	pos  int
}

func newReader(data []byte) *bReader {
	return &bReader{data: data, pos: 0}
}

func (r *bReader) readUntil(delim byte) ([]byte, error) {
	start := r.pos
	for r.pos < len(r.data) {
		if r.data[r.pos] == delim {
			r.pos++
			return r.data[start:r.pos], nil
		}
		r.pos++
	}
	return nil, errors.New("missing delimiter")
}

func (r *bReader) decodeInt() (int, error) {
	s, err := r.readUntil('e')
	if err != nil {
		return 0, fmt.Errorf("could not decode integer: %w", err)
	}

	if len(s) == 2 {
		return 0, fmt.Errorf("bencode integer '%s' is empty", s)
	}

	integer, err := strconv.Atoi(string(s[1 : len(s)-1]))
	if err != nil {
		return 0, fmt.Errorf("could not decode integer '%s': %w", s, err)
	}
	return integer, nil
}

func (r *bReader) decodeString() (string, error) {
	lenStr, err := r.readUntil(':')
	if err != nil {
		return "", fmt.Errorf("could not decode string: %w", err)
	}

	length, err := strconv.Atoi(string(lenStr[:len(lenStr)-1]))
	if err != nil {
		return "", fmt.Errorf("could not decode string '%s' length: %w", r.data[:r.pos], err)
	}
	if length < 0 {
		return "", fmt.Errorf("string length is negative: %d", length)
	}

	end := r.pos + length
	if end > len(r.data) {
		return "", fmt.Errorf("string length exceeds input length: expected %d bytes, got %d bytes", end, len(r.data))
	}

	start := r.pos
	r.pos = end
	return string(r.data[start:end]), nil
}

func (r *bReader) decodeList() ([]interface{}, error) {
	list := make([]interface{}, 0)
	r.pos++

	for {
		if r.pos >= len(r.data) {
			return nil, errors.New("eol reached before end of list")
		}

		if r.data[r.pos] == 'e' {
			break
		}

		val, err := r.decodeElement()
		if err != nil {
			return nil, err
		}

		list = append(list, val)
	}

	return list, nil
}

func (r *bReader) decodeDict() (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	r.pos++

	for {
		if r.pos >= len(r.data) {
			return nil, errors.New("eol reached before end of dict")
		}

		if r.data[r.pos] == 'e' {
			break
		}

		key, err := r.decodeString()
		if err != nil {
			return nil, fmt.Errorf("could not decode dict key: %w", err)
		}

		val, err := r.decodeElement()
		if err != nil {
			return nil, fmt.Errorf("could not decode dict value: %w", err)
		}

		dict[key] = val
	}

	return dict, nil
}

func (r *bReader) decodeElement() (interface{}, error) {
	if r.data[r.pos:] == nil || len(r.data[r.pos:]) == 0 {
		return nil, errors.New("entry is empty")
	}

	switch firstChar := r.data[r.pos]; {
	case firstChar == 'i':
		return r.decodeInt()
	case firstChar == 'l':
		return r.decodeList()
	case firstChar == 'd':
		return r.decodeDict()
	case firstChar >= '0' && firstChar <= '9':
		return r.decodeString()
	default:
		return nil, errors.New("invalid type encountered: character not 'i', 'l', 'd', or '0'-'9'")
	}
}

func Decode(s []byte) (interface{}, error) {
	val, err := newReader(s).decodeElement()
	if err != nil {
		return nil, fmt.Errorf("could not decode bencode: %w", err)
	}
	return val, nil
}
