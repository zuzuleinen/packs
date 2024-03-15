package handlers

import (
	"log"
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

func CalculatePacks(packSizeRepo *db.PackSizeRepository, logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			itemsNo, err := strconv.Atoi(r.PathValue("itemsNo"))
			if err != nil {
				logger.Println("error converting `itemsNo` to integer")
				jsonError(w, "could not process itemsNo", http.StatusBadRequest)
				return
			}

			cfg := domain.NewConfig()
			for _, v := range packSizeRepo.FindAll() {
				cfg.AddPackSize(v.Size)
			}

			calc := domain.NewCalculator(cfg)
			calculatedPacks := calc.Packs(itemsNo)

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
