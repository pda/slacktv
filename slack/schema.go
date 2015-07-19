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

type RTMStartResponse struct {
	Ok    bool
	Error string
	Url   string
	Self  *SelfResource
	Team  *TeamResource
	Users []UserResource
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
