const STORAGE_KEY = 'clusterguard-theme';

export type Theme = 'dark' | 'light';

export function getTheme(): Theme {
    const saved = localStorage.getItem(STORAGE_KEY) as Theme | null;
    if (saved === 'light' || saved === 'dark') return saved;
    return 'dark';
}

export function applyTheme(theme: Theme) {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem(STORAGE_KEY, theme);
    const btn = document.getElementById('theme-toggle');
    if (btn) btn.textContent = theme === 'dark' ? '☀️ Светлая' : '🌙 Тёмная';
}

export function toggleTheme(): Theme {
    const next = getTheme() === 'dark' ? 'light' : 'dark';
    applyTheme(next);
    return next;
}

export function initTheme() {
    applyTheme(getTheme());
}
