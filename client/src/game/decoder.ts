enum Type {
    null = 0,
    bool = 1,
    uint8 = 2,
    int32 = 3,
    int64 = 4,
    float32 = 5,
    buffer = 6,
    string = 7,
    array = 8,
    object = 9,
}

// Decode implements a custom decoding scheme for the wasm worker
export function Decode<T = any>(binary: Uint8Array): T {
    let output: any = undefined;
    switch (binary[0]) {
        case Type.bool: {
            output = Boolean(binary[1]);
            break;
        }
        case Type.uint8: {
            output = binary[1];
            break;
        }
        case Type.int32: {
            output = new DataView(binary.buffer).getInt32(1, true);
            break;
        }
        case Type.int64: {
            output = new DataView(binary.buffer).getBigInt64(1, true);
            break;
        }
        case Type.float32: {
            output = new DataView(binary.buffer).getFloat32(1, true);
            break;
        }
        case Type.buffer: {
            const length = new DataView(binary.buffer).getUint32(1, true);
            output = binary.slice(5, 5 + length);
            break;
        }
        case Type.string: {
            const length = new DataView(binary.buffer).getUint32(1, true);
            output = new TextDecoder().decode(binary.slice(5, 5 + length));
            break;
        }
        case Type.array: {
            const array = new Array(new DataView(binary.buffer).getUint32(1, true));
            let offset = 5;
            for (let idx = 0; idx < array.length; idx++) {
                const valueLength = new DataView(binary.buffer).getUint32(offset, true);
                array[idx] = Decode(binary.slice(offset + 4, offset + 4 + valueLength));
                offset += valueLength + 4;
            }
            output = array;
            break;
        }
        case Type.object: {
            const obj: Record<string, unknown> = {};
            const maxLength = new DataView(binary.buffer).getUint32(1, true);
            let offset = 5;
            while (offset < maxLength + 5) {
                const keyLength = new DataView(binary.buffer).getUint32(offset + 1, true);
                const key = new TextDecoder().decode(binary.slice(offset + 5, offset + 5 + keyLength));
                const valueLength = new DataView(binary.buffer).getUint32(offset + 5 + keyLength, true);
                const value = Decode(binary.slice(offset + 9 + keyLength, offset + 9 + keyLength + valueLength));
                obj[key] = value;
                offset += 9 + keyLength + valueLength;
            }
            output = obj;
            break;
        }
        default:
            break;
    }
    return output as T;
}