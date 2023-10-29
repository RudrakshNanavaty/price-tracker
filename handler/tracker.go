package handler

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"price-tracker/database"
	"price-tracker/entities"
)

func GetProductPrice(url string) (float64, error) {
	cmd := exec.Command("python3", "./scripts/get_price.py", url)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(string(output), 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func StartTracking(product *entities.Product, db *database.DB) {
	desiredTime := time.Date(1, time.January, 1, 12, 0, 0, 0, time.Local)

	for {
		if time.Now().Hour() == desiredTime.Hour() {

			newPrice, err := GetProductPrice(product.URL)
			if err != nil {
				return
			}

			oldPrice, err := db.GetPrice(product.ID)
			if err != nil {
				return
			}

			if newPrice != oldPrice {
				db.UpdatePrice(product.ID, newPrice)
				fmt.Printf("Price updated.\n%s\n%f -> %f\n", product.Name, oldPrice, newPrice)
			}

			fmt.Printf("Price updated.\n%s\n%f -> %f\n", product.Name, oldPrice, newPrice)

			// Sleep for a day before checking again.
			time.Sleep(24 * time.Hour)
		} else {
			// sleep for 15 minutes for rate limiting
			time.Sleep(15 * time.Minute)
		}
	}
}
