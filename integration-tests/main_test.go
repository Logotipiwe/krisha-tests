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

const (
	sellHousePath = "/map/prodazha/doma-dachi/almaty/"
	newApText     = "Новое объявление"
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
	client.SetAutoGrantLimit("")
	time.Sleep(200 * time.Millisecond)
}

func TestParser(t *testing.T) {
	ownerChatID := utils.GetOwnerChatID()
	if ownerChatID == 0 {
		panic("Owner chat id is not set. Provide it from env.")
	}
	client.SetAutoGrantLimit("")

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
		//t.Skip()
		t.Run("unknown for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("IejijiEIAjjIGi")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, `Не понял команду. Попробуйте /help, чтобы получить список команд`, answers[0].Text)
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
			assert.True(t, strings.Contains(answers[0].Text, `/info - сводка`))
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
		t.Run("/chats for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("1")
			client.SendUpdate(1, "/start")
			client.SendUpdate(10, "/start")
			client.SendUpdate(100, "/start")
			client.SendUpdate(1000, "/start")
			sleepForNotification() //TODO speed up just db log
			client.SendUpdateFromOwner("/chats")
			answers, err := client.GetAnswersToOwnerChat()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.True(t, strings.Contains(answers[0].Text, `Известные чаты`))
			assert.True(t, strings.Contains(answers[0].Text, `Title 1 (1)`))
			assert.True(t, strings.Contains(answers[0].Text, `Title 10 (10)`))
			assert.True(t, strings.Contains(answers[0].Text, `Title 100 (100)`))
			assert.True(t, strings.Contains(answers[0].Text, `Title 1000 (1000)`))
		})

		t.Run("Enable parser for admin", func(t *testing.T) {
			cleanupBeforeTest(t)
			status, err := client.CreateNAps(14)
			assert.Nil(t, err)
			assert.Equal(t, status, 201)
			err = client.SendUpdateFromOwner("https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)

			assert.Equal(t, 4, len(answers))
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
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, 4, len(answers)) //Values checked in prev test

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
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, len(answers), 4) //Values checked in prev test

			err = client.AddAp(15, "title 15", 1500, []string{"1", "11"})
			assert.Nil(t, err)
			sleepForNotification()

			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, newApText+`: https://krisha.kz/a/show/15`)
			assert.Equal(t, answers[0].Images, []string{"1", "11"})
		})
		t.Run("Notifying about selling house", func(t *testing.T) {
			cleanupBeforeTest(t)
			status, err := client.CreateNApsByPath(14, sellHousePath)
			assert.Nil(t, err)
			assert.Equal(t, status, 201)
			err = client.SendUpdateFromOwner("https://krisha.kz/map/prodazha/doma-dachi/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, len(answers), 4)

			err = client.AddApByPath(15, "title 15", 1500, []string{"1", "11"}, sellHousePath)
			assert.Nil(t, err)
			sleepForNotification()

			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, newApText+`: https://krisha.kz/a/show/15`)
			assert.Equal(t, answers[0].Images, []string{"1", "11"})
		})
		t.Run("Not notifying about new ap in other path", func(t *testing.T) {
			cleanupBeforeTest(t)
			status, err := client.CreateNApsByPath(14, sellHousePath)
			assert.Nil(t, err)
			assert.Equal(t, status, 201)
			err = client.SendUpdateFromOwner("https://krisha.kz/map/prodazha/doma-dachi/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, len(answers), 4)

			err = client.AddApByPath(15, "title 15", 1500, []string{"1", "11"}, sellHousePath+"kek/")
			assert.Nil(t, err)
			sleepForNotification()

			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Empty(t, answers)
		})
		t.Run("Notifying about exactly one new ap when several added in different paths", func(t *testing.T) {
			cleanupBeforeTest(t)
			status, err := client.CreateNApsByPath(14, sellHousePath)
			assert.Nil(t, err)
			assert.Equal(t, status, 201)
			err = client.SendUpdateFromOwner("https://krisha.kz/map/prodazha/doma-dachi/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.Equal(t, len(answers), 4)

			err = client.AddApByPath(16, "title 16", 1600, []string{"1", "11"}, sellHousePath+"kek/")
			err = client.AddApByPath(15, "title 15", 1500, []string{"1", "11"}, sellHousePath)
			assert.Nil(t, err)
			sleepForNotification()

			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, newApText+`: https://krisha.kz/a/show/15`)
			assert.Equal(t, answers[0].Images, []string{"1", "11"})
		})
		t.Run("Granting rights with incorrect input", func(t *testing.T) {
			cleanupBeforeTest(t)
			err := client.SendUpdateFromOwner("/grant")
			assert.Nil(t, err)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, "Какому чату выдать доступ? И через пробел - лимит", answers[0].Text)

			err = client.SendUpdateFromOwner("someId")
			assert.Nil(t, err)
			answers, err = client.GetAnswers()
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, "Ошибка expected 2 args instead of 1", answers[0].Text)
			assert.Nil(t, err)
		})
		t.Run("Granting rights for chat '1'", func(t *testing.T) {
			cleanupBeforeTest(t)
			err := client.SendUpdateFromOwner("/grant")
			assert.Nil(t, err)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, "Какому чату выдать доступ? И через пробел - лимит")

			err = client.SendUpdateFromOwner("1 100")
			assert.Nil(t, err)
			answers, err = client.GetAnswers()
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, "Доступ выдан для чата 1 с лимитом 100")
			assert.Nil(t, err)
		})
		t.Run("Denying for chat '1'", func(t *testing.T) {
			cleanupBeforeTest(t)

			err := client.SendUpdateFromOwner("/grant")
			assert.Nil(t, err)
			err = client.SendUpdateFromOwner("1 100")
			assert.Nil(t, err)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)

			err = client.SendUpdateFromOwner("/deny")
			assert.Nil(t, err)
			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, "Какому чату запретить доступ?", answers[0].Text)

			err = client.SendUpdateFromOwner("1")
			assert.Nil(t, err)
			answers, err = client.GetAnswersToOwnerChat()
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, "Доступ запрещен чату 1", answers[0].Text)
			assert.Nil(t, err)
		})
		t.Run("Denying for some other chat", func(t *testing.T) {
			cleanupBeforeTest(t)
			err := client.SendUpdateFromOwner("/deny")
			assert.Nil(t, err)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, "Какому чату запретить доступ?")

			err = client.SendUpdateFromOwner("111111111")
			assert.Nil(t, err)
			answers, err = client.GetAnswers()
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, "У этого чата и так нет доступа. Спасибо", answers[0].Text)
			assert.Nil(t, err)
		})
		t.Run("Admin dashboard empty", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("54")
			client.SendUpdateFromOwner("/info")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.True(t, strings.Contains(answers[0].Text, "Стандартный интервал: 1сек"))
			assert.True(t, strings.Contains(answers[0].Text, "Авто лимит: 54"))
			assert.True(t, strings.Contains(answers[0].Text, "Нет активных парсеров"))
		})
		t.Run("Admin dashboard with some chats", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("54")
			client.CreateNAps(25)
			client.SendUpdate(10, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			client.SendUpdate(11, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			client.SendUpdate(12, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			sleepForNotification()
			client.SendUpdateFromOwner("/info")
			answers, err := client.GetAnswersToOwnerChat()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.True(t, strings.Contains(answers[0].Text, "Стандартный интервал: 1сек"))
			assert.True(t, strings.Contains(answers[0].Text, "Авто лимит: 54"))
			assert.True(t, strings.Contains(answers[0].Text, "10 - interval: 1, aps: 25, explicit: false"))
			assert.True(t, strings.Contains(answers[0].Text, "11 - interval: 1, aps: 25, explicit: false"))
			assert.True(t, strings.Contains(answers[0].Text, "12 - interval: 1, aps: 25, explicit: false"))
		})
		t.Run("Curr aps count displayed on dashboard correctly", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("54")
			client.CreateNAps(25)
			client.SendUpdate(10, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			sleepForNotification()
			client.SendUpdateFromOwner("/info")
			answers, err := client.GetAnswersToOwnerChat()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.True(t, strings.Contains(answers[0].Text, "10 - interval: 1, aps: 25, explicit: false"))
			client.AddAp(26, "26", 26, []string{"26"})
			client.AddAp(27, "27", 27, []string{"27"})
			sleepForNotification()
			client.SendUpdateFromOwner("/info")
			answers, err = client.GetAnswersToOwnerChat()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.True(t, strings.Contains(answers[0].Text, "10 - interval: 1, aps: 27, explicit: false"))
		})
		t.Run("Admin dashboard only active chats", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("54")
			client.CreateNAps(1)
			client.SendUpdate(10, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			client.SendUpdate(11, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			client.SendUpdate(12, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			client.SendUpdate(13, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			client.SendUpdateFromOwner("/deny")
			client.SendUpdateFromOwner("11")
			client.SendUpdateFromOwner("/deny")
			client.SendUpdateFromOwner("12")
			client.GetAnswersToOwnerChat()
			client.SendUpdateFromOwner("/info")
			sleepForNotification()
			answers, err := client.GetAnswersToOwnerChat()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.True(t, strings.Contains(answers[0].Text, "Стандартный интервал: 1сек"))
			assert.True(t, strings.Contains(answers[0].Text, "Авто лимит: 54"))
			assert.True(t, strings.Contains(answers[0].Text, "10 - interval"))
			assert.False(t, strings.Contains(answers[0].Text, "11 - interval"))
			assert.False(t, strings.Contains(answers[0].Text, "12 - interval"))
			assert.True(t, strings.Contains(answers[0].Text, "13 - interval"))
		})
		t.Run("Deny for user with auto grant", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.CreateNAps(15)
			client.SetAutoGrantLimit("20")
			client.SendUpdateFromOwner("/deny")
			client.SendUpdateFromOwner("7")
			client.SendUpdate(7, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			answers, err := client.GetAnswersToChat(7)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "У вас нет доступа к боту. Обратитесь к администратору", answers[0].Text)
		})
		t.Run("User with auto grant - has access after deny cancel", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.CreateNAps(15)
			client.SetAutoGrantLimit("20")
			client.SendUpdateFromOwner("/deny")
			client.SendUpdateFromOwner("7")
			client.SendUpdate(7, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			_, err := client.GetAnswersToChat(7)
			assert.Nil(t, err)
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("7 30")
			client.GetAnswers()
			client.SendUpdate(7, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			answers, err := client.GetAnswersToChat(7)
			assert.Nil(t, err)
			assert.Equal(t, "Фильтр успешно установлен и парсер запущен", answers[0].Text)
		})
	})

	t.Run("For users", func(t *testing.T) {
		//t.Skip()
		t.Run("/start for user 1 without rights", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdate(1, "/start")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "У вас нет доступа к боту. Обратитесь к администратору", answers[0].Text)
		})
		t.Run("/start for user 1 with rights", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 200")
			client.GetAnswers()

			client.SendUpdate(1, "/start")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.True(t, strings.Contains(answers[0].Text, `Привет! Это бот для получения уведомлений о новых квартирах по фильтрам`))
			assert.True(t, strings.Contains(answers[0].Text, `/help - общая информация и инструкция`))
		})
		t.Run("/help for user 1 with rights", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 200")
			client.GetAnswers()

			client.SendUpdate(1, "/help")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.True(t, strings.Contains(answers[0].Text, `Вы можете писать /stop или /start`))
			assert.True(t, strings.Contains(answers[0].Text, `Инструкция - /filterHelp`))
		})
		t.Run("/filterHelp for user 1 with rights", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 200")
			client.GetAnswers()

			client.SendUpdate(1, "/filterHelp")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)
			assert.True(t, strings.Contains(answers[0].Text, `1. Зайти на https://krisha.kz/map/arenda/kvartiry/almaty/`))
			assert.True(t, strings.Contains(answers[0].Text, `2. Выбрать нужные фильтры в панели над картой`))
			assert.True(t, strings.Contains(answers[0].Text, `3. ВАЖНО - нажать синюю кнопку`))
			assert.True(t, strings.Contains(answers[0].Text, `/stop - остановить уведомления`))
		})
		t.Run("Enable parser for user 1 with rights", func(t *testing.T) {
			cleanupBeforeTest(t)

			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 200")
			client.GetAnswers()
			client.CreateNAps(190)

			err := client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)

			assert.Equal(t, 4, len(answers))
			assert.Equal(t, answers[0].Text, `Фильтр успешно установлен и парсер запущен`)
			assert.Equal(t, answers[1].Text, `Квартир: 190`)
			assert.Equal(t, answers[2].Text, `Начинаю собирать существующие квартиры, это займет немного времени...`)
			assert.Equal(t, answers[3].Text, `Существующие квартиры собраны, начинаю присылать уведомления о новых...`)
		})
		t.Run("Disable enabled parser for user 1 with rights", func(t *testing.T) {
			cleanupBeforeTest(t)

			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 200")
			client.GetAnswers()
			client.CreateNAps(190)

			err := client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Equal(t, answers[0].Text, `Фильтр успешно установлен и парсер запущен`)

			err = client.SendUpdate(1, "/stop")
			assert.Nil(t, err)
			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, `Парсер остановлен`)
		})
		t.Run("Notifying about one new ap for user 1 with rights", func(t *testing.T) {
			cleanupBeforeTest(t)

			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 200")
			client.GetAnswers()
			client.CreateNAps(190)

			err := client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Equal(t, answers[0].Text, `Фильтр успешно установлен и парсер запущен`)

			err = client.AddAp(191, "title 191", 19100, []string{"1", "2"})
			assert.Nil(t, err)
			sleepForNotification()
			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, answers[0].Text, newApText+`: https://krisha.kz/a/show/191`)
			assert.Equal(t, answers[0].Images, []string{"1", "2"})
		})
		t.Run("Exceed limit for user 1 with rights", func(t *testing.T) {
			cleanupBeforeTest(t)

			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 200")
			client.GetAnswers()
			client.CreateNAps(220)

			err := client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Equal(t, answers[0].Text, `Превышен лимит в 200 квартир в вашем фильтре. Попробуйте другой фильтр`)
		})
		t.Run("Stop parser for user when deny", func(t *testing.T) {
			cleanupBeforeTest(t)

			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 200")
			client.GetAnswers()
			client.CreateNAps(190)

			err := client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.NotEmpty(t, answers)

			err = client.SendUpdateFromOwner("/deny")
			assert.Nil(t, err)
			err = client.SendUpdateFromOwner("1")
			assert.Nil(t, err)
			answers, err = client.GetAnswersToChat(1)
			assert.NotEmpty(t, answers)
			assert.Equal(t, len(answers), 1)
			assert.Equal(t, "Парсер остановлен, обратитесь к администратору", answers[0].Text)

			//Check that new ap doesn't trigger notification
			err = client.AddAp(191, "title 191", 19100, []string{"1"})
			assert.Nil(t, err)
			sleepForNotification()
			answers, err = client.GetAnswersToChat(1)
			assert.Nil(t, err)
			assert.Empty(t, answers)
		})
		t.Run("Allow same filter after grant limit", func(t *testing.T) {
			cleanupBeforeTest(t)

			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 201")
			client.GetAnswers()
			client.CreateNAps(220)

			err := client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, `Превышен лимит в 201 квартир в вашем фильтре. Попробуйте другой фильтр`, answers[0].Text)

			err = client.SendUpdateFromOwner("/grant")
			assert.Nil(t, err)
			err = client.SendUpdateFromOwner("1 230")
			assert.Nil(t, err)
			answers, err = client.GetAnswersToChat(1)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "Ваш лимит изменен на 230 квартир", answers[0].Text)
			err = client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			assert.Nil(t, err)
			sleepForNotification()
			answers, err = client.GetAnswersToChat(1)
			assert.Nil(t, err)
			assert.Equal(t, 4, len(answers))
			assert.Equal(t, "Фильтр успешно установлен и парсер запущен", answers[0].Text)
		})
	})
	t.Run("For users with auto grant", func(t *testing.T) {
		//t.Skip()
		t.Run("Enable parser from new user (exceed limit)", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("200")
			client.CreateNAps(220)
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "Превышен лимит в 200 квартир в вашем фильтре. Попробуйте другой фильтр", answers[0].Text)
		})
		t.Run("Enable parser from new user", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("200")
			client.CreateNAps(180)
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 4, len(answers))
			assert.Equal(t, "Фильтр успешно установлен и парсер запущен", answers[0].Text)
		})
		t.Run("Deny stops parser", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("30")
			client.CreateNAps(20)
			client.SendUpdate(1, "/start")
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			client.GetAnswers()
			client.SendUpdateFromOwner("/deny")
			client.SendUpdateFromOwner("1")
			sleepForNotification() //TODO speed up just db log
			sleepForNotification() //TODO speed up just db log
			sleepForNotification() //TODO speed up just db log
			answers, err := client.GetAnswersToChat(1)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "Парсер остановлен, обратитесь к администратору", answers[0].Text)
		})
		t.Run("Deny restricts parser restart", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("30")
			client.CreateNAps(20)
			client.SendUpdate(1, "/start")
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			client.GetAnswers()
			client.SendUpdateFromOwner("/deny")
			client.SendUpdateFromOwner("1")
			client.GetAnswers()
			client.SendUpdate(1, "/start")
			answers, err := client.GetAnswersToChat(1)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "У вас нет доступа к боту. Обратитесь к администратору", answers[0].Text)
		})
		t.Run("Allow 0 restricts parser restart", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("30")
			client.CreateNAps(20)
			client.SendUpdate(1, "/start")
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			client.GetAnswers()
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 0")
			client.GetAnswers()
			client.SendUpdate(1, "/start")
			answers, err := client.GetAnswersToChat(1)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "У вас нет доступа к боту. Обратитесь к администратору", answers[0].Text)
		})

		t.Run("Limit changes with auto limit parameter (when restart parser)", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("201")
			client.CreateNAps(180)
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 4, len(answers))
			assert.Equal(t, "Фильтр успешно установлен и парсер запущен", answers[0].Text)

			client.SendUpdate(1, "/stop")
			client.SetAutoGrantLimit("150")
			client.GetAnswers()
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			answers, err = client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "Превышен лимит в 150 квартир в вашем фильтре. Попробуйте другой фильтр", answers[0].Text)
		})
		t.Run("Explicit grant accounted instead of auto (deny) (auto first)", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("10")
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 5")
			client.CreateNAps(7)
			client.GetAnswers()

			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "Превышен лимит в 5 квартир в вашем фильтре. Попробуйте другой фильтр", answers[0].Text)
		})
		t.Run("Explicit grant accounted instead of auto (deny) (explicit first)", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 5")
			client.SetAutoGrantLimit("10")
			client.CreateNAps(7)
			client.GetAnswers()

			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, "Превышен лимит в 5 квартир в вашем фильтре. Попробуйте другой фильтр", answers[0].Text)
		})
		t.Run("Explicit grant accounted instead of auto (allow) (grant first)", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("10")
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 20")
			client.CreateNAps(15)
			client.GetAnswers()

			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 4, len(answers))
			assert.Equal(t, answers[0].Text, `Фильтр успешно установлен и парсер запущен`)
		})
		t.Run("Explicit grant accounted instead of auto (allow) (autogrant - /start - explicit grant)", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.SetAutoGrantLimit("10")
			client.SendUpdate(1, "/start")
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 20")
			client.CreateNAps(15)
			client.GetAnswers()

			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 4, len(answers))
			assert.Equal(t, answers[0].Text, `Фильтр успешно установлен и парсер запущен`)
		})

		t.Run("Explicit grant continue work when disable auto grant (auto first)", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.CreateNAps(4)
			client.SetAutoGrantLimit("20")
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 10")
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			client.GetAnswers()

			err := client.SetAutoGrantLimit("")
			assert.Nil(t, err)
			client.AddAp(5, "title 5", 500, []string{"1"})
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, newApText+": https://krisha.kz/a/show/5", answers[0].Text)
		})
		t.Run("Explicit grant continue work when disable auto grant (explicit first)", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.CreateNAps(4)
			client.SendUpdateFromOwner("/grant")
			client.SendUpdateFromOwner("1 10")
			client.SetAutoGrantLimit("20")
			client.SendUpdate(1, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			time.Sleep(1 * time.Second)
			client.GetAnswers()

			err := client.SetAutoGrantLimit("")
			assert.Nil(t, err)
			client.AddAp(5, "title 5", 500, []string{"1"})
			time.Sleep(1 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, 1, len(answers))
			assert.Equal(t, newApText+": https://krisha.kz/a/show/5", answers[0].Text)
		})
		t.Run("Auto stop when limit changed", func(t *testing.T) {
			cleanupBeforeTest(t)
			client.CreateNAps(10)
			client.SetAutoGrantLimit("20")
			client.SendUpdate(19, "https://krisha.kz/map/arenda/kvartiry/almaty/?test=params")
			sleepForNotification()
			client.GetAnswers()
			client.SetAutoGrantLimit("5")
			time.Sleep(3 * time.Second)
			answers, err := client.GetAnswers()
			assert.Nil(t, err)
			assert.Equal(t, "Парсер остановлен из-за изменения лимита на кол-во квартир. Попробуйте другой фильтр", answers[0].Text)
			client.AddAp(21, "title 21", 2100, []string{"21"})
			sleepForNotification()
			messages, err := client.GetAnswersToChat(19)
			assert.Empty(t, messages)
		})
	})
	fmt.Println()
}

func sleepForNotification() {
	time.Sleep(500 * time.Millisecond)
}
