// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package actions

import (
	"context"
	"fmt"

	"github.com/andreychh/coopera/internal/buildinfo"
	"github.com/andreychh/coopera/pkg/bot/api"
	"github.com/andreychh/coopera/pkg/bot/endpoints"
)

type RequestSender interface {
	SendRequest(ctx context.Context, method api.Method, reqBody, respBody any) error
}

type DisplayBuildInfoAction struct {
	client RequestSender
}

func (a DisplayBuildInfoAction) Execute(ctx context.Context, update api.Update) error {
	if update.Message == nil {
		return nil
	}
	_, err := endpoints.SendMessage(a.client).Call(
		ctx,
		api.SendMessageRequest{
			ChatID: api.ChatID{
				ChatID:          new(update.Message.Chat.ID),
				ChannelUsername: nil,
			},
			Text:      a.text(buildinfo.Read()),
			ParseMode: api.ParseModeHTML,
		},
	)
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}
	return nil
}

func (a DisplayBuildInfoAction) text(info buildinfo.BuildInfo) string {
	return fmt.Sprintf(`<b>Build Information</b>
<pre>
Version:    %s
Commit:     %s
Branch:     %s
Build Time: %s
Tree State: %s
Built By:   %s
Go Version: %s
Platform:   %s
</pre>`,
		info.Version,
		info.Commit,
		info.Branch,
		info.BuildTime,
		info.TreeState,
		info.BuiltBy,
		info.GoVersion,
		info.Platform,
	)
}

func DisplayBuildInfo(client RequestSender) DisplayBuildInfoAction {
	return DisplayBuildInfoAction{
		client: client,
	}
}
