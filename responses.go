package main

import (
	"encoding/json"
)

type ApiResponse struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type GroupChat struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type UserOrGroupChat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Title     string `json:"title"`
}

type Message struct {
	MessageId           int             `json:"message_id"`
	From                User            `json:"from"`
	Date                int             `json:"date"`
	Chat                UserOrGroupChat `json:"chat"`
	ForwardFrom         User            `json:"forward_from"`
	ForwardDate         int             `json:"forward_date"`
	ReplyToMessage      *Message        `json:"reply_to_message"`
	Text                string          `json:"text"`
	Audio               Audio           `json:"audio"`
	Document            Document        `json:"document"`
	Photo               []PhotoSize     `json:"photo"`
	Sticker             Sticker         `json:"sticker"`
	Video               Video           `json:"video"`
	Contact             Contact         `json:"contact"`
	Location            Location        `json:"location"`
	NewChatParticipant  User            `json:"new_chat_participant"`
	LeftChatParticipant User            `json:"left_chat_participant"`
	NewChatTitle        string          `json:"new_chat_title"`
	NewChatPhoto        string          `json:"new_chat_photo"`
	DeleteChatPhoto     bool            `json:"delete_chat_photo"`
	GroupChatCreated    bool            `json:"group_chat_created"`
}

type PhotoSize struct {
	FileId   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type Audio struct {
	FileId   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

type Document struct {
	FileId   string    `json:"file_id"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

type Sticker struct {
	FileId   string    `json:"file_id"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Thumb    PhotoSize `json:"thumb"`
	FileSize int       `json:"file_size"`
}

type Video struct {
	FileId   string    `json:"file_id"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Duration int       `json:"duration"`
	Thumb    PhotoSize `json:"thumb"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
	Caption  string    `json:"caption"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserId      string `json:"user_id"`
}

type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type UserProfilePhotos struct {
	TotalCount int         `json:"total_count"`
	Photos     []PhotoSize `json:"photos"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        map[string]map[string]string `json:"keyboard"`
	ResizeKeyboard  bool                         `json:"resize_keyboard"`
	OneTimeKeyboard bool                         `json:"one_time_keyboard"`
	Selective       bool                         `json:"selective"`
}

type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"`
}

type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"force_reply"`
}
