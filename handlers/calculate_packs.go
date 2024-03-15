package handlers

import (
	"net/http"
	"strconv"

	"gymshark/db"
	"gymshark/domain"
)

type pack struct {
	Qty          int `json:"qty"`
	ItemsPerUnit int `json:"items_per_unit"`
}

type calculatePacksResponse struct {
	Packs []pack `json:"packs"`
}

func CalculatePacks(packSizeRepo *db.PackSizeRepository) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			itemsNo, err := strconv.Atoi(r.PathValue("itemsNo"))
			if err != nil {
				jsonError(w, "could not process itemsNo", http.StatusBadRequest)
				return
			}

			dbSizes := packSizeRepo.FindAll()

			cfg := domain.NewConfig()
			for _, v := range dbSizes {
				cfg.AddPackSize(v.Size)
			}

			calculator := domain.NewCalculator(cfg)

			calculatedPacks := calculator.Packs(itemsNo)

			var resp calculatePacksResponse
			resp.Packs = make([]pack, 0)

			for _, v := range calculatedPacks {
				resp.Packs = append(resp.Packs, pack{
					Qty:          v.Qty,
					ItemsPerUnit: v.ItemsPerUnit,
				})
			}

			jsonResponse(w, resp)
		},
	)
}
