package slack

type ChannelResource struct {
	Id   string
	Name string
}

type ChannelResponse struct {
	Response
	Channel ChannelResource
}

type Response struct {
	Ok bool
}

type RtmStartResponse struct {
	Ok    bool           `json:"ok"`
	Url   string         `json:"url"`
	Self  *SelfResource  `json:"self"`
	Team  *TeamResource  `json:"team"`
	Users []UserResource `json:"users"`
}

type SelfResource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TeamResource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	EmailDomain string `json:"email_domain"`
	Domain      string `json:"domain"`
}

type UserResource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
