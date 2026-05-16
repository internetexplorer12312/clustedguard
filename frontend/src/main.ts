import './style.css';
import './app.css';
import { initTheme, toggleTheme, getTheme } from './theme';
import { drawLineChart, ChartPoint } from './charts';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { ru, statusLabel, roleLabel, alertKindLabel } from './i18n';
import {
    ListServers, CreateServer, UpdateServer, DeleteServer,
    CheckServer, CheckAllServers, ListClusters, CreateCluster,
    UpdateCluster, DeleteCluster, CheckCluster, GetDashboardStats,
    GetMetricsHistory, CollectServerMetrics, ListAlerts, MarkAlertRead, DeleteAlert,
} from '../wailsjs/go/main/App';
import type { main } from '../wailsjs/go/models';

type ServerDTO = main.ServerDTO;
type ClusterSummaryDTO = main.ClusterSummaryDTO;
type DashboardStatsDTO = main.DashboardStatsDTO;
type AlertDTO = main.AlertDTO;
type ServerInputDTO = main.ServerInputDTO;
type ClusterInputDTO = main.ClusterInputDTO;
type View = 'dashboard' | 'servers' | 'clusters' | 'alerts' | 'detail';

let currentView: View = 'dashboard';
let selectedServerId: string | null = null;
let servers: ServerDTO[] = [];
let clusters: ClusterSummaryDTO[] = [];
let alerts: AlertDTO[] = [];
const root = document.querySelector('#app')!;

function escapeHtml(s: string): string {
    const el = document.createElement('div');
    el.textContent = s;
    return el.innerHTML;
}

function metricBar(pct: number, cls: string, thr: number): string {
    const level = thr > 0 && pct >= thr ? ' crit' : pct >= 80 ? ' warn' : '';
    return `<div class="metric-bar"><div class="metric-bar-fill ${cls}${level}" style="width:${Math.min(pct || 0, 100)}%"></div></div>`;
}

function statusBadge(status: string): string {
    const cls = ['online', 'offline', 'degraded', 'unknown'].includes(status) ? status : 'unknown';
    return `<span class="badge ${cls}">${statusLabel(status)}</span>`;
}

function clusterName(id: string): string {
    const c = clusters.find(x => x.id === id);
    return c ? c.name : '—';
}

function showToast(msg: string, isError = false) {
    const t = document.createElement('div');
    t.className = `toast${isError ? ' error' : ''}`;
    t.textContent = msg;
    document.body.appendChild(t);
    setTimeout(() => t.remove(), 3500);
}

function updateThemeButton() {
    const btn = document.getElementById('theme-toggle');
    if (btn) btn.textContent = getTheme() === 'dark' ? ru.nav.themeLight : ru.nav.themeDark;
}

function updateAlertBadge(count: number) {
    const b = document.getElementById('alert-badge');
    if (!b) return;
    b.style.display = count > 0 ? 'inline' : 'none';
    if (count > 0) b.textContent = String(count);
}

