import '../../public/wasm_exec';
import * as enums from '../types/enums';
import type { Payload, PayloadInput, PayloadWasm } from '../types/types';

declare global {
    var sharedArray: Uint8Array | undefined;
}

declare class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): void;
}

declare function handleInput(payload: PayloadInput): void;

let sharedBuffer: SharedArrayBuffer | undefined;
let n = 0;

self.onmessage = async (event: MessageEvent<Payload<any>>) => {
    n++;
    switch (event.data.kind) {
        case enums.PayloadKind.wasm: {
            const payload: PayloadWasm = event.data.payload;
            sharedBuffer = payload.pipe;
            global.sharedArray = new Uint8Array(sharedBuffer);
            const go = new Go();
            const result = await WebAssembly.instantiate(payload.wasm, go.importObject);
            go.run(result.instance);
            return;
        }
        case enums.PayloadKind.input: {
            const payload: PayloadInput = event.data.payload;
            try {
                handleInput(payload);
            } catch {}
            return;
        }
        default:
            console.error("Unknown payload kind:", event.data.kind);
            return;
    }
}
