package domain

const (
	ActionRegister    string = "register"
	ActionLeaveChat   string = "leave-chat"
	ActionJoinChat    string = "join-chat"
	ActionSandMessage string = "send-message"
	ActionSendPrivate string = "send-private"
)

type Base struct {
	Action string `json:"action"`
}

type Register struct {
	Base        // actionRegister
	Name string `json:"name"`
}

type JoinChatroom struct {
	Base       // actionJoinChat
	ChatID int `json:"chat_id"`
}

type LeaveChat struct {
	Base       // actionLeaveChat
	ChatID int `json:"chat_id"`
}

type SendMessageToAll struct {
	Base           // actionSandMessage
	Message string `json:"message"`
}

type SendMessageToOne struct {
	Base           // actionSendPrivate
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}
