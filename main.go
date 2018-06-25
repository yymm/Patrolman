package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nlopes/slack"
)

// Config : config.json内の構造体
type Config struct {
	Sites []struct {
		Selector string `json:"selector"`
		URL      string `json:"url"`
	} `json:"sites"`
	Slack struct {
		Token     string `json:"token"`
		Channelid string `json:"channelid"`
	} `json:"slack"`
	Intervalhour float32 `json:"intervalhour"`
}

func slackNotify(token, channelID, message string) {
	api := slack.New(token)
	params := slack.PostMessageParameters{
		AsUser: true, // ユーザーとしてpostする(招待されたページのみpost可能)
		//LinkNames: 1, // メンションを可能にする
		// 以下のオプションを使うと招待されていないチャンネルでも
		// 自由なIcon、ユーザー名でpostできる
		//IconURL:  "https://pbs.twimg.com/media/DUnrnDHVQAAiXff.jpg",
		//Username: "まつぼっくり",
	}
	_, _, err := api.PostMessage(channelID, message, params)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}

func loadConfig() (Config, error) {
	// Hack: structを用いたデフォルト引数を
	// 使って任意名のconfigファイルを読めるようにする
	var config Config
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, err
	}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return config, err
	}
	return config, nil
}

// func getHTMLSize(url string) (int, error) {
// 	res, err := http.Get(url)
// 	if err != nil {
// 		log.Printf("Error(http.Get): %v\n", err)
// 		return 0, err
// 	}
// 	defer res.Body.Close()
// 	bytes, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Printf(" Error(Read HTML): %v\n", err)
// 		return 0, err
// 	}
// 	return len(bytes), nil
// }

func mainLoop(config *Config) {
	t := time.Duration(config.Intervalhour * float32(time.Hour))
	log.Printf("%v間隔で指定ページのパトロールを開始します\n", t)
	previousText := make([]string, len(config.Sites))
	for {
		for i, site := range config.Sites {
			// URL先のHTMLサイズを取得
			// size, err := getHTMLSize(site.URL)
			// URL先のselector文字列を取得
			text, err := getWebScraping(site.URL, site.Selector)
			if err != nil {
				break
			}

			// 更新確認
			if previousText[i] == "" {
				// 初期値
				previousText[i] = text
			} else {
				// 変更があったらslackに通知
				if previousText[i] != text {
					message := fmt.Sprintf("更新を検知しました: %v (%v)", site.URL, time.Now())
					log.Printf("%v\n", message)
					if config.Slack.Token != "" && config.Slack.Channelid != "" {
						slackNotify(config.Slack.Token, config.Slack.Channelid, message)
						slackNotify(config.Slack.Token, config.Slack.Channelid, text)
					}
					previousText[i] = text
				}
			}
		}
		// HACK: 周期処理 https://qiita.com/ruiu/items/1ea0c72088ad8f2b841e
		time.Sleep(t)
	}
}

func getWebScraping(url string, selector string) (string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Printf("Error(goquery.NewDocument): %v\n", err)
		return "", err
	}
	t := doc.Find(selector).Text()
	return t, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic("Config File Error...")
	}
	mainLoop(&config)
}
