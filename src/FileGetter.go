package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

/**
 * FileGetter.go
 * Androidアプリからファイルを取得するクラス
 **/
type FileGetter struct {
	PackageName string
}

// FileType -st

/**
 * FileType
 * ファイルタイプを示すEnum
 **/
type FileType int

const (
	FILE FileType = iota
	DIR
	LINK
	UNKNOWN
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

/**
 * File
 * ファイル情報を持つ構造体
 **/
type File struct {
	PackageName string
	FilePath    string
	FileName    string
	FileType    FileType
}

/**
 * Get
 * ファイルを取得
 **/
func (f *File) Get() {
	// ファイルじゃない場合は無視
	if f.FileType != FILE {
		return
	}

	// ファイル取得
	var fileGetCommand = exec.Command("adb", "exec-out", "run-as", f.PackageName, "cat", "/data/data/"+f.PackageName+f.FilePath+"/"+f.FileName)
	var out, _ = fileGetCommand.Output()

	//  ディレクトリ作成
	exec.Command("mkdir", "-p", "./"+f.PackageName+f.FilePath).Run()

	// ファイル書き込み
	var content = []byte(out)
	ioutil.WriteFile("./"+f.PackageName+f.FilePath+"/"+f.FileName, content, os.ModePerm)
}

// File -ed

/**
 * 現在ディレクトリにAndroidのrun-as領域にあるファイルを全て書き出す
 **/
func (f *FileGetter) GetAll() {
	f.GetPathAll("")
}

/**
 * 現在ディレクトリにAndroidのrun-as領域にあるファイルを指定ディレクトリ以下全て書き出す
 **/
func (f *FileGetter) GetPathAll(path string) {
	// ファイル一覧取得
	var files = f.ListInPath(path)
	for _, file := range files {
		if file.FileType == FILE {
			// ファイルの場合はゲット
			file.Get()
		} else if file.FileType == DIR {
			// ディレクトリの場合は再帰処理
			f.GetPathAll(path + "/" + file.FileName)
		}
	}
}

/**
 * ファイルリストを取得する
 **/
func (f *FileGetter) ListInPath(path string) []File {
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
		}

		// リンクは処理のしようがないので無視
		if fileType == UNKNOWN || fileType == LINK {
			continue
		}

		// ファイルを構築
		var file = File{}
		file.PackageName = f.PackageName
		file.FilePath = path
		file.FileName = regAns[0][4]
		file.FileType = fileType

		// 結果にアペンド
		files = append(files, file)
	}

	return files
}
