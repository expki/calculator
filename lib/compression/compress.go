package compression

import "github.com/klauspost/compress/zstd"

var encoder *zstd.Encoder = func() *zstd.Encoder {
	encoder, err := zstd.NewWriter(
		nil,
		zstd.WithEncoderLevel(zstd.SpeedFastest),
		zstd.WithSingleSegment(true),
		zstd.WithEncoderCRC(false),
		zstd.WithEncoderConcurrency(1),
		zstd.WithEncoderPadding(1),
		zstd.WithNoEntropyCompression(true),
	)
	if err != nil {
		panic(err)
	}
	return encoder
}()

func Compress(in []byte) (out []byte) {
	out = encoder.EncodeAll(in, out)
	return out
}