function renderServerTable(list: ServerDTO[], compact = false): string {
    if (!list.length) return `<p class="empty">${ru.servers.empty}</p>`;
    const actionHeader = compact ? '' : `<th>${ru.common.actions}</th>`;
    const metricsHeader = compact ? '' : `<th>${ru.servers.cols.cpu}</th><th>${ru.servers.cols.mem}</th><th>${ru.servers.cols.disk}</th>`;
    const rows = list.map(s => {
        const actions = compact ? '' : `<td class="actions-cell">
            <button class="btn btn-sm" data-action="metrics" data-id="${s.id}">${ru.common.metrics}</button>
            <button class="btn btn-sm" data-action="check" data-id="${s.id}">${ru.common.check}</button>
            <button class="btn btn-sm" data-action="edit" data-id="${s.id}">${ru.common.edit}</button>
            <button class="btn btn-sm btn-danger" data-action="delete" data-id="${s.id}">${ru.common.delete}</button></td>`;
        const m = `<td>${(s.cpuPercent||0).toFixed(0)}%</td><td>${(s.memPercent||0).toFixed(0)}%</td><td>${(s.diskPercent||0).toFixed(0)}%</td>`;
        return `<tr><td><strong>${escapeHtml(s.name)}</strong></td><td>${escapeHtml(s.host)}:${s.port}</td>
            <td>${escapeHtml(roleLabel(s.role))}</td><td>${statusBadge(s.status)}</td>${m}
            <td>${s.latencyMs ? s.latencyMs + ' мс' : '—'}</td><td>${escapeHtml(clusterName(s.clusterId))}</td>${actions}</tr>`;
    }).join('');
    return `<div class="table-wrap"><table><thead><tr>
        <th>${ru.servers.cols.name}</th><th>${ru.servers.cols.address}</th><th>${ru.servers.cols.role}</th><th>${ru.servers.cols.status}</th>
        ${metricsHeader}<th>${ru.servers.cols.latency}</th><th>${ru.servers.cols.cluster}</th>${actionHeader}
        </tr></thead><tbody>${rows}</tbody></table></div>`;
}

function bindServerActions(container: HTMLElement) {
    container.querySelectorAll('[data-action]').forEach(btn => {
        btn.addEventListener('click', async () => {
            const id = (btn as HTMLElement).dataset.id!;
            const action = (btn as HTMLElement).dataset.action!;
            if (action === 'metrics') { selectedServerId = id; currentView = 'detail'; render(); await loadData(); }
            else if (action === 'check') { try { await CheckServer(id); await loadData(); showToast(ru.servers.checkDone); } catch (e) { showToast(String(e), true); } }
            else if (action === 'edit') { const s = servers.find(x => x.id === id); if (s) openServerModal(s); }
            else if (action === 'delete') {
                if (!confirm(ru.servers.deleteConfirm)) return;
                try { await DeleteServer(id); await loadData(); showToast(ru.servers.deleted); } catch (e) { showToast(String(e), true); }
            }
        });
    });
}

function renderAlerts() {
    const el = document.getElementById('view-content');
    if (!el) return;
    if (!alerts.length) { el.innerHTML = `<p class="empty">${ru.alerts.empty}</p>`; return; }
    el.innerHTML = `<div class="alerts-list">${alerts.map(a => `
        <div class="alert-item ${a.read ? 'read' : ''}">
            <div><div class="msg">${escapeHtml(a.message)}</div><div class="meta">${escapeHtml(a.serverName)} · ${alertKindLabel(a.kind)} · ${new Date(a.createdAt*1000).toLocaleString('ru-RU')}</div></div>
            <div><button class="btn btn-sm" data-mark="${a.id}">${ru.alerts.markRead}</button>
            <button class="btn btn-sm btn-danger" data-del-alert="${a.id}">×</button></div>
        </div>`).join('')}</div>`;
    el.querySelectorAll('[data-mark]').forEach(btn => btn.addEventListener('click', async () => { await MarkAlertRead((btn as HTMLElement).dataset.mark!); await loadData(); }));
    el.querySelectorAll('[data-del-alert]').forEach(btn => btn.addEventListener('click', async () => { await DeleteAlert((btn as HTMLElement).dataset.delAlert!); await loadData(); }));
}

