import '../../public/wasm_exec';
import * as enums from '../types/enums';
import type * as types from '../types/types';

declare global {
    var sharedArray: Uint8Array | undefined;
}

declare class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): void;
}

declare function handleInput(payload: types.PayloadInput): void;

let sharedBuffer: SharedArrayBuffer | undefined;
let n = 0;

self.onmessage = async (event: MessageEvent<types.Payload<any>>) => {
    n++;
    switch (event.data.kind) {
        case enums.PayloadKind.wasm: {
            const payload: types.PayloadWasm = event.data.payload;
            sharedBuffer = payload.pipe;
            global.sharedArray = new Uint8Array(sharedBuffer);
            const go = new Go();
            const result = await WebAssembly.instantiate(payload.wasm, go.importObject);
            go.run(result.instance);
            return;
        }
        case enums.PayloadKind.input: {
            const payload: types.PayloadInput = event.data.payload;
            try {
                handleInput(payload);
            } catch (_) {}
            return;
        }
        default:
            console.error("Unknown payload kind:", event.data.kind);
            return;
    }
}
