package handlers

type RequestBody struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type ResponseBody struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
	Refresh     string `json:"refresh_token"`
	Scope       string `json:"scope"`
}

type DiskPayload struct {
	Files []Files `json:"files"`
}

type Files struct {
	Name string `json:"name"`
}
