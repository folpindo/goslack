package apitypes

type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ChannelInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsMember bool   `json:"is_member"`
}

type Self struct {
	ID string `json:"id"`
}

type RtmStart struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	URL   string `json:"url"`
	Self  Self   `json:"self"`
}

type UserList struct {
	Ok      bool       `json:"ok"`
	Error   string     `json:"error"`
	Members []UserInfo `json:"members"`
}

type ChannelList struct {
	Ok       bool          `json:"ok"`
	Error    string        `json:"error"`
	Channels []ChannelInfo `json:"channels"`
}
