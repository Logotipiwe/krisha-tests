package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetOwnerChatID() int64 {
	chatIdStr := os.Getenv("OWNER_TG_CHAT_ID")
	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return chatId
}

func GetTargetHost() string {
	return os.Getenv("TARGET_HOST")
}

func GetApsMockHost() string {
	return os.Getenv("APS_MOCK_HOST")
}
