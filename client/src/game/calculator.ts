import type * as state from '../types/state';

export function renderCalculator(ctx: CanvasRenderingContext2D, canvas: HTMLCanvasElement, state: state.StateCalculator, x: number, y: number) {
    const width = 200; // Half width of the rectangle
    const height = 350; // Half height of the rectangle
    const radius = 10; // Corner radius

    ctx.fillStyle = '#003242';

    // Rectangle
    function drawRoundedRect(x: number, y: number, width: number, height: number, radius: number) {
        // Set shadow properties
        ctx.shadowColor = 'black';
        ctx.shadowBlur = 10;
        ctx.shadowOffsetX = 5;
        ctx.shadowOffsetY = 5;
        // Draw rounded rectangle
        ctx.beginPath();
        ctx.moveTo(x + radius, y);
        ctx.lineTo(x + width - radius, y);
        ctx.arcTo(x + width, y, x + width, y + radius, radius);
        ctx.lineTo(x + width, y + height - radius);
        ctx.arcTo(x + width, y + height, x + width - radius, y + height, radius);
        ctx.lineTo(x + radius, y + height);
        ctx.arcTo(x, y + height, x, y + height - radius, radius);
        ctx.lineTo(x, y + radius);
        ctx.arcTo(x, y, x + radius, y, radius);
        ctx.closePath();
        ctx.fill();
        // Add white border
        ctx.strokeStyle = 'white';
        ctx.lineWidth = 0.5; // Adjust the border width as needed
        ctx.stroke();
        // Reset shadow properties to avoid affecting other drawings
        ctx.shadowColor = 'transparent';
        ctx.shadowBlur = 0;
        ctx.shadowOffsetX = 0;
        ctx.shadowOffsetY = 0;
    }
    drawRoundedRect(x-width, y-height, width*2, height*2, radius);

    // Header
    ctx.fillStyle = 'white'; // Set text color to black
    ctx.textBaseline = 'top'; // Set text baseline to top
    ctx.fillText('Calculator', x - width + 20, y - height + 20);
}