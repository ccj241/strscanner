package shellscan

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ch = make(chan string, 500)

func readfile(dirPath, suffix string, thread int) {
	var pathArray [500]string
	var count = 0
	filepath.Walk(dirPath, func(filename string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			fisuffix := filepath.Ext(fi.Name())
			fisuffix = strings.Trim(fisuffix, ".")
			if fisuffix == suffix || suffix == "*" && fi.Size() <= 10240*1024 {
				pathArray[count] = filename
				count++
				if count >= thread {
					count = 0
					//读取出每一个符合条件的文件之后进行处理
					for _, filecache := range pathArray[:thread] {
						//log.Println(filecache)
						regStr(filecache)
						//如果当前线程数大于设定的线程，暂停1.5s
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
	for _, filecache := range pathArray[:thread] {
		//log.Println(filecache)
		regStr(filecache)
	}
}

func regStr(filename string) {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}

	//defer rulefile.Close()
	rulefile, _ := os.Open("./shellscan/rule.txt")
	br := bufio.NewReader(rulefile)
	for {
		a, _, c := br.ReadLine()
		if string(a) != "" {
			rule := strings.Replace(string(a), "\n", "", -1)
			//log.Println("正则：" + rule)
			r, _ := regexp.Compile(rule)
			if r.MatchString(string(fd)) {
				log.Println(r.FindString(string(fd)))
				ch <- filename
				log.Println(filename)
				return
			}
		}
		if c == io.EOF {
			break
		}
	}

}

func ShellScan(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("shellscan.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		//读取所有文件，输出结果
		path, suffix, thread := r.FormValue("path"), r.FormValue("suffix"), r.FormValue("thread")
		log.Println(path)
		threadnum, _ := strconv.Atoi(thread)
		_, err := os.Stat(path)
		if err != nil {
			log.Println(err)
		} else {
			readfile(path, suffix, threadnum)
			for {
				i, ok := <-ch
				if ok {
					fmt.Fprintf(w, "    路径%s\n", i)
				}
				if len(ch) <= 0 { // 如果现有数据量为0，跳出循环
					return
				}
			}
			close(ch)
		}

	}

}
