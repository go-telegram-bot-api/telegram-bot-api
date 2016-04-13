package tgbotapi

import (
	"net/url"
)

// NewMessage creates a new Message.
//
// chatID is where to send it, text is the message text.
func NewMessage(chatID int64, text string) MessageConfig {
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
func NewForward(chatID int64, fromChatID int64, messageID int) ForwardConfig {
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
//
// Note that you must send animated GIFs as a document.
func NewPhotoUpload(chatID int64, file interface{}) PhotoConfig {
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
func NewPhotoShare(chatID int64, fileID string) PhotoConfig {
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
func NewAudioUpload(chatID int64, file interface{}) AudioConfig {
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
func NewAudioShare(chatID int64, fileID string) AudioConfig {
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
func NewDocumentUpload(chatID int64, file interface{}) DocumentConfig {
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
func NewDocumentShare(chatID int64, fileID string) DocumentConfig {
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
func NewStickerUpload(chatID int64, file interface{}) StickerConfig {
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
func NewStickerShare(chatID int64, fileID string) StickerConfig {
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
func NewVideoUpload(chatID int64, file interface{}) VideoConfig {
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
func NewVideoShare(chatID int64, fileID string) VideoConfig {
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
func NewVoiceUpload(chatID int64, file interface{}) VoiceConfig {
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
func NewVoiceShare(chatID int64, fileID string) VoiceConfig {
	return VoiceConfig{
		BaseFile: BaseFile{
			BaseChat:    BaseChat{ChatID: chatID},
			FileID:      fileID,
			UseExisting: true,
		},
	}
}

// NewContact allows you to send a shared contact.
func NewContact(chatID int64, phoneNumber, firstName string) ContactConfig {
	return ContactConfig{
		BaseChat: BaseChat{
			ChatID: chatID,
		},
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
	}
}

// NewLocation shares your location.
//
// chatID is where to send it, latitude and longitude are coordinates.
func NewLocation(chatID int64, latitude float64, longitude float64) LocationConfig {
	return LocationConfig{
		BaseChat: BaseChat{
			ChatID: chatID,
		},
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewVenue allows you to send a venue and its location.
func NewVenue(chatID int64, title, address string, latitude, longitude float64) VenueConfig {
	return VenueConfig{
		BaseChat: BaseChat{
			ChatID: chatID,
		},
		Title:     title,
		Address:   address,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewChatAction sets a chat action.
// Actions last for 5 seconds, or until your next action.
//
// chatID is where to send it, action should be set via Chat constants.
func NewChatAction(chatID int64, action string) ChatActionConfig {
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

// NewInlineQueryResultArticle creates a new inline query article.
func NewInlineQueryResultArticle(id, title, messageText string) InlineQueryResultArticle {
	return InlineQueryResultArticle{
		Type:  "article",
		ID:    id,
		Title: title,
		InputMessageContent: InputTextMessageContent{
			Text: messageText,
		},
	}
}

// NewInlineQueryResultGIF creates a new inline query GIF.
func NewInlineQueryResultGIF(id, url string) InlineQueryResultGIF {
	return InlineQueryResultGIF{
		Type: "gif",
		ID:   id,
		URL:  url,
	}
}

// NewInlineQueryResultMPEG4GIF creates a new inline query MPEG4 GIF.
func NewInlineQueryResultMPEG4GIF(id, url string) InlineQueryResultMPEG4GIF {
	return InlineQueryResultMPEG4GIF{
		Type: "mpeg4_gif",
		ID:   id,
		URL:  url,
	}
}

// NewInlineQueryResultPhoto creates a new inline query photo.
func NewInlineQueryResultPhoto(id, url string) InlineQueryResultPhoto {
	return InlineQueryResultPhoto{
		Type: "photo",
		ID:   id,
		URL:  url,
	}
}

// NewInlineQueryResultVideo creates a new inline query video.
func NewInlineQueryResultVideo(id, url string) InlineQueryResultVideo {
	return InlineQueryResultVideo{
		Type: "video",
		ID:   id,
		URL:  url,
	}
}