async function renderServerDetail(id: string) {
    const el = document.getElementById('view-content');
    const s = servers.find(x => x.id === id);
    if (!el || !s) return;
    const hist = await GetMetricsHistory(id, 60);
    const cpuPts: ChartPoint[] = hist.map(h => ({ timestamp: h.timestamp, value: h.cpuPercent }));
    const memPts: ChartPoint[] = hist.map(h => ({ timestamp: h.timestamp, value: h.memPercent }));
    const diskPts: ChartPoint[] = hist.map(h => ({ timestamp: h.timestamp, value: h.diskPercent }));
    el.innerHTML = `<div class="detail-header"><h2>${escapeHtml(s.name)}</h2></div>
        <div class="detail-metrics">
            <div class="metric-tile"><div class="label">${ru.detail.cpu}</div><div class="value">${(s.cpuPercent||0).toFixed(1)}%</div>${metricBar(s.cpuPercent||0,'cpu',s.cpuThreshold)}</div>
            <div class="metric-tile"><div class="label">${ru.detail.mem}</div><div class="value">${(s.memPercent||0).toFixed(1)}%</div>${metricBar(s.memPercent||0,'mem',s.memThreshold)}</div>
            <div class="metric-tile"><div class="label">${ru.detail.disk}</div><div class="value">${(s.diskPercent||0).toFixed(1)}%</div>${metricBar(s.diskPercent||0,'disk',s.diskThreshold)}</div>
        </div>
        <div class="charts-grid">
            <div class="chart-card"><canvas id="chart-cpu"></canvas></div>
            <div class="chart-card"><canvas id="chart-mem"></canvas></div>
            <div class="chart-card"><canvas id="chart-disk"></canvas></div>
        </div>`;
    const st = getComputedStyle(document.documentElement);
    drawLineChart(document.getElementById('chart-cpu') as HTMLCanvasElement, cpuPts, st.getPropertyValue('--chart-cpu'), ru.detail.chartCpu);
    drawLineChart(document.getElementById('chart-mem') as HTMLCanvasElement, memPts, st.getPropertyValue('--chart-mem'), ru.detail.chartMem);
    drawLineChart(document.getElementById('chart-disk') as HTMLCanvasElement, diskPts, st.getPropertyValue('--chart-disk'), ru.detail.chartDisk);
}

function render() {
    root.innerHTML = `<aside class="sidebar"><div class="logo">Cluster<span>Guard</span></div><nav>
        <button class="nav-item ${currentView==='dashboard'?'active':''}" data-view="dashboard">${ru.nav.dashboard}</button>
        <button class="nav-item ${currentView==='servers'?'active':''}" data-view="servers">${ru.nav.servers}</button>
        <button class="nav-item ${currentView==='clusters'?'active':''}" data-view="clusters">${ru.nav.clusters}</button>
        <button class="nav-item ${currentView==='alerts'?'active':''}" data-view="alerts">${ru.nav.alerts}<span class="nav-badge" id="alert-badge" style="display:none">0</span></button>
        </nav><div class="sidebar-footer"><button class="btn theme-btn" id="theme-toggle" type="button"></button></div></aside><main class="main" id="main-content"></main>`;
    document.querySelectorAll('.nav-item').forEach(btn => btn.addEventListener('click', () => {
        currentView = (btn as HTMLElement).dataset.view as View;
        if (currentView !== 'detail') selectedServerId = null;
        render(); loadData();
    }));
    document.getElementById('theme-toggle')?.addEventListener('click', () => {
        toggleTheme(); updateThemeButton();
        if (selectedServerId && currentView === 'detail') renderServerDetail(selectedServerId);
    });
    updateThemeButton(); renderView(); loadData();
}

