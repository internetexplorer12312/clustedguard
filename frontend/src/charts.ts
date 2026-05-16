export interface ChartPoint {
    timestamp: number;
    value: number;
}

export function drawLineChart(
    canvas: HTMLCanvasElement,
    points: ChartPoint[],
    color: string,
    label: string,
    maxY = 100,
) {
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const dpr = window.devicePixelRatio || 1;
    const w = canvas.clientWidth;
    const h = canvas.clientHeight;
    canvas.width = w * dpr;
    canvas.height = h * dpr;
    ctx.scale(dpr, dpr);

    const pad = { top: 20, right: 12, bottom: 24, left: 36 };
    const plotW = w - pad.left - pad.right;
    const plotH = h - pad.top - pad.bottom;

    ctx.clearRect(0, 0, w, h);

    const styles = getComputedStyle(document.documentElement);
    const border = styles.getPropertyValue('--border').trim() || '#333';
    const muted = styles.getPropertyValue('--text-muted').trim() || '#888';

    ctx.fillStyle = muted;
    ctx.font = '11px sans-serif';
    ctx.fillText(label, pad.left, 14);

    ctx.strokeStyle = border;
    ctx.lineWidth = 1;
    for (let i = 0; i <= 4; i++) {
        const y = pad.top + (plotH * i) / 4;
        ctx.beginPath();
        ctx.moveTo(pad.left, y);
        ctx.lineTo(pad.left + plotW, y);
        ctx.stroke();
        const val = Math.round(maxY - (maxY * i) / 4);
        ctx.fillStyle = muted;
        ctx.fillText(`${val}%`, 4, y + 4);
    }

    if (points.length < 2) {
        ctx.fillStyle = muted;
        ctx.fillText('Нет данных', pad.left + plotW / 2 - 30, pad.top + plotH / 2);
        return;
    }

    ctx.strokeStyle = color;
    ctx.lineWidth = 2;
    ctx.beginPath();
    points.forEach((p, i) => {
        const x = pad.left + (plotW * i) / (points.length - 1);
        const y = pad.top + plotH - (Math.min(p.value, maxY) / maxY) * plotH;
        if (i === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
    });
    ctx.stroke();

    ctx.lineTo(pad.left + plotW, pad.top + plotH);
    ctx.lineTo(pad.left, pad.top + plotH);
    ctx.closePath();
    ctx.fillStyle = color + '22';
    ctx.fill();
}

export function formatBytes(n: number): string {
    if (n >= 1e12) return (n / 1e12).toFixed(1) + ' ТБ';
    if (n >= 1e9) return (n / 1e9).toFixed(1) + ' ГБ';
    if (n >= 1e6) return (n / 1e6).toFixed(1) + ' МБ';
    if (n >= 1e3) return (n / 1e3).toFixed(1) + ' КБ';
    return n + ' Б';
}
