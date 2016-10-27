package main

import (
	"fmt"
	"os"
)

func main() {
	var packageName string
	var clean bool

	// args1 パッケージ名
	if len(os.Args) >= 2 {
		packageName = os.Args[1]
	} else {
		fmt.Println("取得するアプリデータのパッケージ名を入力してください")
		fmt.Scanln(&packageName)
	}

	// args2 クリーンの可否
	clean = false
	if len(os.Args) >= 3 {
		if os.Args[2] == "clean" {
			clean = true
		}
	}
	if clean == true {
		os.RemoveAll(packageName)
	}

	fmt.Println("パッケージ `" + packageName + "` からデータを取得します")

	var fileGetter = FileGetter{packageName}
	fileGetter.GetAll()

	fmt.Println("取得完了")
}
