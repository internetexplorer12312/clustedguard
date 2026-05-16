/** Русские строки интерфейса */
export const ru = {
    appTitle: 'ClusterGuard',
    nav: {
        dashboard: 'Обзор',
        servers: 'Серверы',
        clusters: 'Кластеры',
        alerts: 'Алерты',
        theme: 'Тема',
        themeLight: '☀️ Светлая тема',
        themeDark: '🌙 Тёмная тема',
    },
    common: {
        loading: 'Загрузка…',
        cancel: 'Отмена',
        save: 'Сохранить',
        create: 'Создать',
        edit: 'Изменить',
        delete: 'Удалить',
        check: 'Проверить',
        checkAll: 'Проверить все',
        refresh: 'Обновить',
        back: 'Назад',
        none: '— Не выбран —',
        metrics: 'Метрики',
        actions: 'Действия',
    },
    dashboard: {
        title: 'Обзор',
        totalServers: 'Всего серверов',
        online: 'В сети',
        clusters: 'Кластеров',
        unreadAlerts: 'Непрочитанные алерты',
        serversOverview: 'Сводка по серверам',
        checkAllServers: 'Проверить все серверы',
    },
    servers: {
        title: 'Серверы',
        add: 'Добавить сервер',
        edit: 'Изменить сервер',
        empty: 'Серверов пока нет. Добавьте первый на вкладке «Серверы».',
        deleteConfirm: 'Удалить этот сервер?',
        deleted: 'Сервер удалён',
        created: 'Сервер добавлен',
        updated: 'Сервер обновлён',
        checkDone: 'Проверка завершена',
        checkAllDone: 'Все серверы проверены',
        cols: {
            name: 'Имя',
            address: 'Адрес',
            role: 'Роль',
            status: 'Статус',
            latency: 'Задержка',
            cpu: 'ЦП',
            mem: 'ОЗУ',
            disk: 'Диск',
            cluster: 'Кластер',
        },
        form: {
            name: 'Имя',
            host: 'Хост',
            port: 'Порт',
            role: 'Роль',
            cluster: 'Кластер',
            notes: 'Заметки',
            agentSection: 'Агент ClusterGuard',
            agentHint: 'На сервере должен быть установлен и запущен агент ClusterGuard (порт по умолчанию 9100).',
            agentPort: 'Порт агента',
            agentToken: 'Токен агента',
            cpuThreshold: 'Порог ЦП, %',
            memThreshold: 'Порог ОЗУ, %',
            diskThreshold: 'Порог диска, %',
        },
    },
    clusters: {
        title: 'Кластеры',
        add: 'Добавить кластер',
        edit: 'Изменить кластер',
        empty: 'Кластеров пока нет. Создайте группу для объединения серверов.',
        noDescription: 'Без описания',
        serversCount: 'серверов',
        online: 'в сети',
        offline: 'недоступно',
        checkCluster: 'Проверить кластер',
        deleteConfirm: 'Удалить этот кластер?',
        deleted: 'Кластер удалён',
        created: 'Кластер создан',
        updated: 'Кластер обновлён',
        checkDone: 'Проверка кластера завершена',
        form: { name: 'Название', description: 'Описание' },
    },
    alerts: {
        title: 'Алерты',
        empty: 'Алертов нет',
        markRead: 'Прочитано',
    },
    detail: {
        title: 'Метрики сервера',
        updated: 'Метрики обновлены',
        cpu: 'ЦП',
        mem: 'ОЗУ',
        disk: 'Диск',
        chartCpu: 'Загрузка ЦП, %',
        chartMem: 'Загрузка ОЗУ, %',
        chartDisk: 'Загрузка диска, %',
    },
    status: {
        online: 'В сети',
        offline: 'Недоступен',
        degraded: 'Замедление',
        unknown: 'Неизвестно',
    },
    role: {
        master: 'Главный',
        worker: 'Рабочий',
        any: 'Любая',
    },
    checkType: {
        tcp: 'TCP',
        http: 'HTTP',
    },
    alertKind: {
        cpu: 'ЦП',
        memory: 'ОЗУ',
        disk: 'Диск',
    },
} as const;

export function statusLabel(status: string): string {
    const map: Record<string, string> = ru.status;
    return map[status] ?? ru.status.unknown;
}

export function roleLabel(role: string): string {
    const map: Record<string, string> = ru.role;
    return map[role] ?? role;
}

export function alertKindLabel(kind: string): string {
    const map: Record<string, string> = ru.alertKind;
    return map[kind] ?? kind;
}
