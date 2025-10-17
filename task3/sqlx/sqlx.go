package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employee struct {
	Id         int64   `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

type Book struct {
	Id     int64   `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

const dsn = "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

var db *sqlx.DB

func initDB() error {
	var err error
	// sqlx.Connect 等价于 database/sql 的 Open + Ping
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	// 配置连接池
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	fmt.Println("连接数据库成功")
	return nil
}

func main() {
	fmt.Println("sqlx操作Mysql")
	initDB()
	defer db.Close()

	fmt.Println("---------------------------employees-------------------------")
	employeeDemo()

	fmt.Println("---------------------------books-------------------------")
	bookDemo()
}

func employeeDemo() {
	// 创建 employees 表
	createTableSql := `CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT,
		name VARCHAR(50) NOT NULL,
		department VARCHAR(100) NOT NULL,
		salary DOUBLE NOT NULL,
		PRIMARY KEY (id)
		)`
	_, err := db.Exec(createTableSql)
	if err != nil {
		log.Fatalln("创建表失败")
	}
	// 插入数据
	INSERTSQL := `insert into employees(name,department,salary) values(?,?,?)`
	db.Exec(INSERTSQL, "小二", "技术部", 9000)
	db.Exec(INSERTSQL, "大黑", "财务部", 6000)
	db.Exec(INSERTSQL, "包包", "技术部", 9500)

	// 使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息
	var employees []Employee
	db.Select(&employees, "SELECT * FROM employees WHERE department=?", "技术部")
	fmt.Printf("技术部员工列表:%+v\n", employees)

	//使用Sqlx查询 employees 表中工资最高的员工信息
	var employee Employee
	err = db.Get(&employee, "SELECT * FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("工资最高的员工信息: %+v\n", employee)
	}
}

func bookDemo() {
	// 创建 books 表
	createTableSql := `CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT,
		title VARCHAR(50) NOT NULL,
		author VARCHAR(50) NOT NULL,
		price DOUBLE NOT NULL,
		PRIMARY KEY (id)
		)`
	_, err := db.Exec(createTableSql)
	if err != nil {
		log.Fatalln("创建表失败")
	}
	// 插入数据
	INSERTSQL := `insert into books(title, author, price) values(?,?,?)`
	db.Exec(INSERTSQL, "c++开发实战", "阿三", 40)
	db.Exec(INSERTSQL, "java开发实战", "刘帅", 55)
	db.Exec(INSERTSQL, "go开发实战", "年华", 58)

	// 查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全
	var books []Book
	err = db.Select(&books, "SELECT * FROM books WHERE price > ?", 50)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	fmt.Printf("价格大于50元的书籍:%+v\n", books)
}
