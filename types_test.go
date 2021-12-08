package tgbotapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := Error{
		Message: "errorMessage",
	}

	assert.Equal(t, "errorMessage", err.Error())
}

func TestUpdateSentFrom(t *testing.T) {
	cases := []struct {
		Update Update
		User   *User
	}{{
		Update: Update{Message: &Message{
			From: &User{ID: 123932},
		}},
		User: &User{ID: 123932},
	}, {
		Update: Update{Message: &Message{}},
		User:   nil,
	}, {
		Update: Update{EditedMessage: &Message{
			From: &User{ID: 123933},
		}},
		User: &User{ID: 123933},
	}, {
		Update: Update{EditedMessage: &Message{}},
		User:   nil,
	}, {
		Update: Update{InlineQuery: &InlineQuery{
			From: &User{ID: 123934},
		}},
		User: &User{ID: 123934},
	}, {
		Update: Update{ChosenInlineResult: &ChosenInlineResult{
			From: &User{ID: 123935},
		}},
		User: &User{ID: 123935},
	}, {
		Update: Update{CallbackQuery: &CallbackQuery{
			From: &User{ID: 123936},
		}},
		User: &User{ID: 123936},
	}, {
		Update: Update{ShippingQuery: &ShippingQuery{
			From: &User{ID: 123937},
		}},
		User: &User{ID: 123937},
	}, {
		Update: Update{PreCheckoutQuery: &PreCheckoutQuery{
			From: &User{ID: 123938},
		}},
		User: &User{ID: 123938},
	}, {
		Update: Update{},
		User:   nil,
	}}

	for _, testCase := range cases {
		assert.Equal(t, testCase.User, testCase.Update.SentFrom())
	}
}

func TestUpdateCallbackData(t *testing.T) {
	assert.Equal(t, "", (&Update{}).CallbackData())
	assert.Equal(t, "", (&Update{CallbackQuery: &CallbackQuery{}}).CallbackData())
	assert.Equal(t, "data", (&Update{CallbackQuery: &CallbackQuery{Data: "data"}}).CallbackData())
}

func TestUpdateFromChat(t *testing.T) {
	cases := []struct {
		Update *Update
		Chat   *Chat
	}{{
		Update: &Update{Message: &Message{
			Chat: &Chat{ID: 12321},
		}},
		Chat: &Chat{ID: 12321},
	}, {
		Update: &Update{Message: &Message{}},
		Chat:   nil,
	}, {
		Update: &Update{EditedMessage: &Message{
			Chat: &Chat{ID: 12322},
		}},
		Chat: &Chat{ID: 12322},
	}, {
		Update: &Update{ChannelPost: &Message{
			Chat: &Chat{ID: 12323},
		}},
		Chat: &Chat{ID: 12323},
	}, {
		Update: &Update{EditedChannelPost: &Message{
			Chat: &Chat{ID: 12324},
		}},
		Chat: &Chat{ID: 12324},
	}, {
		Update: &Update{CallbackQuery: &CallbackQuery{
			Message: &Message{Chat: &Chat{ID: 12325}},
		}},
		Chat: &Chat{ID: 12325},
	}, {
		Update: &Update{CallbackQuery: &CallbackQuery{
			Message: nil,
		}},
		Chat: nil,
	}, {
		Update: &Update{},
		Chat:   nil,
	}}

	for _, testCase := range cases {
		assert.Equal(t, testCase.Chat, testCase.Update.FromChat())
	}
}

func TestUserStringNoUserName(t *testing.T) {
	user := User{
		ID:           0,
		FirstName:    "Test",
		LastName:     "Test",
		UserName:     "",
		LanguageCode: "en",
		IsBot:        false,
	}

	assert.Equal(t, "Test Test", user.String())
}

func TestUserStringWithUserName(t *testing.T) {
	user := User{
		ID:           0,
		FirstName:    "Test",
		LastName:     "Test",
		UserName:     "@test",
		LanguageCode: "en",
	}

	assert.Equal(t, "@test", user.String())
}

func TestUserStringNil(t *testing.T) {
	var user *User

	assert.Equal(t, "", user.String())
}

