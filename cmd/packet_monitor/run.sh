#!/bin/bash

FILE="packet_monitor"  # 実行ファイル

# エラー処理
if [ ! -e $FILE ]; then
  echo -e "${FILE}が見つかりませんでした．\ngoコマンドを使用してビルドしてください．" 1>&2
  exit 1
elif [ "$EUID" -ne 0 ]; then
  echo "root権限で実行してください．" 1>&2
  exit 1
fi

# # 10秒間全てのパケットを観測
iptables -I DOCKER-USER 1 -p all -j NFQUEUE --queue-num 2
echo "キューを設定しました．"
timeout 10 ./packet_monitor
iptables -D DOCKER-USER  -p all -j NFQUEUE --queue-num 2
echo "キューの設定を元に戻しました．"
