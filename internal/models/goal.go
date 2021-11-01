package models

type GoalEventSubEvent struct {
	ID                   string  `json:"id"`
	BroadcasterUserID    string  `json:"broadcaster_user_id"`
	BroadcasterUserName  string  `json:"broadcaster_user_name"`
	BroadcasterUserLogin string  `json:"broadcaster_user_login"`
	Type                 string  `json:"type"`
	Description          string  `json:"description"`
	IsAchieved           *bool   `json:"is_achieved,omitempty"`
	CurrentAmount        int64   `json:"current_amount"`
	TargetAmount         int64   `json:"target_amount"`
	StartedAt            string  `json:"started_at"`
	EndedAt              *string `json:"ended_at,omitempty"`
}

type GoalEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        GoalEventSubEvent    `json:"event"`
}
