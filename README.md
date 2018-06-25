# Patrolman
- 複数サイトの監視と差分検知
- Slack通知
- 設定はconfig.json

# config.jsonの例
```json
{
  "sites": [
    {
      "selector": "div.contents",
      "url": "http://localhost:8000"
    },
    {
      "selector": "div > ul",
      "url": "http://localhost:8000"
    }
  ],
  "slack": {
    "token": "<slackbotのtoken>",
    "channelId": "<通知を流すChannelのID>"
  },
  "intervalHour": 10
}
```

# 使い方
bin以下に各プラットフォーム向けのバイナリがあります。

- Windows(64bit): bin/Patrolman.exe
- Linux(64bit): bin/Patrolman_linux
- Mac(64bit): bin/Patrolman_mac

同階層にconfig.jsonを用意してバイナリを実行します。

# 実装
intervalHourごとにselectorで指定したサイトの文字列を取得し、前に取得した値と変化があったら通知。
