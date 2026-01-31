// SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package api

type Method string

const (
	MethodSendMessage     Method = "sendMessage"
	MethodEditMessageText Method = "editMessageText"
	MethodAnswerCallback  Method = "answerCallbackQuery"
	MethodGetMe           Method = "getMe"
	MethodGetUpdates      Method = "getUpdates"
)

type ChatType string

const (
	ChatTypeUnknown    ChatType = ""
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeSuperGroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)

type ParseMode string

const (
	ParseModeUnknown    ParseMode = ""
	ParseModeHTML       ParseMode = "HTML"
	ParseModeMarkdown   ParseMode = "Markdown"
	ParseModeMarkdownV2 ParseMode = "MarkdownV2"
)
