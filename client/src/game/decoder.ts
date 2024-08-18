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
export function Decode<T = any>(binary: Uint8Array): [T, number] {
    let output: any = undefined;
    let outputLength = 0;
    switch (binary[0]) {
        case Type.null:
            output = null;
            outputLength = 1;
            break;
        case Type.bool: {
            output = Boolean(binary[1]);
            outputLength = 2;
            break;
        }
        case Type.uint8: {
            output = binary[1];
            outputLength = 2;
            break;
        }
        case Type.int32: {
            output = new DataView(binary.buffer).getInt32(1, true);
            outputLength = 5;
            break;
        }
        case Type.int64: {
            output = new DataView(binary.buffer).getBigInt64(1, true);
            outputLength = 9;
            break;
        }
        case Type.float32: {
            output = new DataView(binary.buffer).getFloat32(1, true);
            outputLength = 5;
            break;
        }
        case Type.buffer: {
            const length = new DataView(binary.buffer).getUint32(1, true);
            output = binary.slice(5, 5 + length);
            outputLength = 5 + length;
            break;
        }
        case Type.string: {
            const length = new DataView(binary.buffer).getUint32(1, true);
            output = new TextDecoder().decode(binary.slice(5, 5 + length));
            outputLength = 5 + length;
            break;
        }
        case Type.array: {
            const length = new DataView(binary.buffer).getUint32(1, true);
            const array = new Array(length);
            let offset = 5;
            let valueLength = 0;
            for (let idx = 0; idx < length; idx++) {
                [array[idx], valueLength] = Decode(binary.slice(offset));
                if (valueLength == 0) {
                    throw new Error(`unable to decode array item ${idx}`);
                }
                offset += valueLength;
            }
            output = array;
            outputLength = offset;
            break;
        }
        case Type.object: {
            const obj: Record<string, unknown> = {};
            const count = new DataView(binary.buffer).getUint32(1, true);
            let offset = 5;
            let valueLength = 0;
            for (let idx = 0; idx < count; idx++) {
                const keyLength = new DataView(binary.buffer).getUint32(offset, true);
                offset += 4;
                const key = new TextDecoder().decode(binary.slice(offset, offset + keyLength));
                offset += keyLength;
                [obj[key], valueLength] = Decode(binary.slice(offset));
                if (valueLength == 0) {
                    throw new Error(`unable to decode object item ${key}`);
                }
                offset += valueLength;
            }
            output = obj;
            outputLength = offset;
            break;
        }
        default:
            output = null;
            outputLength = 0;
            break;
    }
    return [output as T, outputLength];
}