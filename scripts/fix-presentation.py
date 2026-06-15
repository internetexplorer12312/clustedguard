#!/usr/bin/env python3
"""Обновление презентации ВКР: схемы, скриншоты, код."""
from __future__ import annotations

import re
import shutil
import subprocess
import zipfile
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

REPO = Path("/home/lesha/clusterguard")
ASSETS = Path("/home/lesha/ClusterGuard-Otchet-Materialy")
PPTX_IN = REPO / "Desktopnaya-sistema-monitoringa-klastera-serverov.pptx"
PPTX_OUT = REPO / "Desktopnaya-sistema-monitoringa-klastera-serverov.pptx"
BACKUP = REPO / "Desktopnaya-sistema-monitoringa-klastera-serverov.backup.pptx"
WORK = Path("/tmp/pptx-clusterguard")
MEDIA = WORK / "ppt/media"
DIAG = WORK / "diagrams"


def font(size: int, bold: bool = False):
    paths = [
        "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf" if bold else "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
        "/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf" if bold else "/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
    ]
    for p in paths:
        if Path(p).exists():
            return ImageFont.truetype(p, size)
    return ImageFont.load_default()


def box(d: ImageDraw.ImageDraw, xy, text, fill, outline=(37, 99, 235), tc=(255, 255, 255), f=None):
    d.rounded_rectangle(xy, radius=12, fill=fill, outline=outline, width=2)
    x0, y0, x1, y1 = xy
    if f is None:
        f = font(18)
    lines = text.split("\n")
    lh = 22
    total_h = len(lines) * lh
    ty = y0 + (y1 - y0 - total_h) / 2
    for line in lines:
        tw = d.textlength(line, font=f)
        d.text((x0 + (x1 - x0 - tw) / 2, ty), line, fill=tc, font=f)
        ty += lh


def arrow(d, p1, p2, color=(37, 99, 235)):
    d.line([p1, p2], fill=color, width=3)
    import math

    ang = math.atan2(p2[1] - p1[1], p2[0] - p1[0])
    sz = 12
    a1 = ang + math.pi * 0.85
    a2 = ang - math.pi * 0.85
    d.polygon(
        [
            p2,
            (p2[0] + sz * math.cos(a1), p2[1] + sz * math.sin(a1)),
            (p2[0] + sz * math.cos(a2), p2[1] + sz * math.sin(a2)),
        ],
        fill=color,
    )


def save_fit(img: Image.Image, path: Path, size: tuple[int, int]):
    path.parent.mkdir(parents=True, exist_ok=True)
    img.convert("RGB").resize(size, Image.Resampling.LANCZOS).save(path, "PNG")


def copy_fit(src: Path, dst: Path, size: tuple[int, int]):
    save_fit(Image.open(src), dst, size)


def diagram_context(path: Path):
    w, h = 1528, 1528
    img = Image.new("RGB", (w, h), (248, 250, 252))
    d = ImageDraw.Draw(img)
    d.text((60, 40), "Контекст системы ClusterGuard", fill=(30, 41, 59), font=font(32, True))
    box(d, (540, 220, 980, 340), "Системный\nадминистратор", (30, 64, 175))
    box(d, (540, 520, 980, 700), "ClusterGuard\nDesktop (Wails)", (37, 99, 235))
    for i, name in enumerate(["Сервер 1\nАгент :9100", "Сервер 2\nАгент :9100", "Сервер N\nАгент :9100"]):
        y = 860 + i * 150
        box(d, (200 + i * 30, y, 520 + i * 30, y + 110), name, (22, 163, 74))
    arrow(d, (760, 340), (760, 520))
    arrow(d, (700, 700), (380, 860))
    arrow(d, (760, 700), (760, 860))
    arrow(d, (820, 700), (1100, 860))
    d.text((60, h - 60), "HTTP: /health, /metrics + X-ClusterGuard-Token", fill=(100, 116, 139), font=font(20))
    save_fit(img, path, (1528, 1528))


