import type { Member } from '@lib/schema/state';

export default function RenderCursor(ctx: CanvasRenderingContext2D, canvas: HTMLCanvasElement, state: Member) {
    // Generate colour (golden angle distribution)
    const hue = (state.Id * 137.508) % 360; 
    const color = `hsl(${hue}, 100%, 50%)`;

    // Draw cursor
    ctx.beginPath();
    ctx.arc(state.X, state.Y, 5, 0, Math.PI * 2);
    ctx.fillStyle = color;
    ctx.fill();
    ctx.closePath();

    // Draw name
    ctx.font = `12px sans-serif`;
    ctx.fillStyle = color;
    ctx.textAlign = 'center';
    ctx.textBaseline = 'top';
    ctx.fillText(`Player ${state.Id}`, state.X, state.Y+15);
}
