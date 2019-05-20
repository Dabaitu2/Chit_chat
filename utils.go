package main

import (
	"chitchat/data"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration
var logger *log.Logger

// 快捷打印
func p(a ...interface{}) {
	fmt.Println(a)
}

// 一些初始化配置
func init() {
	// 加载预先配置的文件
	loadConfig()
	// 666 代表 rw- rw- rw-， 指文件的当前属主, 属主组, 其他用户分别都有读写权限，新创建的文件不能有执行权限，故不能是777
	// 777 就是 rwx rwx rwx, unix 还可以用 chmod u+x [target] 来指定权限,
	// u指属主, + 表示增加，x表示可执行权限，其他同理, 共有ugo三种形式, g是属主组，o是其他用户
	file, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	// 所有的log都有INFO前缀, 并写入file中，写入的数据为日期，时间，和当前文件名
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)

}

// 加载json中的设置
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	// 创建json解析器
	decoder := json.NewDecoder(file)
	config = Configuration{}
	// 将decoder解析的结果送入config中
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

// 跳转到错误信息页的快捷方法
func error_message(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	// Redirect方法获取request中一些信息，将其导航到url指明的位置，并表明状态
	// 写返回操作由writer完成
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

// 检查用户是否登陆，是否有一个会话
// 如果会话合法, err会是nil
func session(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

// 解析html模板文件
// html_template -> template 对象
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	// 创建一个名为layout的html template对象
	t = template.New("layout")
	for _, file := range filenames {
		// 把传入的filename格式化写为templates下的路径形式，然后将这个字符串塞进files数组
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	// 解析该路径下的文件并形成模板对象
	t = template.Must(t.ParseFiles(files...))
	return
}

// template对象 + data -> response 的html
func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	// 通过writer将data写入templates模板, 并命名为layout
	templates.ExecuteTemplate(writer, "layout", data)
}

// log函数
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

// version
func version() string {
	return "0.1"
}