def diagram_architecture(path: Path):
    w, h = 3454, 550
    img = Image.new("RGB", (w, h), (255, 255, 255))
    d = ImageDraw.Draw(img)
    d.text((30, 15), "Архитектура ClusterGuard", fill=(30, 41, 59), font=font(28, True))
    layers = [
        ("Frontend", "TypeScript, Wails UI", (59, 130, 246)),
        ("Backend Go", "app.go, services", (99, 102, 241)),
        ("LevelDB", "~/.config/clusterguard", (139, 92, 246)),
    ]
    x = 80
    for title, sub, col in layers:
        box(d, (x, 120, x + 520, 280), f"{title}\n{sub}", col, f=font(22))
        if x > 80:
            arrow(d, (x - 40, 200), (x, 200))
        x += 620
    box(d, (2550, 120, 3350, 280), "Агенты на серверах\nGET /health · /metrics", (22, 163, 74), f=font(22))
    arrow(d, (2180, 200), (2550, 200))
    d.text((80, 380), "Опрос каждые 30 с  →  метрики  →  пороги  →  алерты в UI", fill=(71, 85, 105), font=font(22))
    save_fit(img, path, (3454, 550))


def diagram_layers_git(path: Path):
    w, h = 3442, 900
    img = Image.new("RGB", (w, h), (255, 255, 255))
    d = ImageDraw.Draw(img)
    d.text((30, 10), "Слои backend + Git Flow", fill=(30, 41, 59), font=font(26, True))
    for i, (t, s) in enumerate(
        [
            ("domain", "Server, Alert, Cluster"),
            ("service", "Monitor, Alert, Collector"),
            ("repository", "LevelDB + HTTP agent"),
        ]
    ):
        y = 70 + i * 95
        box(d, (60, y, 1500, y + 70), f"{t}\n{s}", (37, 99, 235), f=font(18))
    git = ASSETS / "git/Risunok-3-22-git-graph.png"
    if git.exists():
        g = Image.open(git).convert("RGB")
        g = g.resize((1900, int(1900 * g.height / g.width)), Image.Resampling.LANCZOS)
        img.paste(g, (1550, 55))
    else:
        d.text((1600, 200), "git log --graph --all", fill=(30, 41, 59), font=font(24))
    d.text((60, h - 40), "main (v1.0.x) · develop · feature/* · hotfix/*", fill=(100, 116, 139), font=font(18))
    save_fit(img, path, (3442, 568))


def diagram_goals(path: Path):
    w, h = 1114, 1114
    img = Image.new("RGB", (w, h), (239, 246, 255))
    d = ImageDraw.Draw(img)
    d.text((80, 60), "Результат разработки", fill=(30, 41, 59), font=font(34, True))
    items = [
        "Desktop-приложение\n(Wails + Go)",
        "Агент на каждой ноде\n(HTTP :9100)",
        "Мониторинг + алерты\n(CPU, RAM, диск)",
        "Локальное хранение\n(LevelDB)",
    ]
    y = 200
    for t in items:
        box(d, (120, y, 990, y + 150), t, (37, 99, 235), f=font(24))
        y += 190
    save_fit(img, path, (1114, 1114))


def code_image(path: Path, title: str, lines: list[str]):
    w, h = 1144, 1144
    img = Image.new("RGB", (w, h), (15, 17, 23))
    d = ImageDraw.Draw(img)
    d.rectangle((0, 0, w, 50), fill=(37, 99, 235))
    d.text((20, 12), title, fill=(255, 255, 255), font=font(20, True))
    y = 70
    f = font(17)
    for line in lines:
        col = (248, 250, 252)
        if "if " in line or "return" in line:
            col = (147, 197, 253)
        if "unauthorized" in line or "401" in line:
            col = (248, 113, 113)
        d.text((24, y), line, fill=col, font=f)
        y += 28
    save_fit(img, path, (1144, 1144))


def title_bg(path: Path, subtitle: str = ""):
    w, h = 1402, 2104
    img = Image.new("RGB", (w, h), (15, 23, 42))
    d = ImageDraw.Draw(img)
    d.rounded_rectangle((80, 500, w - 80, 1500), radius=30, fill=(30, 41, 59))
    d.text((120, 620), "Cluster", fill=(255, 255, 255), font=font(72, True))
    d.text((120, 720), "Guard", fill=(59, 130, 246), font=font(72, True))
    d.text((120, 860), "Мониторинг кластера серверов", fill=(203, 213, 225), font=font(28))
    if subtitle:
        d.text((120, 920), subtitle, fill=(148, 163, 184), font=font(22))
    save_fit(img, path, (1402, 2104))


