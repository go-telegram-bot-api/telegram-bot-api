package tgbotapi

import (
	"reflect"
	"testing"
)

type TestCase struct {
	Name                string
	InHTML              string
	InStrict            bool
	ExpectedText        string
	ExpectedOutEntities []MessageEntity
	ExpectedErr         error
}

func TestHtmlToEntities(t *testing.T) {

	testCases := []TestCase{
		{
			Name:         "b-tag",
			InHTML:       `<b>bold text</b>`,
			InStrict:     false,
			ExpectedText: "bold text",
			ExpectedOutEntities: []MessageEntity{
				{
					Type:   "bold",
					Offset: 0,
					Length: 9,
					Tag:    "b",
				},
			},
			ExpectedErr: nil,
		},
		{
			Name:         "strong-and-em-tags",
			InHTML:       `This is <strong>strong</strong> and <em>italic</em> text.`,
			InStrict:     false,
			ExpectedText: "This is strong and italic text.",
			ExpectedOutEntities: []MessageEntity{
				{
					Type:   "bold",
					Offset: 8,
					Length: 6,
					Tag:    "strong",
				},
				{
					Type:   "italic",
					Offset: 19,
					Length: 6,
					Tag:    "em",
				},
			},
			ExpectedErr: nil,
		},
		{
			Name:         "a-tag",
			InHTML:       `Click <a href="https://www.w3.org/">here</a>!`,
			InStrict:     false,
			ExpectedText: "Click here!",
			ExpectedOutEntities: []MessageEntity{
				{
					Type:   "text_link",
					Offset: 6,
					Length: 4,
					URL:    "https://www.w3.org/",
					Tag:    "a",
				},
			},
			ExpectedErr: nil,
		},
		{
			Name:         "a-with-emoji",
			InHTML:       `👉 <a href="https://www.w3.org/">more</a>`,
			InStrict:     false,
			ExpectedText: "👉 more",
			ExpectedOutEntities: []MessageEntity{
				{
					Type:   "text_link",
					Offset: 3,
					Length: 4,
					URL:    "https://www.w3.org/",
					Tag:    "a",
				},
			},
			ExpectedErr: nil,
		},
		{
			Name:         "flag-emoji",
			InHTML:       `🇺🇸 <u>'merica</u>`,
			InStrict:     false,
			ExpectedText: "🇺🇸 'merica",
			ExpectedOutEntities: []MessageEntity{
				{
					Type:   "underline",
					Offset: 5,
					Length: 7,
					Tag:    "u",
				},
			},
			ExpectedErr: nil,
		},
		{
			Name:         "korean-char",
			InHTML:       `<b>내 호버크라프트는 장어로 가득 차있다</b>`,
			InStrict:     false,
			ExpectedText: "내 호버크라프트는 장어로 가득 차있다",
			ExpectedOutEntities: []MessageEntity{
				{
					Type:   "bold",
					Offset: 0,
					Length: 20,
					Tag:    "b",
				},
			},
			ExpectedErr: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			text, entities, err := HtmlToEntities(tt.InHTML, tt.InStrict)
			if err != tt.ExpectedErr {
				t.Errorf("Expected err \"%s\", got \"%s\"", tt.ExpectedErr, err)
				return
			}
			if text != tt.ExpectedText {
				t.Errorf("Expected text \"%s\", got \"%s\"", tt.ExpectedText, text)
				return
			}
			if len(entities) != len(tt.ExpectedOutEntities) {
				t.Errorf("Expected %d entities, got %d", len(tt.ExpectedOutEntities), len(entities))
				return
			}
			if !reflect.DeepEqual(entities, tt.ExpectedOutEntities) {
				t.Errorf("Expected entities %+v, got %+v", tt.ExpectedOutEntities, entities)
				return
			}
		})
	}
}
