package compression

import "github.com/klauspost/compress/zstd"

var decoder *zstd.Decoder = func() *zstd.Decoder {
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		panic(err)
	}
	return decoder
}()

func Decompress(in []byte) (out []byte, err error) {
	_, err = decoder.DecodeAll(in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
