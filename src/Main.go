package main

import (
	"fmt"
	"os"
)

func main() {
	var packageName string

	if len(os.Args) >= 2 {
		fmt.Println("取得するアプリデータのパッケージ名を入力してください")
		packageName = os.Args[1]
	} else {
		fmt.Scanln(&packageName)
	}

	fmt.Println("パッケージ `" + packageName + "` からデータを取得します")

	var fileGetter = FileGetter{packageName}
	fileGetter.GetAll()

	fmt.Println("取得完了")
}
