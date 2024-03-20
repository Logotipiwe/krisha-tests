package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"integration-tests/model"
	"integration-tests/utils"
	"io"
	"net/http"
	"time"
)

func SendUpdate(chatID int64, text string) error {
	b := []byte(fmt.Sprintf(`{"ChatID":%v,"Text":"%v"}`, chatID, text))
	fmt.Printf("Sending `%v` from chat %v\n", text, chatID)
	request, _ := http.NewRequest("POST", utils.GetTargetHost()+"/tests/sendMessage", bytes.NewBuffer(b))
	client := http.Client{}
	//client.Timeout = 2 * time.Second
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return err
}

func SendUpdateFromOwner(text string) error {
	return SendUpdate(utils.GetOwnerChatID(), text)
}

func IsTargetTestsEnabled() (bool, error) {
	req, _ := http.NewRequest("GET", utils.GetTargetHost()+"/tests/enabled", nil)
	client := http.Client{}
	client.Timeout = 2 * time.Second
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	answerStr := string(body)
	if answerStr == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

func GetAnswers() ([]model.SentMockMessage, error) {
	req, _ := http.NewRequest("GET", utils.GetTargetHost()+"/tests/getAnswerMessages", nil)
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var answers []model.SentMockMessage
	body, err := io.ReadAll(response.Body)
	err = json.Unmarshal(body, &answers)
	fmt.Printf("Got messages %v\n", answers)
	return answers, err
}

func GetAnswersToOwnerChat() ([]model.SentMockMessage, error) {
	return GetAnswersToChat(utils.GetOwnerChatID())
}

func GetAnswersToChat(chatID int64) ([]model.SentMockMessage, error) {
	answers, err := GetAnswers()
	if err != nil {
		return answers, err
	}
	var result []model.SentMockMessage
	for _, m := range answers {
		if m.ChatID == chatID {
			result = append(result, m)
		}
	}
	return result, err
}

func ResetDB() error {
	fmt.Println("Performing db reset...")
	request, _ := http.NewRequest("POST", utils.GetTargetHost()+"/tests/reset", nil)
	client := http.Client{}
	client.Timeout = 5 * time.Second
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return err
}

func SetAutoGrantLimit(n string) error {
	fmt.Println("Setting AUTO_GRANT_LIMIT...")
	request, _ := http.NewRequest("POST", utils.GetTargetHost()+"/tests/set-auto-grant-limit?n="+n, nil)
	client := http.Client{}
	client.Timeout = 2 * time.Second
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return err
}

// SetEnv TODO do with saving in config server
func SetEnv(key, value string) error {
	fmt.Println("Setting env " + key + " to value " + value)
	request, _ := http.NewRequest("POST", utils.GetTargetHost()+"/tests/set-env?key="+key+"&value="+value, nil)
	client := http.Client{}
	client.Timeout = 2 * time.Second
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return err
}

func GetEnv(key string) (string, error) {
	fmt.Println("Getting env " + key)
	request, _ := http.NewRequest("GET", utils.GetTargetHost()+"/tests/get-env?key="+key, nil)
	client := http.Client{}
	client.Timeout = 2 * time.Second
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var answer []byte
	_, err = response.Body.Read(answer)
	return string(answer), err
}
