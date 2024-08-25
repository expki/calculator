import pako from 'pako';

// WIP, for some reason it doesn't like decoding golang compression output
export function Decompress(input: Uint8Array): [Uint8Array, Error | undefined] {
    const decoder = new pako.Inflate({
    });    
    decoder.push(input, true);
    let err: Error | undefined = undefined;
    switch (decoder.err) {
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
    let result: Uint8Array; 
    if (typeof decoder.result === 'string') {
        console.warn('decompressed to string value');
        result = new TextEncoder().encode(decoder.result);
    } else {
        result = decoder.result;
    }
    return [result, err];
}
