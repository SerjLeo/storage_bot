package telegram

type Update struct {
	Id      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text            string
	Chat            Chat             `json:"chat"`
	From            User             `json:"from"`
	Reply           *IncomingMessage `json:"reply_to_message"`
	Pinned          *IncomingMessage `json:"pinned_message"`
	Caption         string           `json:"caption"`
	CaptionEntities []MessageEntity  `json:"caption_entities"`
	Entities        []MessageEntity  `json:"entities"`
}

type MessageEntity struct {
	Type string `json:"type"`
	Url  string `json:"url"`
	User User   `json:"user"`
}

type User struct {
	Name string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}

type Response struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}
