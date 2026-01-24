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
