package tgbotapi

import (
	"net/url"
)

// NewMessage creates a new Message.
//
// chatID is where to send it, text is the message text.
func NewMessage(chatID int, text string) MessageConfig {
	return MessageConfig{
		BaseChat: BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text: text,
		DisableWebPagePreview: false,
	}
}

// NewForward creates a new forward.
//
// chatID is where to send it, fromChatID is the source chat,
// and messageID is the ID of the original message.
func NewForward(chatID int, fromChatID int, messageID int) ForwardConfig {
	return ForwardConfig{
		BaseChat:   BaseChat{ChatID: chatID},
		FromChatID: fromChatID,
		MessageID:  messageID,
	}
}

// NewPhotoUpload creates a new photo uploader.
//
// chatID is where to send it, file is a string path to the file,
// FileReader, or FileBytes.
func NewPhotoUpload(chatID int, file interface{}) PhotoConfig {
	return PhotoConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			File:        file,
			UseExisting: false,
		},
	}
}

// NewPhotoShare shares an existing photo.
// You may use this to reshare an existing photo without reuploading it.
//
// chatID is where to send it, fileID is the ID of the file
// already uploaded.
func NewPhotoShare(chatID int, fileID string) PhotoConfig {
	return PhotoConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			FileID:      fileID,
			UseExisting: true,
		},
	}
}

// NewAudioUpload creates a new audio uploader.
//
// chatID is where to send it, file is a string path to the file,
// FileReader, or FileBytes.
func NewAudioUpload(chatID int, file interface{}) AudioConfig {
	return AudioConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			File:        file,
			UseExisting: false,
		},
	}
}

// NewAudioShare shares an existing audio file.
// You may use this to reshare an existing audio file without
// reuploading it.
//
// chatID is where to send it, fileID is the ID of the audio
// already uploaded.
func NewAudioShare(chatID int, fileID string) AudioConfig {
	return AudioConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			FileID:      fileID,
			UseExisting: true,
		},
	}
}

// NewDocumentUpload creates a new document uploader.
//
// chatID is where to send it, file is a string path to the file,
// FileReader, or FileBytes.
func NewDocumentUpload(chatID int, file interface{}) DocumentConfig {
	return DocumentConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			File:        file,
			UseExisting: false,
		},
	}
}

// NewDocumentShare shares an existing document.
// You may use this to reshare an existing document without
// reuploading it.
//
// chatID is where to send it, fileID is the ID of the document
// already uploaded.
func NewDocumentShare(chatID int, fileID string) DocumentConfig {
	return DocumentConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			FileID:      fileID,
			UseExisting: true,
		},
	}
}

// NewStickerUpload creates a new sticker uploader.
//
// chatID is where to send it, file is a string path to the file,
// FileReader, or FileBytes.
func NewStickerUpload(chatID int, file interface{}) StickerConfig {
	return StickerConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			File:        file,
			UseExisting: false,
		},
	}
}

// NewStickerShare shares an existing sticker.
// You may use this to reshare an existing sticker without
// reuploading it.
//
// chatID is where to send it, fileID is the ID of the sticker
// already uploaded.
func NewStickerShare(chatID int, fileID string) StickerConfig {
	return StickerConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			FileID:      fileID,
			UseExisting: true,
		},
	}
}

// NewVideoUpload creates a new video uploader.
//
// chatID is where to send it, file is a string path to the file,
// FileReader, or FileBytes.
func NewVideoUpload(chatID int, file interface{}) VideoConfig {
	return VideoConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			File:        file,
			UseExisting: false,
		},
	}
}

// NewVideoShare shares an existing video.
// You may use this to reshare an existing video without reuploading it.
//
// chatID is where to send it, fileID is the ID of the video
// already uploaded.
func NewVideoShare(chatID int, fileID string) VideoConfig {
	return VideoConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			FileID:      fileID,
			UseExisting: true,
		},
	}
}

// NewVoiceUpload creates a new voice uploader.
//
// chatID is where to send it, file is a string path to the file,
// FileReader, or FileBytes.
func NewVoiceUpload(chatID int, file interface{}) VoiceConfig {
	return VoiceConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			File:        file,
			UseExisting: false,
		},
	}
}

// NewVoiceShare shares an existing voice.
// You may use this to reshare an existing voice without reuploading it.
//
// chatID is where to send it, fileID is the ID of the video
// already uploaded.
func NewVoiceShare(chatID int, fileID string) VoiceConfig {
	return VoiceConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			FileID:      fileID,
			UseExisting: true,
		},
	}
}

// NewLocation shares your location.
//
// chatID is where to send it, latitude and longitude are coordinates.
func NewLocation(chatID int, latitude float64, longitude float64) LocationConfig {
	return LocationConfig{
		BaseChat: BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
			ReplyMarkup:      nil,
		},
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewChatAction sets a chat action.
// Actions last for 5 seconds, or until your next action.
//
// chatID is where to send it, action should be set via Chat constants.
func NewChatAction(chatID int, action string) ChatActionConfig {
	return ChatActionConfig{
		BaseChat: BaseChat{ChatID: chatID},
		Action:   action,
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
		URL: u,
	}
}

// NewWebhookWithCert creates a new webhook with a certificate.
//
// link is the url you wish to get webhooks,
// file contains a string to a file, FileReader, or FileBytes.
func NewWebhookWithCert(link string, file interface{}) WebhookConfig {
	u, _ := url.Parse(link)

	return WebhookConfig{
		URL:         u,
		Certificate: file,
	}
}
