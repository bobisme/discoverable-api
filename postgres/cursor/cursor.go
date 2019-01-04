package cursor

import (
	"bytes"
	"encoding/base64"

	"github.com/tinylib/msgp/msgp"
)

func Encode(e msgp.Encodable) (string, error) {
	buf := new(bytes.Buffer)
	err := msgp.Encode(buf, e)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString([]byte(buf.Bytes())), nil
}

func Decode(cursor string, d msgp.Decodable) error {
	b, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	err = msgp.Decode(buf, d)
	if err != nil {
		return err
	}
	return nil
}
