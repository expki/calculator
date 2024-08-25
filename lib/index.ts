import { Decode, DecodeWithCompression } from './encoding/decoder';
import { Compress } from './compression/compress';
import { Decompress } from './compression/decompress'; 

export const compression = {
    Compress,
    Decompress,
};

export const encoding = {
    Decode,
    DecodeWithCompression,
};

export default {
    compression,
    encoding,
};