function renderView() {
    const main = document.getElementById('main-content')!;
    if (currentView === 'dashboard') {
        main.innerHTML = `<header class="header"><h1>${ru.dashboard.title}</h1><div class="header-actions"><button class="btn btn-primary" id="check-all">${ru.dashboard.checkAllServers}</button></div></header><div class="content" id="view-content"><p class="empty">${ru.common.loading}</p></div>`;
        document.getElementById('check-all')?.addEventListener('click', handleCheckAll);
    } else if (currentView === 'servers') {
        main.innerHTML = `<header class="header"><h1>${ru.servers.title}</h1><div class="header-actions"><button class="btn" id="check-all-srv">${ru.common.checkAll}</button><button class="btn btn-primary" id="add-server">${ru.servers.add}</button></div></header><div class="content" id="view-content"></div>`;
        document.getElementById('add-server')?.addEventListener('click', () => openServerModal());
        document.getElementById('check-all-srv')?.addEventListener('click', handleCheckAll);
    } else if (currentView === 'alerts') {
        main.innerHTML = `<header class="header"><h1>${ru.alerts.title}</h1></header><div class="content" id="view-content"></div>`;
    } else if (currentView === 'detail') {
        main.innerHTML = `<header class="header"><h1>${ru.detail.title}</h1><div class="header-actions"><button class="btn" id="back-servers">${ru.common.back}</button><button class="btn btn-primary" id="refresh-metrics">${ru.common.refresh}</button></div></header><div class="content" id="view-content"></div>`;
        document.getElementById('back-servers')?.addEventListener('click', () => { currentView='servers'; selectedServerId=null; render(); loadData(); });
        document.getElementById('refresh-metrics')?.addEventListener('click', () => { if (selectedServerId) refreshDetail(selectedServerId); });
    } else {
        main.innerHTML = `<header class="header"><h1>${ru.clusters.title}</h1><div class="header-actions"><button class="btn btn-primary" id="add-cluster">${ru.clusters.add}</button></div></header><div class="content" id="view-content"></div>`;
        document.getElementById('add-cluster')?.addEventListener('click', () => openClusterModal());
    }
}

async function loadData() {
    try {
        [servers, clusters] = await Promise.all([ListServers(), ListClusters()]);
        const stats = await GetDashboardStats();
        alerts = await ListAlerts(50);
        updateAlertBadge(stats.unreadAlerts);
        if (currentView === 'dashboard') renderDashboard(stats);
        else if (currentView === 'servers') renderServers();
        else if (currentView === 'alerts') renderAlerts();
        else if (currentView === 'detail' && selectedServerId) await renderServerDetail(selectedServerId);
        else renderClustersView();
    } catch (e) { showToast(String(e), true); }
}

function renderDashboard(stats: DashboardStatsDTO) {
    const el = document.getElementById('view-content'); if (!el) return;
    el.innerHTML = `<div class="stats-grid">
        <div class="stat-card"><div class="label">${ru.dashboard.totalServers}</div><div class="value">${stats.totalServers}</div></div>
        <div class="stat-card online"><div class="label">${ru.dashboard.online}</div><div class="value">${stats.onlineServers}</div></div>
        <div class="stat-card clusters"><div class="label">${ru.dashboard.clusters}</div><div class="value">${stats.totalClusters}</div></div>
        <div class="stat-card"><div class="label">${ru.dashboard.unreadAlerts}</div><div class="value">${stats.unreadAlerts||0}</div></div>
        </div><h2 style="margin-bottom:12px;font-size:16px;">${ru.dashboard.serversOverview}</h2>${renderServerTable(servers,true)}`;
}

function renderServers() { const el = document.getElementById('view-content'); if (!el) return; el.innerHTML = renderServerTable(servers); bindServerActions(el); }

function renderClustersView() {
    const el = document.getElementById('view-content'); if (!el) return;
    if (!clusters.length) { el.innerHTML = `<p class="empty">${ru.clusters.empty}</p>`; return; }
    el.innerHTML = `<div class="cluster-cards">${clusters.map(c => `<article class="cluster-card"><h3>${escapeHtml(c.name)}</h3><p>${escapeHtml(c.description||ru.clusters.noDescription)}</p>
        <div class="cluster-meta"><span>${c.totalServers} ${ru.clusters.serversCount}</span><span class="online">${c.onlineCount} ${ru.clusters.online}</span><span class="offline">${c.offlineCount} ${ru.clusters.offline}</span></div>
        <div class="cluster-card-actions"><button class="btn btn-sm btn-primary" data-cluster-check="${c.id}">${ru.clusters.checkCluster}</button>
        <button class="btn btn-sm" data-cluster-edit="${c.id}">${ru.common.edit}</button><button class="btn btn-sm btn-danger" data-cluster-delete="${c.id}">${ru.common.delete}</button></div></article>`).join('')}</div>`;
    el.querySelectorAll('[data-cluster-check]').forEach(btn => btn.addEventListener('click', async () => {
        try { await CheckCluster((btn as HTMLElement).dataset.clusterCheck!); await loadData(); showToast(ru.clusters.checkDone); } catch(e){showToast(String(e),true);}
    }));
    el.querySelectorAll('[data-cluster-edit]').forEach(btn => btn.addEventListener('click', () => {
        const c = clusters.find(x => x.id === (btn as HTMLElement).dataset.clusterEdit); if (c) openClusterModal(c);
    }));
    el.querySelectorAll('[data-cluster-delete]').forEach(btn => btn.addEventListener('click', async () => {
        if (!confirm(ru.clusters.deleteConfirm)) return;
        try { await DeleteCluster((btn as HTMLElement).dataset.clusterDelete!); await loadData(); showToast(ru.clusters.deleted); } catch(e){showToast(String(e),true);}
    }));
}

