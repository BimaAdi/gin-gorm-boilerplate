package schemas

type OauthLoginJsonRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	ResponseType string `json:"response_type" binding:"required"`
	ClientId     string `json:"client_id" binding:"required"`
	RedirectUri  string `json:"redirect_uri" binding:"required"`
	Scope        string `json:"scope" binding:"required"`
	State        string `json:"state" binding:"required"`
}

type Oauth2TokenJsonRequest struct {
	GrantType    string `json:"grant_type" binding:"required"`
	ClientId     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	RedirectUri  string `json:"redirect_uri" binding:"required"`
	Code         string `json:"code" binding:"required"`
}
