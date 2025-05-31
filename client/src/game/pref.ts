
export default function RenderPreformance(ctx: CanvasRenderingContext2D, fps: number, gpuLoad: number, cpuLoad: number) {
    const fpsView = Math.round(fps).toString().padStart(3, ' ');
    const gpuLoadPercentage = (100 * gpuLoad).toFixed(2).padStart(5, ' ');
    const cpuLoadPercentage = (100 * cpuLoad).toFixed(2).padStart(5, ' ');

    ctx.font = '18px "Consolas", "Menlo", "Courier New", monospace';

    ctx.fillStyle = 'lightgrey';
    ctx.fillText(`FPS: ${fpsView}   `, 75, 25);

    ctx.fillStyle = 'lightgreen';
    ctx.fillText(`GPU: ${gpuLoadPercentage}%`, 75, 50);

    ctx.fillStyle = "lightblue";
    ctx.fillText(`CPU: ${cpuLoadPercentage}%`, 75, 75);
}