async function handleCheckAll() {
    const btn = document.getElementById('check-all') || document.getElementById('check-all-srv');
    btn?.classList.add('loading');
    try { await CheckAllServers(); await loadData(); showToast(ru.servers.checkAllDone); } catch(e){showToast(String(e),true);} finally { btn?.classList.remove('loading'); }
}

async function refreshDetail(id: string) {
    try { await CollectServerMetrics(id); await loadData(); showToast(ru.detail.updated); } catch(e){showToast(String(e),true);}
}

function openServerModal(server?: ServerDTO) {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    const isEdit = !!server;
    const clusterOpts = clusters.map(c =>
        `<option value="${c.id}" ${server?.clusterId === c.id ? 'selected' : ''}>${escapeHtml(c.name)}</option>`,
    ).join('');

    overlay.innerHTML = `
        <div class="modal"><h2>${isEdit ? ru.servers.edit : ru.servers.add}</h2>
        <form id="server-form">
            <div class="form-group"><label>${ru.servers.form.name}</label>
                <input name="name" required value="${server?.name || ''}" /></div>
            <div class="form-row">
                <div class="form-group"><label>${ru.servers.form.host}</label>
                    <input name="host" required value="${server?.host || ''}" placeholder="192.168.1.10" /></div>
                <div class="form-group"><label>${ru.servers.form.role}</label>
                    <select name="role">
                        <option value="master" ${server?.role === 'master' ? 'selected' : ''}>${ru.role.master}</option>
                        <option value="worker" ${!server || server.role === 'worker' ? 'selected' : ''}>${ru.role.worker}</option>
                        <option value="any" ${server?.role === 'any' ? 'selected' : ''}>${ru.role.any}</option>
                    </select></div>
            </div>
            <div class="form-group"><label>${ru.servers.form.cluster}</label>
                <select name="clusterId"><option value="">${ru.common.none}</option>${clusterOpts}</select></div>

            <section class="form-section">
                <div class="form-section-title">${ru.servers.form.agentSection}</div>
                <p class="form-section-hint">${ru.servers.form.agentHint}</p>
                <div class="form-row">
                    <div class="form-group"><label>${ru.servers.form.agentPort}</label>
                        <input name="agentPort" type="number" min="1" max="65535" required value="${server?.agentPort || 9100}" /></div>
                    <div class="form-group"><label>${ru.servers.form.agentToken}</label>
                        <input name="agentToken" value="${server?.agentToken || ''}" placeholder="CLUSTERGUARD_TOKEN" autocomplete="off" /></div>
                </div>
            </section>

            <div class="form-row">
                <div class="form-group"><label>${ru.servers.form.cpuThreshold}</label>
                    <input name="cpuThreshold" type="number" min="1" max="100" value="${server?.cpuThreshold || 90}" /></div>
                <div class="form-group"><label>${ru.servers.form.memThreshold}</label>
                    <input name="memThreshold" type="number" min="1" max="100" value="${server?.memThreshold || 90}" /></div>
                <div class="form-group"><label>${ru.servers.form.diskThreshold}</label>
                    <input name="diskThreshold" type="number" min="1" max="100" value="${server?.diskThreshold || 90}" /></div>
            </div>
            <div class="form-group"><label>${ru.servers.form.notes}</label>
                <textarea name="notes">${server?.notes || ''}</textarea></div>
            <div class="modal-actions">
                <button type="button" class="btn" id="cancel-server">${ru.common.cancel}</button>
                <button type="submit" class="btn btn-primary">${isEdit ? ru.common.save : ru.common.create}</button>
            </div>
        </form></div>`;


    document.body.appendChild(overlay);
    overlay.querySelector('#cancel-server')?.addEventListener('click', () => overlay.remove());
    overlay.addEventListener('click', e => { if (e.target === overlay) overlay.remove(); });
    overlay.querySelector('#server-form')?.addEventListener('submit', async e => {
        e.preventDefault();
        const fd = new FormData(e.target as HTMLFormElement);
        const agentPort = parseInt(fd.get('agentPort') as string, 10) || 9100;
        const input: ServerInputDTO = {
            id: server?.id || '',
            name: fd.get('name') as string,
            host: fd.get('host') as string,
            port: agentPort,
            role: fd.get('role') as string,
            tags: server?.tags || [],
            checkType: 'agent',
            checkPath: '/',
            clusterId: fd.get('clusterId') as string,
            notes: fd.get('notes') as string,
            useAgent: true,
            agentPort,
            agentToken: (fd.get('agentToken') as string) || '',
            cpuThreshold: parseFloat(fd.get('cpuThreshold') as string) || 90,
            memThreshold: parseFloat(fd.get('memThreshold') as string) || 90,
            diskThreshold: parseFloat(fd.get('diskThreshold') as string) || 90,
        };
        try {
            if (isEdit) await UpdateServer(input);
            else await CreateServer(input);
            overlay.remove();
            await loadData();
            showToast(isEdit ? ru.servers.updated : ru.servers.created);
        } catch (err) {
            showToast(String(err), true);
        }
    });
}

