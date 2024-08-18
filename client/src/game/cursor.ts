import type * as state from '../types/state';

export function renderCursor(ctx: CanvasRenderingContext2D, canvas: HTMLCanvasElement, state: state.Member) {
    // Draw the new cursor
    ctx.beginPath();
    ctx.arc(state.X*(window.devicePixelRatio || 1), state.Y*(window.devicePixelRatio || 1), 5, 0, Math.PI * 2); // Draw a circle as the cursor
    ctx.fillStyle = 'white'; // Set the cursor color
    ctx.fill();
    ctx.closePath();
}