package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func addProduct(client *resty.Client, productId int) {

	resp, err := client.R().
		Post("http://localhost:10000/add/" + strconv.Itoa(productId))

	if err != nil {
		panic(err)
	}
	fmt.Sprintf("Add product %d response code %d", productId, resp.StatusCode())
}

func doTestScenario(client *resty.Client, username string, ch chan<-string) {

	body, _ := json.Marshal(map[string]string{
		"name":  username,
		"password": "1234",
	})

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("http://localhost:10000/login")

	if err != nil {
		panic(err)
	}
	fmt.Sprintf("login response code %d", resp.StatusCode())

	addProduct(client, 1)
	addProduct(client, 1)
	addProduct(client, 2)

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		Post("http://localhost:10000/checkout")

	if err != nil {
		panic(err)
	}
	fmt.Sprintf("Checkout response code %d", resp.StatusCode())
}

func parallelRun() {
	ch := make(chan string)
	client := resty.New()
	go doTestScenario(client, "Jon Doe", ch)
	go doTestScenario(client, "Jon Doe 2", ch)
	fmt.Println(<-ch)
}

func main() {
	client := resty.New()

	doTestScenario(client, "Jon Doe", nil)
	doTestScenario(client, "Jon Doe 2", nil)
}