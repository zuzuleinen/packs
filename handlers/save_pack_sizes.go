package handlers

import (
	"encoding/json"
	"net/http"

	"gymshark/db"
)

type savePackSizesRequest struct {
	PackSizes []int `json:"packSizes"`
}

func SavePackSizesHandler(packSizeRepo *db.PackSizeRepository) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req savePackSizesRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Failed to parse request body", http.StatusBadRequest)
				return
			}

			packSizeRepo.DeleteAll()

			for _, size := range req.PackSizes {
				err := packSizeRepo.CreateSize(size)
				if err != nil {
					http.Error(w, "Failed to save sizes", http.StatusInternalServerError)
					return
				}
			}

			jsonSuccess(w, "Created", http.StatusCreated)
		},
	)
}