def icon_thumb(src: Path, dst: Path):
    """Миниатюра для слайда 9 (~400px квадрат)."""
    copy_fit(src, dst, (400, 400))


def patch_slide_text(slide_path: Path, replacements: dict[str, str]):
    xml = slide_path.read_text(encoding="utf-8")
    for old, new in replacements.items():
        xml = xml.replace(old, new)
    slide_path.write_text(xml, encoding="utf-8")


NS_A = "http://schemas.openxmlformats.org/drawingml/2006/main"


def replace_code_textbox(slide_path: Path, needle: str, new_text: str):
    """Заменить содержимое текстового блока с кодом (отдельный shape)."""
    import xml.etree.ElementTree as ET

    root = ET.parse(slide_path).getroot()
    changed = False
    for tx in root.iter():
        if not tx.tag.endswith("txBody"):
            continue
        full = "".join((t.text or "") for t in tx.iter() if t.tag.endswith("t"))
        if needle not in full:
            continue
        for child in list(tx):
            if child.tag.endswith("p"):
                tx.remove(child)
        code_p = ET.SubElement(tx, f"{{{NS_A}}}p")
        r = ET.SubElement(code_p, f"{{{NS_A}}}r")
        rpr = ET.SubElement(r, f"{{{NS_A}}}rPr")
        rpr.set("lang", "en-US")
        rpr.set("sz", "1050")
        ET.SubElement(r, f"{{{NS_A}}}t").text = new_text
        changed = True
        break
    if changed:
        ET.register_namespace("a", NS_A)
        ET.register_namespace("p", "http://schemas.openxmlformats.org/presentationml/2006/main")
        ET.register_namespace("r", "http://schemas.openxmlformats.org/officeDocument/2006/relationships")
        ET.ElementTree(root).write(slide_path, encoding="UTF-8", xml_declaration=True)


def testing_banner(path: Path):
    """Широкий коллаж для нижней части слайда 9."""
    imgs = [
        ASSETS / "screenshots/04-2-modulnye-testy.png",
        ASSETS / "screenshots/04-3-agent-api.png",
        ASSETS / "desktop-testing/04-FS-02-katalog-dannyh.png",
    ]
    w, h = 2400, 500
    canvas = Image.new("RGB", (w, h), (240, 253, 244))
    d = ImageDraw.Draw(canvas)
    d.text((20, 10), "Результаты тестирования ClusterGuard", fill=(22, 101, 52), font=font(24, True))
    x = 20
    for src in imgs:
        if not src.exists():
            continue
        im = Image.open(src).convert("RGB")
        im.thumbnail((760, 420), Image.Resampling.LANCZOS)
        canvas.paste(im, (x, 60))
        x += im.width + 20
    save_fit(canvas, path, (2400, 500))


