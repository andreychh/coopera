package views

import (
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/formatting"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
)

func WelcomeMessage() sources.Source[content.Content] {
	text := `👋 <b>Добро пожаловать!</b>

Coopera — это платформа для командной работы, где вклад участников прозрачен и измерим.

🤖 <b>Что умеет система:</b>
• <b>Команды:</b> создание групп и добавление участников.
• <b>Задачи:</b> оценка стоимости в баллах.
• <b>Процесс:</b> назначение исполнителей и смена статусов.
• <b>Аналитика:</b> статистика вклада каждого участника.`
	return sources.Static(formatting.Formatted(content.Text(text), formatting.ParseModeHTML))
}
