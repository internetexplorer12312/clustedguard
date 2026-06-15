# Презентация ВКР — ClusterGuard

Файл: `Desktopnaya-sistema-monitoringa-klastera-serverov.pptx` (10 слайдов).

## Что обновлено

| Слайд | Содержание | Источник |
|-------|------------|----------|
| 1, 10 | Титул / финал | `diagrams/title.png` |
| 2 | Актуальность | `diagrams/context.png` — схема контекста |
| 3 | Цель и задачи | `diagrams/goals.png` |
| 4 | Архитектура | `diagrams/architecture.png` |
| 5 | Агент + API | скрин `04-3-agent-api.png`, код из `agent/main.go` |
| 6 | UI | `E1-Obzor.png`, `E6-Grafiki.png` |
| 7 | Алерты | `E5-Alerty.png`, код `AlertService` |
| 8 | Git + слои | `diagrams/layers_git.png` |
| 9 | Тестирование | миниатюры тестов + баннер внизу |

Бэкап оригинала: `Desktopnaya-sistema-monitoringa-klastera-serverov.backup.pptx`

## Пересборка после правок

```bash
python3 scripts/fix-presentation.py
```

Скриншоты берутся из `/home/lesha/ClusterGuard-Otchet-Materialy/`.

## Перед защитой

1. Замени `[ФИО]`, группу, руководителя на слайде 1.
2. Открой `.pptx` в LibreOffice Impress или PowerPoint и проверь вёрстку.
3. Прогони доклад с таймером (5–7 мин).
