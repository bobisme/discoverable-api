package errors

type errorResource struct {
	Message string `json:"message"`
	Logref  int    `json:"logref,omitempty"`
	Path    string `json:"path,omitempty"`
}

type resp struct {
	Total    int `json:"total"`
	Embedded struct {
		Errors []errorResource `json:"errors"`
	} `json:"_embedded"`
}
