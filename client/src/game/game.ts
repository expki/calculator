import { Decode } from './decoder';
import { renderCalculator } from './calculator';
import type * as state from '../types/state';

const target_tick_rate = 1000 / 60
let ctx: CanvasRenderingContext2D;
let canvas: HTMLCanvasElement;
let sharedBuffer: SharedArrayBuffer;

export function game(inCtx: CanvasRenderingContext2D, inCanvas: HTMLCanvasElement, inSharedBuffer: SharedArrayBuffer) {
    ctx = inCtx;
    canvas = inCanvas;
    sharedBuffer = inSharedBuffer;
    setTimeout(renderLoop, 500);
}

let n = 0;
const pref = new Array<number>(60).fill(0);

async function renderLoop(): Promise<void> {
    try {
        n++;
        const start = performance.now();
        // Load game state
        const sharedBytes = new Uint8Array(sharedBuffer);
        const bytes = new Uint8Array(sharedBytes.length);
        bytes.set(sharedBytes);
        const [state, _ ]: [state.StateExt, number] = Decode<state.State>(bytes);
        state.CpuRender = pref.reduce((a, b) => a + b, 0) / 60;
        const xCenter = canvas.width / 2;
        const yCenter = canvas.height / 2

        // Clear canvas
        ctx.reset();
        ctx.fillStyle = "black";
        ctx.fillRect(0, 0, canvas.width, canvas.height);

        // Draw FPS in top left corner
        ctx.fillStyle = "lightgreen";
        ctx.font = "20px Verdana";
        ctx.fillText(`${state.CpuLogic >= 0.1 ? "" : " "}${(100 * state.CpuRender).toFixed(2)}%`, 10, 25);
        ctx.fillStyle = "lightblue";
        ctx.fillText(`${state.CpuLogic >= 0.1 ? "" : " "}${(100 * state.CpuLogic).toFixed(2)}%`, 10, 50);

        // Draw calculator
        renderCalculator(ctx, canvas, state.Global.Calculator, xCenter, yCenter);

        // Calculate render time
        const end = performance.now() - start;
        pref[n%60] = end / target_tick_rate;

        if (n % 120 === 0) {
            console.log(state);
        }
    } catch(e) {
        console.error(e);
    } finally {
        // Request next frame
        requestAnimationFrame(renderLoop);
    }
}
