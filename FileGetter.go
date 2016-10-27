package main

import (
	"os/exec"
	"regexp"
	"strings"
)

type FileGetter struct {
	PackageName string
}

// FileType -st

type FileType int

const (
	UNKNOWN FileType = iota
	FILE
	DIR
	LINK
)

func (f *FileType) ToString() string {
	switch *f {
	case FILE:
		return "file"
	case DIR:
		return "dir"
	case LINK:
		return "link"
	}

	return "unknown"
}

// FileType -ed

// File -st

type File struct {
	FileName string
	FileType FileType
}

// File -ed

/**
 * 現在ディレクトリにAndroidのrun-as領域にあるファイルを全て書き出す
 **/
func (f *FileGetter) Get() {

}

func (f *FileGetter) searchGet(path string) {

}

/**
 * ファイルリストを取得する
 **/
func (f *FileGetter) list(path string) []File {
	// 指定アプリの指定パスでlsコマンド
	var out, _ = exec.Command("adb", "shell", "run-as", f.PackageName, "ls", "/data/data/"+f.PackageName+path, "-n").Output()
	var outStr = string(out)

	// ファイル解析
	var files = make([]File, 0)
	for {
		// 一行の終了インデックス検索
		var endIndex = strings.Index(outStr, "\r\n")

		// もうデータが残っていない場合は終了
		if endIndex == -1 {
			break
		}

		// 一行をlsの結果から抜き取る
		var line = outStr[0:endIndex]
		outStr = outStr[endIndex+2 : len(outStr)]

		// ファイルタイプ 日付 時間 ファイル名 を正規表現で検索
		var reg, _ = regexp.Compile("(.).* (....-..-..) (..:..) (.*)")
		var regAns = reg.FindAllStringSubmatch(line, -1)

		// ファイルタイプ解析
		var fileType FileType
		switch regAns[0][1] {
		case "-":
			fileType = FILE
		case "d":
			fileType = DIR
		case "l":
			fileType = LINK
		default:
			fileType = UNKNOWN
			continue
		}

		// リンクは処理のしようがないので無視
		if fileType == UNKNOWN || fileType == LINK {
			continue
		}

		// ファイルを構築
		var file = File{}
		file.FileName = regAns[0][4]
		file.FileType = fileType

		// 結果にアペンド
		files = append(files, file)
	}

	return files
}
