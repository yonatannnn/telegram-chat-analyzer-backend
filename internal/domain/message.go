// internal/domain/chat.go
package domain

type Chat struct {
	Name     string    `json:"name" bson:"name"`
	Type     string    `json:"type" bson:"type"`
	ID       int       `json:"id" bson:"id"`
	Messages []Message `json:"messages" bson:"messages"`
}

type Message struct {
	ID               int          `json:"id" bson:"id"`
	Type             string       `json:"type" bson:"type"`
	Date             string       `json:"date" bson:"date"`
	DateUnixtime     string       `json:"date_unixtime" bson:"date_unixtime"`
	Edited           string       `json:"edited,omitempty" bson:"edited,omitempty"`
	EditedUnixtime   string       `json:"edited_unixtime,omitempty" bson:"edited_unixtime,omitempty"`
	From             string       `json:"from" bson:"from"`
	FromID           string       `json:"from_id" bson:"from_id"`
	Text             interface{}  `json:"text" bson:"text"`
	ReplyToMessageID int          `json:"reply_to_message_id,omitempty" bson:"reply_to_message_id,omitempty"`
	TextEntities     []TextEntity `json:"text_entities" bson:"text_entities"`
}

type TextEntity struct {
	Type string `json:"type" bson:"type"`
	Text string `json:"text" bson:"text"`
}
