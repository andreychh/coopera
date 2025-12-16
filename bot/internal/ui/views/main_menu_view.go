package views

import (
	"github.com/andreychh/coopera-bot/internal/ui/protocol"
	"github.com/andreychh/coopera-bot/pkg/botlib/content"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/formatting"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards"
	"github.com/andreychh/coopera-bot/pkg/botlib/content/keyboards/buttons"
	"github.com/andreychh/coopera-bot/pkg/botlib/sources"
)

func MainMenuView() sources.Source[content.Content] {
	text := `📂 <b>Главное меню</b>

Здесь вы можете управлять командами, задачами и отслеживать личный прогресс.

👇 Выберите нужный раздел:`
	return sources.Static[content.Content](
		keyboards.Inline(
			formatting.Formatted(content.Text(text), formatting.ParseModeHTML),
			buttons.Matrix(
				buttons.Row(buttons.CallbackButton("👥 Мои команды", protocol.ToTeamsMenu())),
				buttons.Row(buttons.CallbackButton("📋 Мои задачи", protocol.ToUserTasksMenu())),
				buttons.Row(buttons.CallbackButton("📊 Статистика", protocol.ToUserStatsMenu())),
			),
		),
	)
}
