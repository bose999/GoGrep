// Package main do Ggrep
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// grep の基本機能である検索文字列がファイルの中にあればシステムアウトするという機能だけ実装
// 追加の機能として指定行数移行から処理を可能としている
func main() {
	// 処理可能なMAXプロセス数を実行環境から取得してセット
	setGOMAXPROCS()
	// 1つのゴルーチンで処理する行数
	splitLineCount := 128
	// 引数から値を取得して変数へ格納
	grepString, filePath, processingStartNumberOfLine := getArgsValue()
	// grep処理を実行
	lineCount := startScan(grepString, filePath, processingStartNumberOfLine, splitLineCount)
	if processingStartNumberOfLine != 1 {
		fmt.Println("GoGrep finished. All Line: " + strconv.Itoa(lineCount) + " Check Start Line: " + strconv.Itoa(processingStartNumberOfLine))
	}
}
 
// ファイルを先頭から行でループさせてsplitLineCountの値で1つのゴルーチンで処理させる行数を決めて文字列チェックさせる
func startScan(grepString string, filePath string, processingStartNumberOfLine int, splitLineCount int) int{
	file, e := os.Open(filePath)
	checkError(e)
	// この関数を抜ける際にファイルを閉じる
	defer file.Close()
	// 現在処理してる行数
	lineCount := 1
	// ゴルーチンで処理する文字列配列
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// ファイルを一行ごとで処理をする
		if lineCount >= processingStartNumberOfLine {
			// 現在処理してる行数が処理開始行数以上の処理
			lines = append(lines, scanner.Text())
			if lineCount%splitLineCount == 0 {
				// splitLineCountに達したらゴルーチンで別スレッドで文字列の配列をチェック
				if lineCount != 0 {
					checkLine(lines, grepString)
					lines = []string{}
				}
			}
		}
		lineCount++
	}
	// 最終行で抜けた際の最後の配列の処理
	checkLine(lines, grepString)
	return lineCount
}

// ゴルーチンで配列に入れた文字列をループ処理して該当文字列があれば出力する
func checkLine(lines []string, grepString string) {
	channel := make(chan int)
	go func(sendChannel chan<- int) {
		// ゴルーチンによる別スレッド処理
		for i := range lines {
			if strings.Index(lines[i], grepString) > -1 {
				fmt.Println(lines[i])
			}
		}
		close(sendChannel)
	}(channel)
	for {
		ok := <-channel
		if ok == 0 {
			break
		}
	}
}

// 引数を処理する 引数が2か3の場合は処理を続行し違う場合はExit
func getArgsValue() (string, string, int) {
	// 処理開始行
	processingStartNumberOfLineString := "1"
	// 検索文字列
	grepString := ""
	// 検索するファイルのパス
	filePath := ""

	flag.Parse()
	args := flag.Args()

	if len(args) == 2 {
		// 引数が2の場合
		grepString = args[0]
		filePath = args[1]
	} else if len(args) == 3 {
		// 引数が3の場合
		grepString = args[0]
		filePath = args[1]
		processingStartNumberOfLineString = args[2]
	} else {
		// その他の場合はエラーで処理をエラーで終了させる
		fmt.Println("args error: need 2 or 3 value. grepString, filepath, processingStartNumberOfLine")
		os.Exit(1)
	}

	// 処理開始行を数字に変更するNGな場合はエラーで終了
	processingStartNumberOfLine, e := strconv.Atoi(processingStartNumberOfLineString)
	if e != nil {
		fmt.Println("args error: processingStartNumberOfLine is not number.")
		os.Exit(1)
	}
	return grepString, filePath, processingStartNumberOfLine
}

// set GOMAXPROCS function
func setGOMAXPROCS() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// when error, print error and exit os
func checkError(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}