func TestMessageTime(t *testing.T) {
	message := Message{Date: 12345}

	assert.Equal(t, time.Unix(12345, 0), message.Time())
}

func TestMessageIsCommandWithCommand(t *testing.T) {
	message := Message{Text: "/command"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	assert.True(t, message.IsCommand())
}

func TestIsCommandWithText(t *testing.T) {
	message := Message{Text: "some text"}

	assert.False(t, message.IsCommand())
}

func TestIsCommandWithEmptyText(t *testing.T) {
	message := Message{Text: ""}

	assert.False(t, message.IsCommand())
}

func TestCommandWithCommand(t *testing.T) {
	message := Message{Text: "/command"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	assert.Equal(t, "command", message.Command())
}

func TestCommandWithEmptyText(t *testing.T) {
	message := Message{Text: ""}

	assert.Equal(t, "", message.Command())
}

func TestCommandWithNonCommand(t *testing.T) {
	message := Message{Text: "test text"}

	assert.Equal(t, "", message.Command())
}

func TestCommandWithBotName(t *testing.T) {
	message := Message{Text: "/command@test_bot"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 17}}

	assert.Equal(t, "command", message.Command())
}

func TestCommandWithAtWithBotName(t *testing.T) {
	message := Message{Text: "/command@test_bot"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 17}}

	assert.Equal(t, "command@test_bot", message.CommandWithAt())
}

