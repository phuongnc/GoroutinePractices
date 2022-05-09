package src

type ResizeImagesReq struct {
	Data []ResizeImageItem `json:"data"`
}

type ResizeImageItem struct {
	ID      string   `json:"_id"`
	Bucket  string   `json:"bucket"`
	Prefix  string   `json:"prefix"`
	Key     string   `json:"key"`
	Type    string   `json:"type"`
	Profile []string `json:"profile"`
	Message []string `json:"message"`
}

type ResizeImagesRes struct {
	Data []ResizeImageItem `json:"data"`
}

type UploaderRes struct {
	fileName string
	err      error
}
