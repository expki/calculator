import * as enums from '../types/enums';
import type * as types from '../types/types';

export async function listenForUserInput(canvas: HTMLCanvasElement, logic: Worker) {
    // Register input events
    const keyState: Record<string, boolean> = {};
    let mouseLeftState: boolean | undefined;
    let xState: number | undefined;
    let yState: number | undefined;
    addEventListener("keydown", (ev) => {
        if (keyState[ev.key] === true) return;
        keyState[ev.key] = true;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: {
                keydown: ev.key,
            },
        };
        logic.postMessage(payload);
    });
    addEventListener("keyup", (ev) => {
        if (keyState[ev.key] === false) return;
        keyState[ev.key] = false;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: {
                keyup: ev.key,
            },
        };
        logic.postMessage(payload);
    });
    addEventListener("mousedown", (ev) => {
        let enabled = false;
        const data: types.PayloadInput = {};
        if (ev.button === 0 && mouseLeftState !== true) {
            enabled = true;
            data.mouseleft = true;
        }
        if (ev.clientX !== xState) {
            enabled = true;
            data.mousex = ev.clientX;
        }
        if (ev.clientY !== yState) {
            enabled = true;
            data.mousey = ev.clientY;
        }
        if (!enabled) return;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    });
    addEventListener("mouseup", (ev) => {
        let enabled = false;
        const data: types.PayloadInput = {};
        if (ev.button === 0 && mouseLeftState !== false) {
            enabled = true;
            data.mouseleft = false;
        }
        if (ev.clientX !== xState) {
            enabled = true;
            data.mousex = ev.clientX;
        }
        if (ev.clientY !== yState) {
            enabled = true;
            data.mousey = ev.clientY;
        }
        if (!enabled) return;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    });
    addEventListener("mousemove", (ev) => {
        let enabled = false;
        const data: types.PayloadInput = {};
        if (ev.clientX !== xState) {
            enabled = true;
            data.mousex = ev.clientX;
        }
        if (ev.clientY !== yState) {
            enabled = true;
            data.mousey = ev.clientY;
        }
        if (!enabled) return;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: {
                mousex: ev.clientX,
                mousey: ev.clientY,
            },
        };
        logic.postMessage(payload);
    });
    const resizeScreen = (width: number, height: number): void => {
        let enabled = false;
        if (canvas.width !== width || canvas.height !== height) {
            enabled = true;
            // Set canvas size to window size
            const dpr = window.devicePixelRatio || 1;
            canvas.width = window.innerWidth * dpr;
            window.innerHeight * dpr;
            canvas.style.width = window.innerWidth + "px";
            canvas.style.height = window.innerHeight + "px";
            canvas.getContext("2d").scale(dpr, dpr);
        }
        if (!enabled) return;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: {
                width: width,
                height: height,
            },
        };
        logic.postMessage(payload);
    }
    let resizeTimeout: NodeJS.Timeout | undefined;
    addEventListener("resize", () => {
        clearTimeout(resizeTimeout);
        // resize with backoff
        resizeTimeout = setTimeout(() => {
            resizeScreen(window.innerWidth, window.innerHeight);
        }, 500);
    });
}