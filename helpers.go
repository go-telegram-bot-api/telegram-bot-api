package tgbotapi

import (
	"net/url"
)

// Creates a new Message.
// Perhaps set a ChatAction of CHAT_TYPING while processing.
//
// chatId is where to send it, text is the message text.
func NewMessage(chatId int, text string) MessageConfig {
	return MessageConfig{
		ChatId: chatId,
		Text:   text,
		DisableWebPagePreview: false,
		ReplyToMessageId:      0,
	}
}

// Creates a new forward.
//
// chatId is where to send it, fromChatId is the source chat,
// and messageId is the Id of the original message.
func NewForward(chatId int, fromChatId int, messageId int) ForwardConfig {
	return ForwardConfig{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageId:  messageId,
	}
}

// Creates a new photo uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of CHAT_UPLOAD_PHOTO while processing.
//
// chatId is where to send it, filename is the path to the file.
func NewPhotoUpload(chatId int, filename string) PhotoConfig {
	return PhotoConfig{
		ChatId:           chatId,
		UseExistingPhoto: false,
		FilePath:         filename,
	}
}

// Shares an existing photo.
// You may use this to reshare an existing photo without reuploading it.
//
// chatId is where to send it, fileId is the Id of the file already uploaded.
func NewPhotoShare(chatId int, fileId string) PhotoConfig {
	return PhotoConfig{
		ChatId:           chatId,
		UseExistingPhoto: true,
		FileId:           fileId,
	}
}

// Creates a new audio uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of CHAT_RECORD_AUDIO or CHAT_UPLOAD_AUDIO while processing.
//
// chatId is where to send it, filename is the path to the file.
func NewAudioUpload(chatId int, filename string) AudioConfig {
	return AudioConfig{
		ChatId:           chatId,
		UseExistingAudio: false,
		FilePath:         filename,
	}
}

// Shares an existing audio file.
// You may use this to reshare an existing audio file without reuploading it.
//
// chatId is where to send it, fileId is the Id of the audio already uploaded.
func NewAudioShare(chatId int, fileId string) AudioConfig {
	return AudioConfig{
		ChatId:           chatId,
		UseExistingAudio: true,
		FileId:           fileId,
	}
}

// Creates a new document uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of CHAT_UPLOAD_DOCUMENT while processing.
//
// chatId is where to send it, filename is the path to the file.
func NewDocumentUpload(chatId int, filename string) DocumentConfig {
	return DocumentConfig{
		ChatId:              chatId,
		UseExistingDocument: false,
		FilePath:            filename,
	}
}

// Shares an existing document.
// You may use this to reshare an existing document without reuploading it.
//
// chatId is where to send it, fileId is the Id of the document already uploaded.
func NewDocumentShare(chatId int, fileId string) DocumentConfig {
	return DocumentConfig{
		ChatId:              chatId,
		UseExistingDocument: true,
		FileId:              fileId,
	}
}

// Creates a new sticker uploader.
// This requires a file on the local filesystem to upload to Telegram.
//
// chatId is where to send it, filename is the path to the file.
func NewStickerUpload(chatId int, filename string) StickerConfig {
	return StickerConfig{
		ChatId:             chatId,
		UseExistingSticker: false,
		FilePath:           filename,
	}
}

// Shares an existing sticker.
// You may use this to reshare an existing sticker without reuploading it.
//
// chatId is where to send it, fileId is the Id of the sticker already uploaded.
func NewStickerShare(chatId int, fileId string) StickerConfig {
	return StickerConfig{
		ChatId:             chatId,
		UseExistingSticker: true,
		FileId:             fileId,
	}
}

// Creates a new video uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of CHAT_RECORD_VIDEO or CHAT_UPLOAD_VIDEO while processing.
//
// chatId is where to send it, filename is the path to the file.
func NewVideoUpload(chatId int, filename string) VideoConfig {
	return VideoConfig{
		ChatId:           chatId,
		UseExistingVideo: false,
		FilePath:         filename,
	}
}

// Shares an existing video.
// You may use this to reshare an existing video without reuploading it.
//
// chatId is where to send it, fileId is the Id of the video already uploaded.
func NewVideoShare(chatId int, fileId string) VideoConfig {
	return VideoConfig{
		ChatId:           chatId,
		UseExistingVideo: true,
		FileId:           fileId,
	}
}

// Shares your location.
// Perhaps set a ChatAction of CHAT_FIND_LOCATION while processing.
//
// chatId is where to send it, latitude and longitude are coordinates.
func NewLocation(chatId int, latitude float64, longitude float64) LocationConfig {
	return LocationConfig{
		ChatId:           chatId,
		Latitude:         latitude,
		Longitude:        longitude,
		ReplyToMessageId: 0,
		ReplyMarkup:      nil,
	}
}

// Sets a chat action.
// Actions last for 5 seconds, or until your next action.
//
// chatId is where to send it, action should be set via CHAT constants.
func NewChatAction(chatId int, action string) ChatActionConfig {
	return ChatActionConfig{
		ChatId: chatId,
		Action: action,
	}
}

// Gets user profile photos.
//
// userId is the Id of the user you wish to get profile photos from.
func NewUserProfilePhotos(userId int) UserProfilePhotosConfig {
	return UserProfilePhotosConfig{
		UserId: userId,
		Offset: 0,
		Limit:  0,
	}
}

// Gets updates since the last Offset.
//
// offset is the last Update Id to include.
// You likely want to set this to the last Update Id plus 1.
func NewUpdate(offset int) UpdateConfig {
	return UpdateConfig{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
}

// Creates a new webhook.
//
// link is the url parsable link you wish to get the updates.
func NewWebhook(link string) WebhookConfig {
	u, _ := url.Parse(link)

	return WebhookConfig{
		Url:   u,
		Clear: false,
	}
}
