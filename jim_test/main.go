package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pelletier/go-toml"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func Get_msg(name string, code int) Message {
	return Message{
		Msg:  fmt.Sprintf("this page is [%s]", name),
		Code: code,
	}
}

func GetConfs() (string, string, string) {
	// 解析TOML文件
	mysql_dns := os.Getenv("mysql_dns")
	mysql_user := os.Getenv("mysql_user")
	mysql_pass := os.Getenv("mysql_pass")
	config, err := toml.LoadFile("./conf/config.toml")
	if err != nil {
		log.Println("解析TOML文件时发生错误:", err)
		os.Exit(-1)
	}
	// 读取配置信息
	databaseHost := config.Get("mysql.dns").(string)
	sqlStr := config.Get("mysql.sql").(string)
	name := config.Get("mysql.name").(string)
	// 使用配置信息
	log.Println("数据库dns:", databaseHost)
	replace1 := strings.ReplaceAll(databaseHost, "${mysql_user}", mysql_user)
	replace2 := strings.ReplaceAll(replace1, "${mysql_pass}", mysql_pass)
	replace3 := strings.ReplaceAll(replace2, "${mysql_dns}", mysql_dns)
	log.Println("替换后连接串：", replace3)
	log.Println("=====>", sqlStr)
	return databaseHost, name, sqlStr
}

func ret_msg(name string, code int, err error, c *gin.Context) {
	msg := Get_msg(name, code)
	if err != nil {
		log.Println("异常：", err)
		return
	}
	c.JSON(http.StatusOK, msg)
}

func Pro(c *gin.Context) {
	// 连接数据库
	confs, name, sqlStr := GetConfs()
	fmt.Println("--->", confs)
	db, err := sql.Open("mysql", confs)
	if err != nil {
		log.Println("连接数据库时发生错误:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}
	defer db.Close()

	// 测试连接是否成功
	err = db.Ping()
	if err != nil {
		log.Println("连接测试失败:", err)
		c.JSON(http.StatusCreated, gin.H{"error": err})
		return
	}
	log.Println("连接成功！")

	// 执行查询
	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Println("查询数据时发生错误:", err)
		return
	}
	defer rows.Close()

	// 处理查询结果
	log.Printf("执行sql:%v\n 结果：", sqlStr)
	for rows.Next() {
		var col1, col3 string
		var col2 int64
		err = rows.Scan(&col1, &col2, &col3)
		if err != nil {
			log.Println("扫描行时发生错误:", err)
			return
		}
		log.Println(col1, col2, col3)
	}
	ret_msg(name, http.StatusOK, err, c)
}

func main() {
	// 创建Gin路由引擎
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// 定义一个GET请求的路由处理函数
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})
	router.GET("/test", Pro)
	// 启动HTTP服务器，默认监听在本地的8080端口
	_ = router.Run(":8000")
}
