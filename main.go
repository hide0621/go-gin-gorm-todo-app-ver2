package main

import (
	"strconv" // for strconv.Atoi

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	/*
		sqlite3 は読み込まなければなりません
		しかし実際のデータベースの操作は、gorm で行います
		そのため、sqlite3は使わないことを、 _ で明示しなければなりません
	*/
	_ "github.com/mattn/go-sqlite3"
)

// モデルの宣言（データ構造の宣言）
type Todo struct {
	gorm.Model
	Memo string
}

func main() {
	r := gin.Default()
	r.Static("styles", "./styles")
	r.LoadHTMLGlob("templates/*.html")

	dbInit()

	// List
	r.GET("/", func(c *gin.Context) {
		todos := dbGetAll()
		c.HTML(200, "index.html", gin.H{"todos": todos})
	})

	// Create
	r.POST("/new", func(c *gin.Context) {
		memo := c.PostForm("memo")
		dbCreate(memo)
		c.Redirect(302, "/")
	})

	//Delete
	r.GET("delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbDelete(id)
		c.Redirect(302, "/")
	})

	// Edit
	r.GET("/edit/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		c.HTML(200, "edit.html", gin.H{"todo": todo})
	})

	// Update
	r.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		memo := c.PostForm("memo")
		dbUpdate(id, memo)
		c.Redirect(302, "/")
	})

	// Run Server
	r.Run() // 引数を指定しないと、":8080" が指定されたことになります
}

// データベースのマイグレート（データベースの初期化）
func dbInit() {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()
	db.AutoMigrate(&Todo{})
}

// データの作成
func dbCreate(memo string) {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()
	db.Create(&Todo{Memo: memo})
}

// データの削除
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
}

// データの更新
func dbUpdate(id int, memo string) {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()
	var todo Todo
	db.First(&todo, id)
	todo.Memo = memo
	db.Save(&todo)
}

// データのすべてを取得
func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()
	var todos []Todo
	db.Find(&todos)
	return todos
}

// データの一つを取得
func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "todo.sqlite3")
	if err != nil {
		panic("railed to connect database\n")
	}
	defer db.Close()
	var todo Todo
	db.First(&todo, id)
	return todo
}
