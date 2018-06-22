# Patrolman
- 複数サイトの監視と差分検知
- Slack通知
- 設定はconfig.json

# config.jsonの例
```json
{
  "urls": [
    "http://localhost:8000",
    "http://localhost:8080"
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
intervalHourごとにurlsのサイトのHTMLのサイズを取得し前に取得した値と変化があったら通知。

# 既知の課題
- HTMLのサイズがint型の最大値を超えたときの挙動が不明
- 具体的にどの部分に変更があったかまではわからない
