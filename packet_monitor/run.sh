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

# 30秒間icmpパケットを観測
iptables -A OUTPUT -p icmp,icmpv6 -j NFQUEUE --queue-num 2
echo "キューを設定しました．"
timeout 30 ./packet_monitor
iptables -D OUTPUT -p icmp,icmpv6 -j NFQUEUE --queue-num 2
echo "キューの設定を元に戻しました．"
