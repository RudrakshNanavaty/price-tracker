package handler

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"price-tracker/database"
	"price-tracker/entities"
)

func getProductPrice(url string) (float64, error) {
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

func startTracking(product *entities.Product, db *database.DB) {
	for {
		scriptOutputChannel := make(chan float64)
		scriptErrorChannel := make(chan error)
		dbOutputChannel := make(chan float64)
		dbErrorChannel := make(chan error)

		go func(outputChannel chan float64, errorChannel chan error) {
			newPrice, err := getProductPrice(product.URL)
			if err != nil {
				errorChannel <- err
				return
			}
			outputChannel <- newPrice
		}(scriptOutputChannel, scriptErrorChannel)

		go func(outputChannel chan float64, errorChannel chan error) {
			oldPrice, err := db.GetPrice(product.ID)
			if err != nil {
				errorChannel <- err
				return
			}
			outputChannel <- oldPrice
		}(dbOutputChannel, dbErrorChannel)

		select {
		case err := <-scriptErrorChannel:
			fmt.Println("Error fetching price: ", err)
			return
		case err := <-dbErrorChannel:
			fmt.Println("Error fetching price: ", err)
			return
		case oldPrice := <-scriptOutputChannel:
			newPrice := <-dbOutputChannel

			if newPrice != oldPrice {
				db.UpdatePrice(product.ID, newPrice)
				fmt.Printf("Price updated.\n%s\n%f -> %f\n", product.Name, oldPrice, newPrice)
			}
		}

		// Sleep for a day before checking again.
		time.Sleep(24 * time.Hour)
	}
}
