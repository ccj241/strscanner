package strscan

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var ch = make(chan string, 500)

func FindStr(dirPath, text, suffix string, thread int) {
	var pathArray [500]string
	var count = 0
	//files = make([]string, 0, 30)
	filepath.Walk(dirPath, func(filename string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			fisuffix := filepath.Ext(fi.Name())
			fisuffix = strings.Trim(fisuffix, ".")
			if fisuffix == suffix || suffix == "*" && fi.Size() <= 10240*1024 {
				pathArray[count] = filename
				count++
				if count >= thread {
					count = 0
					for _, filecache := range pathArray[:thread] {
						go findText(filecache, text)
						if len(ch) >= thread {
							time.Sleep(time.Millisecond * 1500)
						}
					}
					time.Sleep(time.Millisecond * 2500)
				}

			}
		}
		return nil
	})
	for _, filecache := range pathArray[:count] {
		go findText(filecache, text)
	}

}

func findText(file string, text string) {
	fi, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	if strings.Index(string(fd), text) > 0 {
		ch <- file
		return
	}
}

func Strscan(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("strscan.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		fmt.Println("path : ", r.FormValue("path"))
		fmt.Println("keywords : ", r.FormValue("keywords"))
		fmt.Println("suffix : ", r.FormValue("suffix"))
		fmt.Println("thread : ", r.FormValue("thread"))
		path, keywords, suffix, thread := r.FormValue("path"), r.FormValue("keywords"), r.FormValue("suffix"), r.FormValue("thread")
		threadnum, _ := strconv.Atoi(thread)
		log.Println("开始扫描")
		_, err := os.Stat(path)
		if err != nil {
			log.Println(err)
		} else {
			FindStr(path, keywords, suffix, threadnum)
			fmt.Fprintf(w, "在文件夹【 %s 】含有关键字【 %s 】的路径有： \n", path, keywords)
			for {
				i, ok := <-ch
				if ok {
					fmt.Fprintf(w, "    %s\n", i)
				}
				if len(ch) <= 0 { // 如果现有数据量为0，跳出循环
					return
				}
			}
			close(ch)
			log.Println("查询完成")
		}
	}
}
