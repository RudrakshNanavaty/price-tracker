package handler

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"price-tracker/entities"
	"price-tracker/database"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetAll() (*[]entities.Product, error) {
	data, _ := h.db.GetAll()

	return data, nil
}

func (h *Handler) GetPrice(url string) (float64, error) {
	outputChannel := make(chan float64)
	errorChannel := make(chan error)

	go fetchPrice(url, outputChannel, errorChannel)

	select {
	case err := <-errorChannel:
		return 0, err
	case price := <-outputChannel:
		return price, nil
	}
}

func (h *Handler) TrackPrice(url string) (*entities.Product, error) {
	outputChannel := make(chan entities.Product)
	errorChannel := make(chan error)

	go fetchProductInfo(url, outputChannel, errorChannel)

	select {
	case err := <-errorChannel:
		fmt.Println("Error fetching price")
		return nil, err

	case product := <-outputChannel:
		h.db.AddProduct(&product)

		go StartTracking(&product, h.db)

		return &product, nil
	}
}

func fetchPrice(url string, outputChannel chan float64, errorChannel chan error) {
	command := exec.Command("python3", "./scripts/get_price.py", url)

	commandOutput, err := command.Output()
	if err != nil {
		fmt.Println("Error running the get_price.py\n", err)
		errorChannel <- err
	}

	price, err := strconv.ParseFloat(string(commandOutput), 64)
	if err != nil {
		fmt.Println("Error parsing the output:\n", err)
		errorChannel <- err
	}

	outputChannel <- price
}

func fetchProductInfo(url string, outputChannel chan entities.Product, errorChannel chan error) (entities.Product, error) {
	command := exec.Command("python3", "./scripts/get_info.py", url)

	output, err := command.Output()
	if err != nil {
		fmt.Println("Error running the get_info.py\n", err)
		errorChannel <- err
	}

	var product entities.Product
	if err := json.Unmarshal([]byte(output), &product); err != nil {
		return entities.Product{}, err
	}

	return product, nil
}
