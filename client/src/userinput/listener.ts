import { PayloadKind } from '../types/enums';
import type { Payload, PayloadInput, ButtonState } from '../types/types';

export async function listenForUserInput(canvas: HTMLCanvasElement, logic: Worker) {
    // Register input events
    const keyState: Record<string, boolean> = {};
    let mouseLeftState: boolean | undefined;
    let xState: number | undefined;
    let yState: number | undefined;
    addEventListener("mousedown", (ev) => {
        if ( !(ev.button === 0 && mouseLeftState !== true) && !(ev.clientX !== xState) && !(ev.clientY !== yState) ) return;
        const rect = canvas.getBoundingClientRect();
        let data: PayloadInput = {};
        if (ev.button === 0 && mouseLeftState !== true) {
            data.mouseleft = true;
        }
        if (ev.clientX !== xState) {
            data.mousex = ev.clientX;
            data.x = ev.clientX - rect.left;
        }
        if (ev.clientY !== yState) {
            data.mousey = ev.clientY;
            data.y = ev.clientY - rect.top;
        }
        data = findButtonAtCoordinates(data);
        const payload: Payload<PayloadInput> = {
            kind: PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    });
    addEventListener("mouseup", (ev) => {
        if ( !(ev.button === 0 && mouseLeftState !== false) && !(ev.clientX !== xState) && !(ev.clientY !== yState) ) return;
        const rect = canvas.getBoundingClientRect();
        const data: PayloadInput = {};
        if (ev.button === 0 && mouseLeftState !== false) {
            data.mouseleft = false;
        }
        if (ev.clientX !== xState) {
            data.mousex = ev.clientX;
            data.x = ev.clientX - rect.left;
        }
        if (ev.clientY !== yState) {
            data.mousey = ev.clientY;
            data.y = ev.clientY - rect.top;
        }
        const payload: Payload<PayloadInput> = {
            kind: PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    });
    addEventListener("mousemove", (ev) => {
        if ( !(ev.clientX !== xState) && !(ev.clientY !== yState) ) return;
        const rect = canvas.getBoundingClientRect();
        const data: PayloadInput = {};
        if (ev.clientX !== xState) {
            data.mousex = ev.clientX;
            data.x = ev.clientX - rect.left;
        }
        if (ev.clientY !== yState) {
            data.mousey = ev.clientY;
            data.y = ev.clientY - rect.top;
        }
        const payload: Payload<PayloadInput> = {
            kind: PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    });
    addEventListener("touchstart", (ev) => { // equivalent to mousedown
        const touch = ev.touches[0];
        if ( !(!mouseLeftState) && !(touch.clientX !== xState) && !(touch.clientY !== yState) ) return;
        const rect = canvas.getBoundingClientRect();
        let data: PayloadInput = {};
        if (!mouseLeftState) {
            data.mouseleft = true;
        }
        if (touch.clientX !== xState) {
            data.mousex = touch.clientX;
            data.x = touch.clientX - rect.left;
        }
        if (touch.clientY !== yState) {
            data.mousey = touch.clientY;
            data.y = touch.clientY - rect.top;
        }
        data = findButtonAtCoordinates(data);
        const payload: Payload<PayloadInput> = {
            kind: PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    }, { passive: false });
    addEventListener("touchend", (ev) => { // equivalent to mouseup
        const touch = ev.changedTouches[0];
        if ( !(mouseLeftState) && !(touch.clientX !== xState) && !(touch.clientY !== yState) ) return;
        const rect = canvas.getBoundingClientRect();
        const data: PayloadInput = {};
        if (mouseLeftState) {
            data.mouseleft = false;
        }
        if (touch.clientX !== xState) {
            data.mousex = touch.clientX;
            data.x = touch.clientX - rect.left;
        }
        if (touch.clientY !== yState) {
            data.mousey = touch.clientY;
            data.y = touch.clientY - rect.top;
        }
        const payload: Payload<PayloadInput> = {
            kind: PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    }, { passive: false });
    addEventListener("touchmove", (ev) => { // equivalent to mousemove
        const touch = ev.touches[0];
        if ( !(touch.clientX !== xState) && !(touch.clientY !== yState) ) return;
        const rect = canvas.getBoundingClientRect();
        const data: PayloadInput = {};
        if (touch.clientX !== xState) {
            data.mousex = touch.clientX;
            data.x = touch.clientX - rect.left;
        }
        if (touch.clientY !== yState) {
            data.mousey = touch.clientY;
            data.y = touch.clientY - rect.top;
        }
        const payload: Payload<PayloadInput> = {
            kind: PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    }, { passive: false });
    const resizeScreen = (width: number, height: number): void => {
        if (canvas.width === width && canvas.height === height) return;
        // Set canvas size to window size
        const dpr = window.devicePixelRatio || 1;
        canvas.width = window.innerWidth * dpr;
        canvas.height = window.innerHeight * dpr;
        canvas.style.width = window.innerWidth + "px";
        canvas.style.height = window.innerHeight + "px";
        canvas.getContext("2d").scale(dpr, dpr);
        const payload: Payload<PayloadInput> = {
            kind: PayloadKind.input,
            payload: {
                width: width,
                height: height,
            },
        };
        logic.postMessage(payload);
    }
    let resizeTimeout: NodeJS.Timeout | undefined = undefined;
    addEventListener("resize", () => {
        clearTimeout(resizeTimeout);
        // resize with backoff
        resizeTimeout = setTimeout(() => {
            resizeScreen(window.innerWidth, window.innerHeight);
        }, 500);
    });
}

function findButtonAtCoordinates(data: PayloadInput): PayloadInput {
    // row 1
    if (isButtonAtCoordinates(data, global.buttonMap.clear)) {
        data.clear = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.bracket)) {
        data.bracket = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.percentage)) {
        data.percentage = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.divide)) {
        data.divide = true;
        return data;
    }
    // row 2
    if (isButtonAtCoordinates(data, global.buttonMap.seven)) {
        data.seven = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.eight)) {
        data.eight = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.nine)) {
        data.nine = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.times)) {
        data.times = true;
        return data;
    }
    // row 3
    if (isButtonAtCoordinates(data, global.buttonMap.four)) {
        data.four = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.five)) {
        data.five = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.six)) {
        data.six = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.minus)) {
        data.minus = true;
        return data;
    }
    // row 4
    if (isButtonAtCoordinates(data, global.buttonMap.one)) {
        data.one = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.two)) {
        data.two = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.three)) {
        data.three = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.plus)) {
        data.plus = true;
        return data;
    }
    // row 5
    if (isButtonAtCoordinates(data, global.buttonMap.negate)) {
        data.negate = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.zero)) {
        data.zero = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.decimal)) {
        data.decimal = true;
        return data;
    }
    if (isButtonAtCoordinates(data, global.buttonMap.equals)) {
        data.equals = true;
        return data;
    }
};

function isButtonAtCoordinates(data: PayloadInput, button: ButtonState): boolean {
    return data.x >= button.x && data.x <= button.x + button.width && data.y >= button.y && data.y <= button.y + button.height;
};
