package src

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, httpCode int, success bool, message string, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	m := map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}
	_ = json.NewEncoder(w).Encode(m)
}

func ResizeImage(response http.ResponseWriter, request *http.Request) {
	var objReq ResizeImagesReq
	err := json.NewDecoder(request.Body).Decode(&objReq)
	if err != nil {
		Response(response, http.StatusBadRequest, false, GetMsg(INVALID_PARAMS), nil)
		return
	}
	resObj := objReq.ResizeImages()
	Response(response, http.StatusOK, true, GetMsg(SUCCESS), resObj)
}