def main():
    if not PPTX_IN.exists():
        raise SystemExit(f"Не найден файл: {PPTX_IN}")

    shutil.copy2(PPTX_IN, BACKUP)
    if WORK.exists():
        shutil.rmtree(WORK)
    WORK.mkdir(parents=True)
    with zipfile.ZipFile(PPTX_IN, "r") as z:
        z.extractall(WORK)

    DIAG.mkdir(exist_ok=True)

    # --- Схемы ---
    diagram_context(DIAG / "context.png")
    diagram_architecture(DIAG / "architecture.png")
    diagram_layers_git(DIAG / "layers_git.png")
    diagram_goals(DIAG / "goals.png")
    title_bg(DIAG / "title.png", "ВКР · 2025")

    code_image(
        DIAG / "agent_code.png",
        "agent/main.go — проверка токена",
        [
            'mux.HandleFunc("/metrics", func(w, r) {',
            '    if *token != "" &&',
            '       r.Header.Get("X-ClusterGuard-Token") != *token {',
            '        http.Error(w, "unauthorized", 401)',
            '        return',
            "    }",
            "    // collectMetrics() → JSON",
            "})",
        ],
    )

    # --- Подстановка в media ---
    mapping = {
        "image-1-1.png": DIAG / "title.png",
        "image-2-1.png": DIAG / "context.png",
        "image-3-1.png": DIAG / "goals.png",
        "image-4-1.png": DIAG / "architecture.png",
        "image-5-1.png": ASSETS / "screenshots/04-3-agent-api.png",
        "image-6-1.png": ASSETS / "prilozhenie-e/E1-Obzor.png",
        "image-6-2.png": ASSETS / "prilozhenie-e/E6-Grafiki.png",
        "image-7-1.png": ASSETS / "prilozhenie-e/E5-Alerty.png",
        "image-8-1.png": DIAG / "layers_git.png",
        "image-9-1.png": ASSETS / "screenshots/04-2-modulnye-testy.png",
        "image-9-3.png": ASSETS / "screenshots/04-3-agent-api.png",
        "image-9-5.png": ASSETS / "screenshots/04-5-docker-compose.png",
        "image-10-1.png": DIAG / "title.png",
    }

    sizes = {
        "image-1-1.png": (1402, 2104),
        "image-2-1.png": (1528, 1528),
        "image-3-1.png": (1114, 1114),
        "image-4-1.png": (3454, 550),
        "image-5-1.png": (1144, 1144),
        "image-6-1.png": (1280, 800),
        "image-6-2.png": (1280, 800),
        "image-7-1.png": (1280, 800),
        "image-8-1.png": (3442, 568),
        "image-9-1.png": (400, 400),
        "image-9-3.png": (400, 400),
        "image-9-5.png": (400, 400),
        "image-10-1.png": (1402, 2104),
    }

    for name, src in mapping.items():
        dst = MEDIA / name
        if not src.exists():
            print("WARN missing", src)
            continue
        if name.startswith("image-9-"):
            icon_thumb(src, dst)
        else:
            copy_fit(src, dst, sizes[name])
        print("OK", name, "<-", src.name)

    testing_banner(DIAG / "testing_banner.png")
    save_fit(Image.open(DIAG / "testing_banner.png"), MEDIA / "image-9-7.png", (2400, 500))

    # Стрелка слайд 2
    arr = Image.new("RGBA", (54, 44), (0, 0, 0, 0))
    d = ImageDraw.Draw(arr)
    arrow(d, (5, 22), (48, 22), (37, 99, 235))
    arr.convert("RGB").save(MEDIA / "image-2-2.png")

    # --- Текст слайдов: реальный код ---
    replace_code_textbox(
        WORK / "ppt/slides/slide5.xml",
        "AGENT_TOKEN",
        'if *token != "" && r.Header.Get("X-ClusterGuard-Token") != *token {\n    http.Error(w, "unauthorized", 401)\n    return\n}',
    )
    replace_code_textbox(
        WORK / "ppt/slides/slide7.xml",
        "checkThreshold",
        "if c.value >= c.threshold {\n    notifier.Notify(alert)\n    repo.Save(alert)\n}",
    )
    patch_slide_text(
        WORK / "ppt/slides/slide6.xml",
        {
            "экран «Серверы» (E2)": "экран с графиками метрик (E6)",
            "E1-Obzor) и экран «Серверы» (E2)": "E1-Obzor) и графики метрик (E6-Grafiki)",
            "экран «Серверы» (E6)": "графики метрик (E6-Grafiki)",
        },
    )

    # Увеличить баннер тестов на слайде 9
    s9 = (WORK / "ppt/slides/slide9.xml").read_text(encoding="utf-8")
    s9 = s9.replace(
        '<a:off x="609302" y="3602310"/>  <a:ext cx="169218" cy="135359"/>',
        '<a:off x="473943" y="3408759"/>  <a:ext cx="8196114" cy="744587"/>',
    )
    (WORK / "ppt/slides/slide9.xml").write_text(s9, encoding="utf-8")

    # --- Упаковка ---
    if PPTX_OUT.exists():
        PPTX_OUT.unlink()
    subprocess.run(
        ["zip", "-r", "-q", str(PPTX_OUT), "."],
        cwd=WORK,
        check=True,
    )
    diag_out = REPO / "presentation" / "diagrams"
    diag_out.mkdir(parents=True, exist_ok=True)
    for f in DIAG.glob("*.png"):
        shutil.copy2(f, diag_out / f.name)

    print("\nГотово:", PPTX_OUT)
    print("Бэкап:", BACKUP)
    print("Схемы:", diag_out)


if __name__ == "__main__":
    main()
