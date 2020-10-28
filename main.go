package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//DataFileName name of users dat file
var DataFileName = "users.dat"
var lastCheckDate = time.Now().AddDate(0, 0, -21)

func main() {
	telegramToken := getTokenFromTile("telegram_token", "Telegram token (telegram_token)")
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	go waitNewUsers(bot)
	go checkAndSendReleases(bot)
	select {}
}

//GetUsers gets all users from dat file
func GetUsers() []int64 {
	log.Printf("Try read data from file %v", DataFileName)
	file, err := os.Open(DataFileName)

	result := make([]int64, 0)
	if err != nil && !os.IsExist(err) {
		log.Printf("Dat file not found, create them")
		os.Create(DataFileName)
		return result
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error read users files: %v\n", err)
			break
		}

		in, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		log.Printf("Get line from file: %d. Source line:%v\n", in, line)
		if err != nil {
			log.Printf("Data file corrupted: %v\n", err)
		} else {
			result = append(result, in)
		}
	}

	return result
}

//AddUser add user to file
func AddUser(user int64) {
	log.Printf("Add user %d to file %v", user, DataFileName)
	file, err := os.Open(DataFileName)
	if !os.IsExist(err) {
		log.Printf("Data file not found, create them")
		file, err = os.Create(DataFileName)
		if err != nil {
			panic(err)
		}
	}
	_, err = file.WriteString(fmt.Sprintf("%v\n", user))
	if err != nil {
		panic(err)
	}
	file.Close()
}

func checkAndSendReleases(bot *tgbotapi.BotAPI) {
	for {
		githubToken := getTokenFromTile("github_token", "GitHub Token (github_token)")
		releases := getReleasesAfterDate(lastCheckDate, githubToken)
		connectedUsers := GetUsers()
		for _, release := range releases {
			fmt.Printf("%v: %v at %v\n", release.repName, release.version, release.releaseAt.Format("01-02-2006"))
			log.Printf("Connected %v users", len(connectedUsers))
			for _, user := range connectedUsers {
				msg := tgbotapi.NewMessage(user, fmt.Sprintf("%v: %v at %v: %v\n", release.repName, release.version, release.releaseAt.Format("01-02-2006"), release.url))
				log.Printf("Send message to user %v", user)
				bot.Send(msg)
			}
		}
		lastCheckDate = time.Now()
		time.Sleep(10 * 60 * 1000 * time.Millisecond)
	}
}

func waitNewUsers(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Printf("Wait messages error: %s", err)
		return
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		AddUser(update.Message.Chat.ID)
		bot.Send(msg)
	}
}

func test(t time.Time) {
	fmt.Printf(t.String())
}

func getTokenFromTile(fileName string, tokenKind string) string {
	tokenFile, err := os.Open(fileName)
	if os.IsNotExist(err) {
		log.Fatalf("%s token file not found", tokenKind)
		panic(err)
	}

	scanner := bufio.NewScanner(tokenFile)
	scanner.Scan()
	token := scanner.Text()
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return token
}
