package integration_tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"integration-tests/client"
	"integration-tests/utils"
	"strings"
	"testing"
	"time"
)

func cleanupBeforeTest(t *testing.T) {
	//enabled, err := client.IsTargetTestsEnabled()
	//assert.Nil(t, enabled)
	//assert.
	err := client.ResetDB()
	if err != nil {
		panic(err)
	}
	client.ClearAps()
	client.GetAnswers()
}

func TestKek(t *testing.T) {
	ownerChatID := utils.GetOwnerChatID()
	if ownerChatID == 0 {
		panic("Owner chat id is not set. Provide it from env.")
	}

	t.Run("Test mode enabled", func(t *testing.T) {
		enabled, err := client.IsTargetTestsEnabled()
		if err != nil {
			t.Fatalf("Error checking test mode. " + err.Error())
		}
		if !enabled {
			t.Fatalf("Test mode disabled!")
		}
	})
	t.Run("For admin", func(t *testing.T) {
		t.Run("unknown for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("IejijiEIAjjIGi")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, answers[0].Text, `Не понял команду. Попробуйте /help, чтобы получить список команд`)
		})
		t.Run("/start for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("/start")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.True(t, strings.HasPrefix(answers[0].Text, `Привет! Это бот`))
		})
		t.Run("/help for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("/help")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.True(t, strings.Contains(answers[0].Text, `/grant - выдать доступ`))
			assert.True(t, strings.Contains(answers[0].Text, `Вы можете писать /stop или /start`))
			assert.True(t, strings.Contains(answers[0].Text, `Инструкция - /filterHelp`))
		})
		t.Run("/filterHelp for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("/filterHelp")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.True(t, strings.Contains(answers[0].Text, `1. Зайти на https://krisha.kz/map/arenda/kvartiry/almaty/`))
			assert.True(t, strings.Contains(answers[0].Text, `2. Выбрать нужные фильтры в панели над картой`))
			assert.True(t, strings.Contains(answers[0].Text, `3. ВАЖНО - нажать синюю кнопку`))
			assert.True(t, strings.Contains(answers[0].Text, `/stop - остановить уведомления`))
		})
		t.Run("Enable parser for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			status, err := client.CreateNAps(14)
			assert.Nil(t, err)
			assert.Equal(t, status, 201)
			err = client.SendUpdateFromOwner("https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)

			assert.Equal(t, len(answers), 4)
			assert.Equal(t, answers[0].Text, `Фильтр успешно установлен и парсер запущен`)
			assert.Equal(t, answers[1].Text, `Квартир: 14`)
			assert.Equal(t, answers[2].Text, `Начинаю собирать существующие квартиры, это займет немного времени...`)
			assert.Equal(t, answers[3].Text, `Существующие квартиры собраны, начинаю присылать уведомления о новых...`)
		})
		t.Run("Disable enabled parser for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			status, err := client.CreateNAps(14)
			assert.Nil(t, err)
			assert.Equal(t, status, 201)
			err = client.SendUpdateFromOwner("https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, len(answers), 4) //Values checked in prev test

			err = client.SendUpdateFromOwner("/stop")
			assert.Nil(t, err)
			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, `Парсер остановлен`)
		})
		t.Run("Notifying about one new ap", func(t *testing.T) {
			cleanupBeforeTest(t)
			status, err := client.CreateNAps(14)
			assert.Nil(t, err)
			assert.Equal(t, status, 201)
			err = client.SendUpdateFromOwner("https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, len(answers), 4) //Values checked in prev test

			err = client.SendUpdateFromOwner("/stop")
			assert.Nil(t, err)
			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, `Парсер остановлен`)
		})
	})

	t.Run("For user 1 without rights", func(t *testing.T) {
		t.Run("/start without rights", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdate(1, "/start")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Empty(t, answers)
		})
	})
	fmt.Println()
}
