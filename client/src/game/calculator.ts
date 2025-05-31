import type { State } from '@lib/schema/state';
import type { ButtonProperties, ButtonLayout, ButtonStyle, ButtonState } from '../types/types';

export default function RenderCalculator(ctx: CanvasRenderingContext2D, canvas: HTMLCanvasElement, state: State) {
    // Clear canvas
    const gradient = ctx.createLinearGradient(0, 0, 0, canvas.height);
    gradient.addColorStop(0, '#1a1a1a');
    gradient.addColorStop(1, '#0d1117');
    ctx.fillStyle = gradient;
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Dimensions
    const container = canvas.parentElement;
    const containerWidth = container.clientWidth;
    const containerHeight = container.clientHeight;

    // Aspect ratio 3:4
    const aspectRatio = 3 / 4;
    let width = Math.min(containerWidth * 0.9, 400);
    let height = width / aspectRatio;
    if (height > containerHeight * 0.9) {
      height = containerHeight * 0.9;
      width = height * aspectRatio;
    }

    // Calculate center
    const offsetX = (containerWidth - width) / 2;
    const offsetY = (containerHeight - height) / 2;

    // Button layout
    const displayHeight = height / 6;
    const buttonWidth = width / 4; // 4 columns
    const buttonHeight = (height - displayHeight) / 5; // 5 rows
    global.buttonMap = {
        // row 1
        clear: newButtonState(newButtonLayout(buttonPropertiesClear, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleLight, Boolean(state.Calculator?.clear)),
        bracket: newButtonState(newButtonLayout(buttonPropertiesBracket, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleLight, Boolean(state.Calculator?.bracket)),
        percentage: newButtonState(newButtonLayout(buttonPropertiesPercentage, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleLight, Boolean(state.Calculator?.percentage)),
        divide: newButtonState(newButtonLayout(buttonPropertiesDivide, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleLight, Boolean(state.Calculator?.divide)),
        // row 2
        seven: newButtonState(newButtonLayout(buttonPropertiesSeven, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.seven)),
        eight:  newButtonState(newButtonLayout(buttonPropertiesEight, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.eight)),
        nine: newButtonState(newButtonLayout(buttonPropertiesNine, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.nine)),
        times: newButtonState(newButtonLayout(buttonPropertiesTimes, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleLight, Boolean(state.Calculator?.times)),
        // row 3
        four: newButtonState(newButtonLayout(buttonPropertiesFour, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.four)),
        five: newButtonState(newButtonLayout(buttonPropertiesFive, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.five)),
        six: newButtonState(newButtonLayout(buttonPropertiesSix, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.six)),
        minus: newButtonState(newButtonLayout(buttonPropertiesMinus, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleLight, Boolean(state.Calculator?.minus)),
        // row 4
        one: newButtonState(newButtonLayout(buttonPropertiesOne, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.one)),
        two: newButtonState(newButtonLayout(buttonPropertiesTwo, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.two)),
        three: newButtonState(newButtonLayout(buttonPropertiesThree, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.three)),
        plus: newButtonState(newButtonLayout(buttonPropertiesPlus, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleLight, Boolean(state.Calculator?.plus)),
        // row 5
        negate: newButtonState(newButtonLayout(buttonPropertiesNegate, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.negate)),
        zero: newButtonState(newButtonLayout(buttonPropertiesZero, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.zero)),
        decimal: newButtonState(newButtonLayout(buttonPropertiesDecimal, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleDark, Boolean(state.Calculator?.decimal)),
        equals: newButtonState(newButtonLayout(buttonPropertiesEquals, displayHeight, offsetX, offsetY, buttonWidth, buttonHeight), buttonStyleRed, Boolean(state.Calculator?.equals)),
    };
    const buttonList = Object.values(global.buttonMap);

    // Calculate limits
    const minX = Math.min(...buttonList.map(b => b.x));
    const maxX = Math.max(...buttonList.map(b => b.x + b.width));
    const minY = global.buttonMap.clear.y - buttonHeight;
    const maxY = Math.max(...buttonList.map(b => b.y + b.height));

    // Draw body
    const calcGradient = ctx.createLinearGradient(minX - 15, minY - 15, minX - 15, maxY + 15);
    calcGradient.addColorStop(0, '#2d3748');
    calcGradient.addColorStop(1, '#1a202c');
    ctx.fillStyle = calcGradient;
    
    // Draw body shadow
    ctx.shadowColor = 'rgba(0, 0, 0, 0.3)';
    ctx.shadowBlur = 20;
    ctx.shadowOffsetX = 0;
    ctx.shadowOffsetY = 10;
    fillRoundedRect(ctx, minX - 15, minY - 15, maxX - minX + 30, maxY - minY + 30, 12);
    ctx.shadowColor = 'transparent';
    ctx.shadowBlur = 0;
    ctx.shadowOffsetX = 0;
    ctx.shadowOffsetY = 0;

    // Draw display
    const displayGradient = ctx.createLinearGradient(minX, minY, minX, minY + displayHeight);
    displayGradient.addColorStop(0, '#0f172a');
    displayGradient.addColorStop(1, '#1e293b');
    ctx.fillStyle = displayGradient;
    fillRoundedRect(ctx, minX + 5, minY + 5, maxX - minX - 10, displayHeight - 10, 8);

    // Draw display border
    ctx.strokeStyle = '#374151';
    ctx.lineWidth = 1;
    strokeRoundedRect(ctx, minX + 5, minY + 5, maxX - minX - 10, displayHeight - 10, 8);

    // Draw equation text
    ctx.fillStyle = '#94a3b8';
    ctx.font = `${Math.min(displayHeight * 0.25, 18)}px -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`;
    ctx.textAlign = 'right';
    ctx.textBaseline = 'top';
    ctx.fillText(state.Calculator?.Equation ?? '', maxX - 25, minY + 15);

    // Draw result text
    ctx.fillStyle = '#f8fafc';
    ctx.font = `${Math.min(displayHeight * 0.35, 28)}px -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`;
    ctx.textAlign = 'right';
    ctx.textBaseline = 'bottom';
    ctx.fillText(state.Calculator?.Display ?? 'Disconnected', maxX - 25, minY + displayHeight - 15);

    // Draw buttons
    buttonList.forEach((button) => {
      const buttonX = button.x + buttonPadding;
      const buttonY = button.y + buttonPadding;
      const buttonWidth = button.width - buttonPadding * 2;
      const buttonHeight = button.height - buttonPadding * 2;

      // Create button gradient
      const buttonGradient = ctx.createLinearGradient(buttonX, buttonY, buttonX, buttonY + buttonHeight);
      buttonGradient.addColorStop(0, button.colourPrimary);
      buttonGradient.addColorStop(1, button.colourSecondary);

      // Draw button shadow if depressed
      if (!button.pressed) {
        ctx.fillStyle = 'rgba(0, 0, 0, 0.2)';
        fillRoundedRect(ctx, buttonX + 2, buttonY + 2, buttonWidth, buttonHeight, 6);
      }

      // Draw button background
      ctx.fillStyle = buttonGradient;
      const offsetY = button.pressed ? 2 : 0;
      fillRoundedRect(ctx, buttonX, buttonY + offsetY, buttonWidth, buttonHeight - offsetY, 6);

      // Draw button border
      ctx.strokeStyle = button.borderColor;
      ctx.lineWidth = 1;
      strokeRoundedRect(ctx, buttonX, buttonY + offsetY, buttonWidth, buttonHeight - offsetY, 6);

      // Add inner highlight if depressed
      if (!button.pressed) {
        ctx.strokeStyle = 'rgba(255, 255, 255, 0.1)';
        ctx.lineWidth = 1;
        strokeRoundedRect(ctx, buttonX + 1, buttonY + 1, buttonWidth - 2, buttonHeight - 2, 5);
      }

      // Draw button text
      ctx.fillStyle = button.textColor;
      ctx.font = `${Math.min(buttonHeight * 0.35, 26)}px -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`;
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText(
        button.text,
        buttonX + buttonWidth / 2,
        buttonY + offsetY + (buttonHeight - offsetY) / 2
      );
    });
}

const drawRoundedRect = (ctx: CanvasRenderingContext2D, x: number, y: number, width: number, height: number, radius: number) => {
    ctx.beginPath();
    ctx.moveTo(x + radius, y);
    ctx.lineTo(x + width - radius, y);
    ctx.quadraticCurveTo(x + width, y, x + width, y + radius);
    ctx.lineTo(x + width, y + height - radius);
    ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height);
    ctx.lineTo(x + radius, y + height);
    ctx.quadraticCurveTo(x, y + height, x, y + height - radius);
    ctx.lineTo(x, y + radius);
    ctx.quadraticCurveTo(x, y, x + radius, y);
    ctx.closePath();
};
const fillRoundedRect = (ctx: CanvasRenderingContext2D, x: number, y: number, width: number, height: number, radius: number) => {
    drawRoundedRect(ctx, x, y, width, height, radius);
    ctx.fill();
};
const strokeRoundedRect = (ctx: CanvasRenderingContext2D, x: number, y: number, width: number, height: number, radius: number) => {
    drawRoundedRect(ctx, x, y, width, height, radius);
    ctx.stroke();
};

const buttonPadding = 3;

// Row 1
const buttonPropertiesClear: ButtonProperties = { text: 'C', col: 0, row: 0 };
const buttonPropertiesBracket: ButtonProperties = { text: '( )', col: 1, row: 0 };
const buttonPropertiesPercentage: ButtonProperties = { text: '%', col: 2, row: 0 };
const buttonPropertiesDivide: ButtonProperties = { text: '÷', col: 3, row: 0 };
// Row 2
const buttonPropertiesSeven: ButtonProperties = { text: '7', col: 0, row: 1 };
const buttonPropertiesEight: ButtonProperties = { text: '8', col: 1, row: 1 };
const buttonPropertiesNine: ButtonProperties = { text: '9', col: 2, row: 1 };
const buttonPropertiesTimes: ButtonProperties = { text: '×', col: 3, row: 1 };
// Row 3
const buttonPropertiesFour: ButtonProperties = { text: '4', col: 0, row: 2 };
const buttonPropertiesFive: ButtonProperties = { text: '5', col: 1, row: 2 };
const buttonPropertiesSix: ButtonProperties = { text: '6', col: 2, row: 2 };
const buttonPropertiesMinus: ButtonProperties = { text: '−', col: 3, row: 2 };
// Row 4
const buttonPropertiesOne: ButtonProperties = { text: '1', col: 0, row: 3 };
const buttonPropertiesTwo: ButtonProperties = { text: '2', col: 1, row: 3 };
const buttonPropertiesThree: ButtonProperties = { text: '3', col: 2, row: 3 };
const buttonPropertiesPlus: ButtonProperties = { text: '+', col: 3, row: 3 };
// Row 5
const buttonPropertiesNegate: ButtonProperties = { text: '⁺∕₋', col: 0, row: 4 };
const buttonPropertiesZero: ButtonProperties = { text: '0', col: 1, row: 4 };
const buttonPropertiesDecimal: ButtonProperties = { text: '.', col: 2, row: 4 };
const buttonPropertiesEquals: ButtonProperties = { text: '=', col: 3, row: 4 };

const newButtonLayout = (properties: ButtonProperties, displayHeight: number, offsetX: number, offsetY: number, buttonWidth: number, buttonHeight: number): ButtonLayout => ({
    ...properties,
    x: offsetX + properties.col * buttonWidth,
    y: offsetY + displayHeight + properties.row * buttonHeight,
    width: buttonWidth,
    height: buttonHeight,
});
const newButtonState = (layout: ButtonLayout, style: ButtonStyle, pressed: boolean): ButtonState => ({
    ...layout,
    pressed: pressed,
    colourPrimary: pressed ? style.pressed.colourPrimary : style.depressed.colourPrimary,
    colourSecondary: pressed ? style.pressed.colourSecondary : style.depressed.colourSecondary,
    textColor: style.textColor,
    borderColor: style.borderColor,
});

const buttonStyleRed: ButtonStyle = {
    pressed: {
        colourPrimary: '#dc2626',
        colourSecondary: '#b91c1c',
    },
    depressed: {
        colourPrimary: '#e14848',
        colourSecondary: '#dc2626',
    },
    textColor: '#ffffff',
    borderColor: '#991b1b',
};
const buttonStyleLight: ButtonStyle = {
    pressed: {
        colourPrimary: '#6b7280',
        colourSecondary: '#4b5563',
    },
    depressed: {
        colourPrimary: '#9ca3af',
        colourSecondary: '#6b7280',
    },
    textColor: '#111827',
    borderColor: '#374151',
};
const buttonStyleDark: ButtonStyle = {
    pressed: {
        colourPrimary: '#374151',
        colourSecondary: '#1f2937',
    },
    depressed: {
        colourPrimary: '#4b5563',
        colourSecondary: '#374151',
    },
    textColor: '#f9fafb',
    borderColor: '#6b7280',
};
