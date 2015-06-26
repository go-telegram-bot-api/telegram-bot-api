package tgbotapi

import (
	"net/url"
)

// NewMessage creates a new Message.
// Perhaps set a ChatAction of ChatTyping while processing.
//
// chatID is where to send it, text is the message text.
func NewMessage(chatID int, text string) MessageConfig {
	return MessageConfig{
		ChatID: chatID,
		Text:   text,
		DisableWebPagePreview: false,
		ReplyToMessageID:      0,
	}
}

// NewForward creates a new forward.
//
// chatID is where to send it, fromChatID is the source chat,
// and messageID is the ID of the original message.
func NewForward(chatID int, fromChatID int, messageID int) ForwardConfig {
	return ForwardConfig{
		ChatID:     chatID,
		FromChatID: fromChatID,
		MessageID:  messageID,
	}
}

// NewPhotoUpload creates a new photo uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of ChatUploadPhoto while processing.
//
// chatID is where to send it, filename is the path to the file.
func NewPhotoUpload(chatID int, filename string) PhotoConfig {
	return PhotoConfig{
		ChatID:           chatID,
		UseExistingPhoto: false,
		FilePath:         filename,
	}
}

// NewPhotoShare shares an existing photo.
// You may use this to reshare an existing photo without reuploading it.
//
// chatID is where to send it, fileID is the ID of the file already uploaded.
func NewPhotoShare(chatID int, fileID string) PhotoConfig {
	return PhotoConfig{
		ChatID:           chatID,
		UseExistingPhoto: true,
		FileID:           fileID,
	}
}

// NewAudioUpload creates a new audio uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of ChatRecordAudio or ChatUploadAudio while processing.
//
// chatID is where to send it, filename is the path to the file.
func NewAudioUpload(chatID int, filename string) AudioConfig {
	return AudioConfig{
		ChatID:           chatID,
		UseExistingAudio: false,
		FilePath:         filename,
	}
}

// NewAudioShare shares an existing audio file.
// You may use this to reshare an existing audio file without reuploading it.
//
// chatID is where to send it, fileID is the ID of the audio already uploaded.
func NewAudioShare(chatID int, fileID string) AudioConfig {
	return AudioConfig{
		ChatID:           chatID,
		UseExistingAudio: true,
		FileID:           fileID,
	}
}

// NewDocumentUpload creates a new document uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of ChatUploadDocument while processing.
//
// chatID is where to send it, filename is the path to the file.
func NewDocumentUpload(chatID int, filename string) DocumentConfig {
	return DocumentConfig{
		ChatID:              chatID,
		UseExistingDocument: false,
		FilePath:            filename,
	}
}

// NewDocumentShare shares an existing document.
// You may use this to reshare an existing document without reuploading it.
//
// chatID is where to send it, fileID is the ID of the document already uploaded.
func NewDocumentShare(chatID int, fileID string) DocumentConfig {
	return DocumentConfig{
		ChatID:              chatID,
		UseExistingDocument: true,
		FileID:              fileID,
	}
}

// NewStickerUpload creates a new sticker uploader.
// This requires a file on the local filesystem to upload to Telegram.
//
// chatID is where to send it, filename is the path to the file.
func NewStickerUpload(chatID int, filename string) StickerConfig {
	return StickerConfig{
		ChatID:             chatID,
		UseExistingSticker: false,
		FilePath:           filename,
	}
}

// NewStickerShare shares an existing sticker.
// You may use this to reshare an existing sticker without reuploading it.
//
// chatID is where to send it, fileID is the ID of the sticker already uploaded.
func NewStickerShare(chatID int, fileID string) StickerConfig {
	return StickerConfig{
		ChatID:             chatID,
		UseExistingSticker: true,
		FileID:             fileID,
	}
}

// NewVideoUpload creates a new video uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of ChatRecordVideo or ChatUploadVideo while processing.
//
// chatID is where to send it, filename is the path to the file.
func NewVideoUpload(chatID int, filename string) VideoConfig {
	return VideoConfig{
		ChatID:           chatID,
		UseExistingVideo: false,
		FilePath:         filename,
	}
}

// NewVideoShare shares an existing video.
// You may use this to reshare an existing video without reuploading it.
//
// chatID is where to send it, fileID is the ID of the video already uploaded.
func NewVideoShare(chatID int, fileID string) VideoConfig {
	return VideoConfig{
		ChatID:           chatID,
		UseExistingVideo: true,
		FileID:           fileID,
	}
}

// NewLocation shares your location.
// Perhaps set a ChatAction of ChatFindLocation while processing.
//
// chatID is where to send it, latitude and longitude are coordinates.
func NewLocation(chatID int, latitude float64, longitude float64) LocationConfig {
	return LocationConfig{
		ChatID:           chatID,
		Latitude:         latitude,
		Longitude:        longitude,
		ReplyToMessageID: 0,
		ReplyMarkup:      nil,
	}
}

// NewChatAction sets a chat action.
// Actions last for 5 seconds, or until your next action.
//
// chatID is where to send it, action should be set via CHAT constants.
func NewChatAction(chatID int, action string) ChatActionConfig {
	return ChatActionConfig{
		ChatID: chatID,
		Action: action,
	}
}

// NewUserProfilePhotos gets user profile photos.
//
// userID is the ID of the user you wish to get profile photos from.
func NewUserProfilePhotos(userID int) UserProfilePhotosConfig {
	return UserProfilePhotosConfig{
		UserID: userID,
		Offset: 0,
		Limit:  0,
	}
}

// NewUpdate gets updates since the last Offset.
//
// offset is the last Update ID to include.
// You likely want to set this to the last Update ID plus 1.
func NewUpdate(offset int) UpdateConfig {
	return UpdateConfig{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
}

// NewWebhook creates a new webhook.
//
// link is the url parsable link you wish to get the updates.
func NewWebhook(link string) WebhookConfig {
	u, _ := url.Parse(link)

	return WebhookConfig{
		URL:   u,
		Clear: false,
	}
}
