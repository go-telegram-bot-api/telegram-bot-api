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

// NewFileUpload creates a new photo uploader.
// This requires a file on the local filesystem to upload to Telegram.
// Perhaps set a ChatAction of ChatUploadPhoto while processing.
//
// chatID is where to send it, filename is the path to the file.
func NewFileUpload(chatID int, filename string, fileType FileType) FileConfig {
	return FileConfig{
		ChatID:           chatID,
		UseExistingPhoto: false,
		FilePath:         filename,
		FileType:         fileType,
	}
}

// NewFileShare shares an existing photo.
// You may use this to reshare an existing photo without reuploading it.
//
// chatID is where to send it, fileID is the ID of the file already uploaded.
func NewPhotoShare(chatID int, fileID string, fileType FileType) FileConfig {
	return FileConfig{
		ChatID:           chatID,
		UseExistingPhoto: true,
		FileID:           fileID,
		FileType:         fileType,
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
