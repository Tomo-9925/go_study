package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tomo-9925/go_study/pkg/yaml"
)

func main() {
	// logrusの設定
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	// セキュリティポリシーのパース
	policy, err := yaml.ParseSecurityPolicy("security_policy.yml")
	if err != nil {
		logrus.Fatalln(err)
	}

	// 結果の出力
	logrus.Infof("取得結果: %+v\n", policy)
}
