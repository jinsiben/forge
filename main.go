package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	// ターゲット解析
	var mode int = -1
	var baseurl string = ""
	var tag string = ""
	mode, baseurl, tag = parseTarget(os.Args)
	fmt.Println("MODE ", mode, baseurl)

	// モード解析
	if mode == 1 {
		downloader(baseurl, tag)
	} else if mode == 2 {
		fmt.Println("MODE: ", mode)
	} else {
		fmt.Println("Unknown MODE: ", mode)
	}

}

// parseTarget(os.Arg)
//
func parseTarget(str []string) (int, string, string) {
	//var mode int = -1
	//var baseurl string = ""

	if len(str) <= 2 {
		fmt.Println("Usage:")
		fmt.Println(str[0], " taskurl tag")
		fmt.Println(str[0], " mode taskurl tag")
		return -1, "", ""
	}

	if len(str) == 3 {
		return 1, str[1], str[2]
	}

	fmt.Println("Too many args")
	return 2, "", ""
}

func downloader(baseurl string, tag string) {

	// フォルダ作成（既にあっても無視）
	os.Mkdir(tag, 0777)
	// フォルダ内にmain.goがなければ同階層のmain.goからフォルダにコピー
	fstr := fmt.Sprintf("./%s/%s", tag, "main.go")
	if _, err := os.Stat(fstr); os.IsNotExist(err) {

		from, err := os.Open("main.go")
		if err != nil {
			fmt.Println("同階層にmain.goがない")
			log.Fatal(err)
		}
		defer from.Close()

		to, err := os.OpenFile(fstr, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			panic(err)
		}
		fmt.Println(fstr, " をコピーした")
	} else {
		fmt.Println(fstr, " は既に存在するので上書きしない")
	}

	// ダウンロード
	url := baseurl
	bodystr := ""

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	bodystrorig := string(byteArray)
	bodystr = bodystrorig

	// 入力サンプル数のカウント
	sicnt := strings.Count(bodystr, "<h3>Sample Input")
	fmt.Println("Sample Input FOUND: ", sicnt)

	// 入力解析
	startIdx := 0
	for i := 0; i < sicnt; i++ {
		clue := strings.Index(bodystr, "<h3>Sample Input")
		prefirst := strings.Index(bodystr[clue:], "<pre>")
		prelast := strings.Index(bodystr[clue:], "</pre>")

		//fmt.Println("FOUND ", i+1, startIdx, clue, clue+prefirst, clue+prelast, bodystr[clue+prefirst:clue+prelast])
		inputstr := bodystr[clue+prefirst+len("<pre>") : clue+prelast-1]

		// ファイル作成(inputN.txt)
		filenamestr := fmt.Sprintf("./%s/input%d.txt", tag, i+1)
		fp, err := os.Create(filenamestr)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fp.Close()
		fp.WriteString(inputstr)

		startIdx = clue + prelast + 1
		bodystr = bodystr[startIdx:]
	}

	// 規定数に足りない分だけinputファイルを追加
	minfilenum := 5 // 最低inputファイル数
	for i := sicnt; i < minfilenum; i++ {
		// ファイル作成(inputN.txt)
		filenamestr := fmt.Sprintf("./%s/input%d.txt", tag, i+1)
		fp, err := os.Create(filenamestr)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fp.Close()
		fp.WriteString("")
	}
	fmt.Println("Sample Input File Num: ", minfilenum)
}
