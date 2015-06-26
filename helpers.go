package tgbotapi

import (
	"net/url"
)

func NewMessage(chatId int, text string) MessageConfig {
	return MessageConfig{
		ChatId: chatId,
		Text:   text,
		DisableWebPagePreview: false,
		ReplyToMessageId:      0,
	}
}

func NewForward(chatId int, fromChatId int, messageId int) ForwardConfig {
	return ForwardConfig{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageId:  messageId,
	}
}

func NewPhotoUpload(chatId int, filename string) PhotoConfig {
	return PhotoConfig{
		ChatId:           chatId,
		UseExistingPhoto: false,
		FilePath:         filename,
	}
}

func NewPhotoShare(chatId int, fileId string) PhotoConfig {
	return PhotoConfig{
		ChatId:           chatId,
		UseExistingPhoto: true,
		FileId:           fileId,
	}
}

func NewAudioUpload(chatId int, filename string) AudioConfig {
	return AudioConfig{
		ChatId:           chatId,
		UseExistingAudio: false,
		FilePath:         filename,
	}
}

func NewAudioShare(chatId int, fileId string) AudioConfig {
	return AudioConfig{
		ChatId:           chatId,
		UseExistingAudio: true,
		FileId:           fileId,
	}
}

func NewDocumentUpload(chatId int, filename string) DocumentConfig {
	return DocumentConfig{
		ChatId:              chatId,
		UseExistingDocument: false,
		FilePath:            filename,
	}
}

func NewDocumentShare(chatId int, fileId string) DocumentConfig {
	return DocumentConfig{
		ChatId:              chatId,
		UseExistingDocument: true,
		FileId:              fileId,
	}
}

func NewStickerUpload(chatId int, filename string) StickerConfig {
	return StickerConfig{
		ChatId:             chatId,
		UseExistingSticker: false,
		FilePath:           filename,
	}
}

func NewStickerShare(chatId int, fileId string) StickerConfig {
	return StickerConfig{
		ChatId:             chatId,
		UseExistingSticker: true,
		FileId:             fileId,
	}
}

func NewVideoUpload(chatId int, filename string) VideoConfig {
	return VideoConfig{
		ChatId:           chatId,
		UseExistingVideo: false,
		FilePath:         filename,
	}
}

func NewVideoShare(chatId int, fileId string) VideoConfig {
	return VideoConfig{
		ChatId:           chatId,
		UseExistingVideo: true,
		FileId:           fileId,
	}
}

func NewLocation(chatId int, latitude float64, longitude float64) LocationConfig {
	return LocationConfig{
		ChatId:           chatId,
		Latitude:         latitude,
		Longitude:        longitude,
		ReplyToMessageId: 0,
		ReplyMarkup:      nil,
	}
}

func NewChatAction(chatId int, action string) ChatActionConfig {
	return ChatActionConfig{
		ChatId: chatId,
		Action: action,
	}
}

func NewUserProfilePhotos(userId int) UserProfilePhotosConfig {
	return UserProfilePhotosConfig{
		UserId: userId,
		Offset: 0,
		Limit:  0,
	}
}

func NewUpdate(offset int) UpdateConfig {
	return UpdateConfig{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
}

func NewWebhook(link string) WebhookConfig {
	u, _ := url.Parse(link)

	return WebhookConfig{
		Url:   u,
		Clear: false,
	}
}
