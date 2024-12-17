package domain

type Client struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	AppId         string `json:"appId"`
	PhoneNumberId string `json:"phoneNumberId"`
	BusinessId    string `json:"businessId"`
	AccessToken   string `json:"accessToken"`
}
