package handlers

import (
	"net/http"

	"gymshark/db"
)

type getPackSizesResponse struct {
	PackSizes []int `json:"packSizes"`
}

func GetPackSizes(packSizeRepo *db.PackSizeRepository) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ps := packSizeRepo.FindAll()

			var resp getPackSizesResponse
			resp.PackSizes = make([]int, 0)
			for _, v := range ps {
				resp.PackSizes = append(resp.PackSizes, v.Size)
			}

			jsonResponse(w, resp)
		},
	)
}
