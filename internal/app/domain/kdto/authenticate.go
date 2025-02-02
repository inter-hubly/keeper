package kdto

type Authenticate struct {
	AccessToken string `json:"accessToken"`
	TenantId    string `json:"tenantId"`
}