func TestMessageCommandArgumentsWithArguments(t *testing.T) {
	message := Message{Text: "/command with arguments"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	assert.Equal(t, "with arguments", message.CommandArguments())
}

func TestMessageCommandArgumentsWithMalformedArguments(t *testing.T) {
	message := Message{Text: "/command-without argument space"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	assert.Equal(t, "without argument space", message.CommandArguments())
}

func TestMessageCommandArgumentsWithoutArguments(t *testing.T) {
	message := Message{Text: "/command"}

	assert.Equal(t, "", message.CommandArguments())
}

func TestMessageCommandArgumentsForNonCommand(t *testing.T) {
	message := Message{Text: "test text"}

	assert.Equal(t, "", message.CommandArguments())
}

func TestMessageCommandArgumentsWithFullLength(t *testing.T) {
	message := Message{Text: "/command"}
	message.Entities = []MessageEntity{{Type: "bot_command", Offset: 0, Length: 8}}

	assert.Equal(t, "", message.CommandArguments())
}

func TestMessageEntityParseURLGood(t *testing.T) {
	entity := MessageEntity{URL: "https://www.google.com"}

	url, err := entity.ParseURL()
	assert.NoError(t, err)
	assert.Equal(t, "https://www.google.com", url.String())
}

func TestMessageEntityParseURLBad(t *testing.T) {
	entity := MessageEntity{URL: ""}

	url, err := entity.ParseURL()
	assert.Error(t, err)
	assert.Nil(t, url)
}

func TestChatIsPrivate(t *testing.T) {
	chat := Chat{ID: 10, Type: "private"}

	assert.True(t, chat.IsPrivate())
	assert.False(t, chat.IsGroup())
	assert.False(t, chat.IsChannel())
	assert.False(t, chat.IsSuperGroup())
}

func TestChatIsGroup(t *testing.T) {
	chat := Chat{ID: 10, Type: "group"}

	assert.False(t, chat.IsPrivate())
	assert.True(t, chat.IsGroup())
	assert.False(t, chat.IsChannel())
	assert.False(t, chat.IsSuperGroup())
}

func TestChatIsChannel(t *testing.T) {
	chat := Chat{ID: 10, Type: "channel"}

	assert.False(t, chat.IsPrivate())
	assert.False(t, chat.IsGroup())
	assert.True(t, chat.IsChannel())
	assert.False(t, chat.IsSuperGroup())
}

func TestChatIsSuperGroup(t *testing.T) {
	chat := Chat{ID: 10, Type: "supergroup"}

	assert.False(t, chat.IsPrivate())
	assert.False(t, chat.IsGroup())
	assert.False(t, chat.IsChannel())
	assert.True(t, chat.IsSuperGroup())
}

func TestChatChatConfig(t *testing.T) {
	chat := Chat{ID: 10}

	assert.Equal(t, ChatConfig{ChatID: 10}, chat.ChatConfig())
}

func TestMessageEntityIsMention(t *testing.T) {
	entity := MessageEntity{Type: "mention"}

	assert.True(t, entity.IsMention())
}

func TestMessageEntityIsHashtag(t *testing.T) {
	entity := MessageEntity{Type: "hashtag"}

	assert.True(t, entity.IsHashtag())
}

func TestMessageEntityIsBotCommand(t *testing.T) {
	entity := MessageEntity{Type: "bot_command"}

	assert.True(t, entity.IsCommand())
}

func TestMessageEntityIsUrl(t *testing.T) {
	entity := MessageEntity{Type: "url"}

	assert.True(t, entity.IsURL())
}

func TestMessageEntityIsEmail(t *testing.T) {
	entity := MessageEntity{Type: "email"}

	assert.True(t, entity.IsEmail())
}

func TestMessageEntityIsBold(t *testing.T) {
	entity := MessageEntity{Type: "bold"}

	assert.True(t, entity.IsBold())
}

func TestMessageEntityIsItalic(t *testing.T) {
	entity := MessageEntity{Type: "italic"}

	assert.True(t, entity.IsItalic())
}

func TestMessageEntityIsCode(t *testing.T) {
	entity := MessageEntity{Type: "code"}

	assert.True(t, entity.IsCode())
}

func TestMessageEntityIsPre(t *testing.T) {
	entity := MessageEntity{Type: "pre"}

	assert.True(t, entity.IsPre())
}

func TestMessageEntityIsTextLink(t *testing.T) {
	entity := MessageEntity{Type: "text_link"}

	assert.True(t, entity.IsTextLink())
}

func TestVoiceChatScheduledTime(t *testing.T) {
	scheduledTime := VoiceChatScheduled{StartDate: 3784}

	assert.Equal(t, time.Unix(3784, 0), scheduledTime.Time())
}

func TestFileLink(t *testing.T) {
	file := File{FilePath: "test/test.txt"}

	if file.Link("token") != "https://api.telegram.org/file/bottoken/test/test.txt" {
		t.Fail()
	}
}

func TestChatMemberIsCreator(t *testing.T) {
	member := ChatMember{Status: "creator"}

	assert.True(t, member.IsCreator())
	assert.False(t, member.IsAdministrator())
	assert.False(t, member.HasLeft())
	assert.False(t, member.WasKicked())
}

func TestChatMemberIsAdministrator(t *testing.T) {
	member := ChatMember{Status: "administrator"}

	assert.False(t, member.IsCreator())
	assert.True(t, member.IsAdministrator())
	assert.False(t, member.HasLeft())
	assert.False(t, member.WasKicked())
}

func TestChatMemberHasLeft(t *testing.T) {
	member := ChatMember{Status: "left"}

	assert.False(t, member.IsCreator())
	assert.False(t, member.IsAdministrator())
	assert.True(t, member.HasLeft())
	assert.False(t, member.WasKicked())
}

func TestChatMemberWasKicked(t *testing.T) {
	member := ChatMember{Status: "kicked"}

	assert.False(t, member.IsCreator())
	assert.False(t, member.IsAdministrator())
	assert.False(t, member.HasLeft())
	assert.True(t, member.WasKicked())
}

func TestWebhookInfoIsSet(t *testing.T) {
	assert.True(t, WebhookInfo{URL: "test"}.IsSet())
	assert.False(t, WebhookInfo{}.IsSet())
}

// Ensure all configs are sendable
var (
	_ Chattable = AnimationConfig{}
	_ Chattable = AudioConfig{}
	_ Chattable = BanChatMemberConfig{}
	_ Chattable = BanChatSenderChatConfig{}
	_ Chattable = CallbackConfig{}
	_ Chattable = ChatActionConfig{}
	_ Chattable = ChatAdministratorsConfig{}
	_ Chattable = ChatInfoConfig{}
	_ Chattable = ChatInviteLinkConfig{}
	_ Chattable = CloseConfig{}
	_ Chattable = ContactConfig{}
	_ Chattable = CopyMessageConfig{}
	_ Chattable = CreateChatInviteLinkConfig{}
	_ Chattable = DeleteChatPhotoConfig{}
	_ Chattable = DeleteChatStickerSetConfig{}
	_ Chattable = DeleteMessageConfig{}
	_ Chattable = DeleteMyCommandsConfig{}
	_ Chattable = DeleteWebhookConfig{}
	_ Chattable = DocumentConfig{}
	_ Chattable = EditChatInviteLinkConfig{}
	_ Chattable = EditMessageCaptionConfig{}
	_ Chattable = EditMessageLiveLocationConfig{}
	_ Chattable = EditMessageMediaConfig{}
	_ Chattable = EditMessageReplyMarkupConfig{}
	_ Chattable = EditMessageTextConfig{}
	_ Chattable = FileConfig{}
	_ Chattable = ForwardConfig{}
	_ Chattable = GameConfig{}
	_ Chattable = GetChatMemberConfig{}
	_ Chattable = GetGameHighScoresConfig{}
	_ Chattable = InlineConfig{}
	_ Chattable = InvoiceConfig{}
	_ Chattable = KickChatMemberConfig{}
	_ Chattable = LeaveChatConfig{}
	_ Chattable = LocationConfig{}
	_ Chattable = LogOutConfig{}
	_ Chattable = MediaGroupConfig{}
	_ Chattable = MessageConfig{}
	_ Chattable = PhotoConfig{}
	_ Chattable = PinChatMessageConfig{}
	_ Chattable = PreCheckoutConfig{}
	_ Chattable = PromoteChatMemberConfig{}
	_ Chattable = RestrictChatMemberConfig{}
	_ Chattable = RevokeChatInviteLinkConfig{}
	_ Chattable = SendPollConfig{}
	_ Chattable = SetChatDescriptionConfig{}
	_ Chattable = SetChatPhotoConfig{}
	_ Chattable = SetChatTitleConfig{}
	_ Chattable = SetGameScoreConfig{}
	_ Chattable = ShippingConfig{}
	_ Chattable = StickerConfig{}
	_ Chattable = StopMessageLiveLocationConfig{}
	_ Chattable = StopPollConfig{}
	_ Chattable = UnbanChatMemberConfig{}
	_ Chattable = UnbanChatSenderChatConfig{}
	_ Chattable = UnpinChatMessageConfig{}
	_ Chattable = UpdateConfig{}
	_ Chattable = UserProfilePhotosConfig{}
	_ Chattable = VenueConfig{}
	_ Chattable = VideoConfig{}
	_ Chattable = VideoNoteConfig{}
	_ Chattable = VoiceConfig{}
	_ Chattable = WebhookConfig{}
)

// Ensure all Fileable types are correct.
var (
	_ Fileable = (*PhotoConfig)(nil)
	_ Fileable = (*AudioConfig)(nil)
	_ Fileable = (*DocumentConfig)(nil)
	_ Fileable = (*StickerConfig)(nil)
	_ Fileable = (*VideoConfig)(nil)
	_ Fileable = (*AnimationConfig)(nil)
	_ Fileable = (*VideoNoteConfig)(nil)
	_ Fileable = (*VoiceConfig)(nil)
	_ Fileable = (*SetChatPhotoConfig)(nil)
	_ Fileable = (*EditMessageMediaConfig)(nil)
	_ Fileable = (*SetChatPhotoConfig)(nil)
	_ Fileable = (*UploadStickerConfig)(nil)
	_ Fileable = (*NewStickerSetConfig)(nil)
	_ Fileable = (*AddStickerConfig)(nil)
	_ Fileable = (*MediaGroupConfig)(nil)
	_ Fileable = (*WebhookConfig)(nil)
	_ Fileable = (*SetStickerSetThumbConfig)(nil)
)

// Ensure all RequestFileData types are correct.
var (
	_ RequestFileData = (*FilePath)(nil)
	_ RequestFileData = (*FileBytes)(nil)
	_ RequestFileData = (*FileReader)(nil)
	_ RequestFileData = (*FileURL)(nil)
	_ RequestFileData = (*FileID)(nil)
	_ RequestFileData = (*fileAttach)(nil)
)
