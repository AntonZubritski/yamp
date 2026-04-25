# TOOLS.md - Local Notes

Skills define _how_ tools work. This file is for _your_ specifics — the stuff that's unique to your setup.

## What Goes Here

Things like:

- Camera names and locations
- SSH hosts and aliases
- Preferred voices for TTS
- Speaker/room names
- Device nicknames
- Anything environment-specific

## MashaBot Tool Reference

### 🔍 Information & Internet
- **search** — поиск в интернете (Tavily, до 10 результатов)
- **web_fetch** — извлечение текста с URL

### 🎨 Генерация изображений
- **imagine** — Nano Banana Pro, 1 картинка, 10-30 сек (быстро)
- **imagine_mj** — Midjourney, 4 художественных, 2-10 мин (нестабилен)
- **image_edit** — редактирование/трансформация изображения

### 🎵 Медиа
- **speak** — озвучка текста (TTS), ElevenLabs
- **transcribe** — расшифровка аудио/голосовых (Whisper)
- **generate_video** — генерация видео из текста/изображения

### 💰 Баланс
- **get_balance** — баланс токенов, подписка, дневной бюджет

### 📝 MashaNote (Заметки и Базы данных)
- **notes_list_pages** — дерево страниц
- **notes_search** — поиск по заголовкам и содержимому
- **notes_create_page** — создать страницу (Markdown, автоматом конвертируется)
- **notes_read_page** — прочитать страницу по ID
- **notes_update_page** — обновить заголовок/содержимое
- **notes_delete_page** — удалить (в корзину)
- **notes_share_page** — публичная ссылка
- **notes_get_link** — сгенерировать ссылку без API-вызова
- **notes_list_databases** — все базы данных
- **notes_create_database** — создать БД с колонками
- **notes_get_database** — схема + строки
- **notes_add_row** — добавить строку
- **notes_update_row** — обновить строку
- **notes_delete_row** — удалить строку
- **notes_query_database** — фильтр/сортировка

### 📁 Работа с файлами пользователя
Файлы сохраняются в `~/.openclaw/media/inbound/`. Читать через exec:

| Формат | Команда |
|--------|---------|
| Excel (.xlsx) | `python3 -c "import openpyxl; wb=openpyxl.load_workbook('FILE'); ws=wb.active; [print(list(row)) for row in ws.iter_rows(values_only=True)]"` |
| CSV | `cat FILE` или `python3 -c "import csv; [print(r) for r in csv.reader(open('FILE'))]"` |
| JSON | `python3 -c "import json; print(json.dumps(json.load(open('FILE')), indent=2, ensure_ascii=False))"` |
| PDF | `python3 -c "from PyPDF2 import PdfReader; r=PdfReader('FILE'); [print(p.extract_text()) for p in r.pages]"` |
| TXT/MD/LOG | `cat FILE` |
| Изображения | vision-модель автоматом анализирует |

### 📋 Правила создания БД в MashaNote
1. **Первое свойство** массива `properties` — всегда TEXT (заголовок строки)
   - Пример: `{name: "Название", type: "TEXT"}`
2. Типы: `TEXT`, `DATE`, `NUMBER`, `CHECKBOX`, `SELECT`, `MULTI_SELECT`, `STATUS`, `URL`, `EMAIL`
3. Для SELECT/MULTI_SELECT/STATUS: `options: {choices: [{name: "...", color: "..."}]}`
4. Доступные цвета: `default`, `gray`, `brown`, `orange`, `yellow`, `green`, `blue`, `purple`, `pink`, `red`
5. **Никогда не говори пользователю что БД или контент "не работают" — они работают**
6. При ошибке 401/403 Notes API — повтори через 5 секунд (кэш сбросится)

### 🚫 Ошибки, которые нельзя повторять
- **imagine / imagine_mj / image_edit** — это встроенные инструменты. Не искать их в интернете, не говорить "нет доступа"
- **Голосовое сообщение** → всегда transcribe (не читать вручную)
- **Файл от пользователя** → всегда читать через exec (не говорить "не могу прочитать")
- **No REPL** — следить за контекстом, не отправлять NO_REPLY когда пользователь ждёт ответ

---

Add whatever helps you do your job. This is your cheat sheet.
