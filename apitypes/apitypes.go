package apitypes

type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
