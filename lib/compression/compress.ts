import pako from 'pako';

export function Compress(input: Uint8Array): [Uint8Array, Error | undefined] {
    const encoder = new pako.Deflate({
        level: 1,
        strategy: 0,
        header: {
            text: false,
            hcrc: false,
        },
    });
    encoder.push(input, true);
    let err: Error | undefined = undefined;
    switch (encoder.err) {
    case pako.constants.Z_STREAM_END:
        err = new Error('compress error: Z_STREAM_END');
        break;
    case pako.constants.Z_NEED_DICT:
        err = new Error('compress error: Z_NEED_DICT');
        break;
    case pako.constants.Z_ERRNO:
        err = new Error('compress error: Z_ERRNO');
        break;
    case pako.constants.Z_STREAM_ERROR:
        err = new Error('compress error: Z_STREAM_ERROR');
        break;
    case pako.constants.Z_DATA_ERROR:
        err = new Error('compress error: Z_DATA_ERROR');
        break;
    case pako.constants.Z_BUF_ERROR:
        err = new Error('compress error: Z_BUF_ERROR');
        break;
    }
    return [encoder.result, err];
}
