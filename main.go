package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Urls  []string `json:"urls"`
	Slack struct {
		Token     string `json:"token"`
		Channelid string `json:"channelid"`
	} `json:"slack"`
	Intervalhour time.Duration `json:"intervalhour"`
}

func slackNotify(token, channelId, message string) {
	api := slack.New(token)
	params := slack.PostMessageParameters{
		AsUser: true, // ユーザーとしてpostする(招待されたページのみpost可能)
		//LinkNames: 1, // メンションを可能にする
		// 以下のオプションを使うと招待されていないチャンネルでも
		// 自由なIcon、ユーザー名でpostできる
		//IconURL:  "https://pbs.twimg.com/media/DUnrnDHVQAAiXff.jpg",
		//Username: "まつぼっくり",
	}
	_, _, err := api.PostMessage(channelId, message, params)
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

func getHtmlSize(url string) (int, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error(http.Get): %v\n", err)
		return 0, err
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf(" Error(Read HTML): %v\n", err)
		return 0, err
	}
	return len(bytes), nil
}

func mainLoop(config *Config) {
	previousSizes := make([]int, len(config.Urls))
	for {
		for i, url := range config.Urls {
			// URL先のHTMLサイズを取得
			size, err := getHtmlSize(url)
			if err != nil {
				break
			}
			// 更新確認
			if previousSizes[i] == 0 {
				// 初期値
				previousSizes[i] = size
			} else {
				// 変更があったらslackに通知
				if previousSizes[i] != size {
					message := fmt.Sprintf("Change => %v (%v)", url, time.Now())
					log.Printf("%v\n", message)
					if config.Slack.Token != "" && config.Slack.Channelid != "" {
						slackNotify(config.Slack.Token, config.Slack.Channelid, message)
					}
					previousSizes[i] = size
				}
			}
		}
		// HACK: 周期処理 https://qiita.com/ruiu/items/1ea0c72088ad8f2b841e
		time.Sleep(config.Intervalhour * time.Hour)
	}
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic("Config File Error...")
	}
	mainLoop(&config)
}
