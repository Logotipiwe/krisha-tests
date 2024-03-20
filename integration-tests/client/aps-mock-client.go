package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"integration-tests/utils"
	"net/http"
	"strconv"
	"time"
)

func AddAp(id int64, title string, price int, images []string) error {
	return addAp(id, title, price, images, "/map/arenda/kvartiry/almaty/")
}

func AddApByPath(id int64, title string, price int, images []string, subPath string) error {
	return addAp(id, title, price, images, subPath)
}

func CreateNAps(n int) (int, error) {
	return createNAps(n, "/map/arenda/kvartiry/almaty/")
}

func CreateNApsByPath(n int, subPath string) (int, error) {
	return createNAps(n, subPath)
}

func ClearAps() error {
	return clearAps()
}

func clearAps() error {
	fmt.Println("Clearing mock aps...\n")
	request, _ := http.NewRequest("POST", utils.GetApsMockHost()+"/clear-aps", nil)
	client := http.Client{}
	client.Timeout = 2 * time.Second
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func addAp(id int64, title string, price int, images []string, subPath string) error {
	fmt.Printf("Creating mock ap with title %v and price %v\n", title, price)
	body := make(map[string]any)
	body["id"] = id
	body["title"] = title
	body["price"] = price
	body["images"] = images

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)
	request, _ := http.NewRequest("POST", utils.GetApsMockHost()+"/create-ap?subPath="+subPath, buffer)

	client := http.Client{}
	client.Timeout = 2 * time.Second
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func createNAps(n int, subPath string) (int, error) {
	fmt.Printf("Creating %v aps\n", n)
	request, _ := http.NewRequest(
		"POST", utils.GetApsMockHost()+"/create-n-aps?n="+strconv.Itoa(n)+"&subPath="+subPath, nil)
	client := http.Client{}
	client.Timeout = 2 * time.Second
	response, err := client.Do(request)
	status := response.StatusCode
	if err != nil {
		return status, err
	}
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	return status, err
}
