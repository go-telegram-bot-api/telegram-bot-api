package tgbotapi

// ChatConfig is a base type for all chat identifiers
type ChatConfig struct {
	ChatID             int64
	ChannelUsername    string
	SuperGroupUsername string
}

func (base ChatConfig) params() (Params, error) {
	return base.paramsWithKey("chat_id")
}

func (base ChatConfig) paramsWithKey(key string) (Params, error) {
	params := make(Params)
	return params, params.AddFirstValid(key, base.ChatID, base.ChannelUsername, base.SuperGroupUsername)
}

// BaseChat is base type for all chat config types.
type BaseChat struct {
	ChatConfig
	MessageThreadID     int
	ProtectContent      bool
	ReplyMarkup         interface{}
	DisableNotification bool
	ReplyParameters     ReplyParameters
}

func (chat *BaseChat) params() (Params, error) {
	params, err := chat.ChatConfig.params()
	if err != nil {
		return params, err
	}

	params.AddNonZero("message_thread_id", chat.MessageThreadID)
	params.AddBool("disable_notification", chat.DisableNotification)
	params.AddBool("protect_content", chat.ProtectContent)

	err = params.AddInterface("reply_markup", chat.ReplyMarkup)
	if err != nil {
		return params, err
	}
	err = params.AddInterface("reply_parameters", chat.ReplyParameters)
	return params, err
}

// BaseFile is a base type for all file config types.
type BaseFile struct {
	BaseChat
	File RequestFileData
}

func (file BaseFile) params() (Params, error) {
	return file.BaseChat.params()
}

// BaseEdit is base type of all chat edits.
type BaseEdit struct {
	BaseChatMessage
	InlineMessageID string
	ReplyMarkup     *InlineKeyboardMarkup
}

func (edit BaseEdit) params() (Params, error) {
	params := make(Params)

	if edit.InlineMessageID != "" {
		params["inline_message_id"] = edit.InlineMessageID
	} else {
		p1, err := edit.BaseChatMessage.params()
		if err != nil {
			return params, err
		}
		params.Merge(p1)
	}

	err := params.AddInterface("reply_markup", edit.ReplyMarkup)

	return params, err
}

// BaseSpoiler is base type of structures with spoilers.
type BaseSpoiler struct {
	HasSpoiler bool
}

func (spoiler BaseSpoiler) params() (Params, error) {
	params := make(Params)

	if spoiler.HasSpoiler {
		params.AddBool("has_spoiler", true)
	}

	return params, nil
}

// BaseChatMessage is a base type for all messages in chats.
type BaseChatMessage struct {
	ChatConfig
	MessageID int
}

func (base BaseChatMessage) params() (Params, error) {
	params, err := base.ChatConfig.params()
	if err != nil {
		return params, err
	}
	params.AddNonZero("message_id", base.MessageID)

	return params, nil
}

// BaseChatMessages is a base type for all messages in chats.
type BaseChatMessages struct {
	ChatConfig
	MessageIDs []int
}

func (base BaseChatMessages) params() (Params, error) {
	params, err := base.ChatConfig.params()
	if err != nil {
		return params, err
	}
	err = params.AddInterface("message_ids", base.MessageIDs)

	return params, err
}