function openClusterModal(cluster?: ClusterSummaryDTO) {
    const overlay = document.createElement('div'); overlay.className = 'modal-overlay'; const isEdit = !!cluster;
    overlay.innerHTML = `<div class="modal"><h2>${isEdit?ru.clusters.edit:ru.clusters.add}</h2><form id="cluster-form">
        <div class="form-group"><label>${ru.clusters.form.name}</label><input name="name" required value="${cluster?.name||''}"/></div>
        <div class="form-group"><label>${ru.clusters.form.description}</label><textarea name="description">${cluster?.description||''}</textarea></div>
        <div class="modal-actions"><button type="button" class="btn" id="cancel-cluster">${ru.common.cancel}</button>
        <button type="submit" class="btn btn-primary">${isEdit?ru.common.save:ru.common.create}</button></div></form></div>`;
    document.body.appendChild(overlay);
    overlay.querySelector('#cancel-cluster')?.addEventListener('click', () => overlay.remove());
    overlay.addEventListener('click', e => { if (e.target===overlay) overlay.remove(); });
    overlay.querySelector('#cluster-form')?.addEventListener('submit', async e => {
        e.preventDefault(); const fd = new FormData(e.target as HTMLFormElement);
        const input: ClusterInputDTO = { id: cluster?.id||'', name: fd.get('name') as string, description: fd.get('description') as string, serverIds: cluster?.serverIds||[] };
        try { if (isEdit) await UpdateCluster(input); else await CreateCluster(input); overlay.remove(); await loadData();
            showToast(isEdit?ru.clusters.updated:ru.clusters.created); } catch(err){showToast(String(err),true);}
    });
}

initTheme();
updateThemeButton();
EventsOn('alert', (a: AlertDTO) => { showToast('⚠ ' + a.message, true); loadData(); });
render();
