import lib from '@lib/index';
import RenderCalculator from './calculator';
import RenderPreformance from './pref';
import RenderCursor from './cursor';
import type { LocalState } from '../types/state';
import type { ButtonMap } from '../types/types';

declare global {
    var buttonMap: ButtonMap | undefined;
}

const target_tick_rate = 1000 / 60;
let ctx: CanvasRenderingContext2D;
let canvas: HTMLCanvasElement;
let sharedBuffer: SharedArrayBuffer;

export default function Game(inCtx: CanvasRenderingContext2D, inCanvas: HTMLCanvasElement, inSharedBuffer: SharedArrayBuffer) {
    ctx = inCtx;
    canvas = inCanvas;
    sharedBuffer = inSharedBuffer;
    setTimeout(renderLoop, 500);
}

let n = 0;
let lastFrameTime = performance.now();
const fps = new Array<number>(60).fill(0);
const gpuLoad = new Array<number>(60).fill(0);
const cpuLoad = new Array<number>(60).fill(0);
let prevBytes = new Uint8Array();

function renderLoop(): void {
    try {
        n++;
        const start = performance.now();
        const delta = start - lastFrameTime;
        lastFrameTime = start;

        // Calculate fps
        fps[n % 60] = 1000 / delta;

        // Load game state
        const sharedBytes = new Uint8Array(sharedBuffer);
        const bytes = new Uint8Array(sharedBytes.length);
        bytes.set(sharedBytes);
        const [state, _ ] = lib.encoding.Decode<LocalState>(bytes);

        // Draw calculator
        RenderCalculator(ctx, canvas, state.State);

        // Draw FPS in top left corner
        const stableFps = fps.reduce((a, b) => a + b, 0) / 60;
        const stableGpuLoad = gpuLoad.reduce((a, b) => a + b, 0) / 60;
        const stableCpuLoad = cpuLoad.reduce((a, b) => a + b, 0) / 60;
        RenderPreformance(ctx, stableFps, stableGpuLoad, stableCpuLoad);

        // Draw cursor
        (state.State.Members ?? []).forEach((member) => RenderCursor(ctx, canvas, member));

        // Calculate render time
        const end = performance.now() - start;
        gpuLoad[n%60] = end / target_tick_rate;
        cpuLoad[n%60] = state.CpuLoad ?? 1;

        // Log state
        if (!bytesEqual(bytes, prevBytes)) {
            prevBytes = new Uint8Array(bytes.length);
            prevBytes.set(bytes);
            console.debug(state);
        }
    } catch(e: unknown) {
        console.error(e);
    } finally {
        // Request next frame
        requestAnimationFrame(renderLoop);
    }
}

function bytesEqual(a: Uint8Array, b: Uint8Array): boolean {
    if (a.length !== b.length) return false;
    let result = 0;
    for (let i = 0; i < a.length; i++) {
        result |= a[i] ^ b[i];
    }
    return result === 0;
}