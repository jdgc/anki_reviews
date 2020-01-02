package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const Url = "http://localhost:8765"

type GetCardsResponse struct {
	Result []int  `json:"result"`
	Error  string `json:"error"`
}

type AreDueResponse struct {
	Result []bool `json:"result"`
	Error  string `json:"error"`
}

func main() {
	jsonBody := []byte(`{"action":"findCards","version":6,"params":{"query":"deck:current"}}`)

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(jsonBody))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var response GetCardsResponse

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(bodyBytes, &response)

	cards := response.Result

	due := dueCards(cards)
	fmt.Println("cards due:", due)
}

func dueCards(cards []int) int {
	cardsJson, err := json.Marshal(cards)
	if err != nil {
		panic(err)
	}

	jsonString := fmt.Sprintf(`{"action":"areDue","version":6,"params":{"cards":%s}}`, string(cardsJson))
	jsonBody := []byte(jsonString)
	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(jsonBody))
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var response AreDueResponse

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)

	var count int
	for i := 0; i < len(response.Result); i++ {
		if response.Result[i] {
			count++
		}
	}
	return count
}
