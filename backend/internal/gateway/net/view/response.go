package view

type ErrorResponse struct {
	Error string `json:"error" example:"error"`
}

type MessageResponse struct {
	Message string `json:"message" example:"some info"`
}

type IDResponse struct {
	ID string `json:"id" example:"asidofadsklhf"`
}

type FileUploadResponse struct {
	Name string `json:"name" example:"example.jpg"`
}
