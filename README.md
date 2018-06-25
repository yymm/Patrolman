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
      "url": "http://localhost:8888"
    },
    {
      "url": "http://localhost:5000"
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
intervalHourごとにサイトの文字列を取得し、前に取得した値と変化があったら通知します。

sites内のselectorがある場合は指定の要素の差分を使用し、selectorの指定がない場合は全文の比較を行います。
