package domain

type Client struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	AppId         string `json:"appId"`
	PhoneNumberId string `json:"phoneNumberId"`
	BusinessId    string `json:"businessId"`
	AccessToken   string `json:"accessToken"`
}
