package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 导入驱动（下划线表示只执行 init 函数）
)

var db *sql.DB

func initDB() error {
	dsn := "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = sql.Open("mysql", dsn) // 打开连接（只是初始化，不会立即连接）
	if err != nil {
		return err
	}
	// 验证连接是否成功
	if err = db.Ping(); err != nil {
		return err
	}
	// 配置连接池
	db.SetMaxOpenConns(100) // 设置最大连接数
	db.SetMaxIdleConns(20)  // 设置最大空闲连接数

	fmt.Println("连接数据库成功")
	return nil
}

// 定义Student表的结构体
type Student struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade string `json:"grade"`
	Dr    int    `json:"dr"`
}

type Accounts struct {
	Id      int
	Balance float64 //账户余额
}

type Transactions struct {
	Id            int
	FromAccountId int     //转出账户ID
	ToAccountId   int     //转入账户ID
	Amount        float64 //转账金额
}

func main() {
	fmt.Println("原生Mysql操作")
	if err := initDB(); err != nil {
		fmt.Printf("init db failed, err:%v\n", err)
		return
	}
	defer db.Close()

	fmt.Println("---------------------------CRUD-------------------------")
	curdDemo()

	fmt.Println("---------------------------事务-------------------------")
	transactionDemo()
}

// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
// 在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，
// 向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务
func transactionDemo() {
	err := transfer(1, 2, 100)
	if err != nil {
		fmt.Printf("转账失败：%v\n", err)
	} else {
		fmt.Printf("转账成功\n")
	}
}

func transfer(fromAccountId, toAccountId int, amount float64) error {
	// 1.开启事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer func() {
		if p := recover(); p != nil { // 捕获 panic（如SQL错误导致的崩溃）
			tx.Rollback() // 发生panic时回滚事务
			fmt.Printf("事务 panic 回滚: %v", p)
		}
	}()
	// 2.检查账户 A 的余额是否足够
	var fromBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", fromAccountId).Scan(&fromBalance)
	if err != nil {
		tx.Rollback() // 查询失败回滚事务
		return fmt.Errorf("查询账户余额失败: %v", err)
	}
	// 3.检查余额是否足够
	if fromBalance < amount {
		tx.Rollback() // 余额不足回滚事务
		return fmt.Errorf("账户余额不足,当前余额: %.2f, 需转账: %.2f", fromBalance, amount)
	}
	// 4.从账户 A 扣除金额
	result, err := tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromAccountId)
	if err != nil {
		tx.Rollback() // 扣款失败回滚事务
		return fmt.Errorf("账户扣款失败: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		tx.Rollback() // 扣款失败回滚事务
		return fmt.Errorf("账户扣款失败，影响行数: %d, err: %v", rowsAffected, err)
	}
	// 5.向账户 B 增加金额
	result, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toAccountId)
	if err != nil {
		tx.Rollback() // 入账失败回滚事务
		return fmt.Errorf("账户入账失败: %v", err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		tx.Rollback() // 入账失败回滚事务
		return fmt.Errorf("账户入账失败，影响行数: %d, err: %v", rowsAffected, err)
	}
	// 6.记录转账信息
	_, err = tx.Exec("INSERT INTO transactions (from_account_id, to_account_id, amount) VALUES (?, ?, ?)",
		fromAccountId, toAccountId, amount)
	if err != nil {
		tx.Rollback() // 记录失败回滚事务
		return fmt.Errorf("记录转账信息失败: %v", err)
	}
	// 7.提交事务
	if err = tx.Commit(); err != nil {
		tx.Rollback() // 提交失败回滚事务
		return fmt.Errorf("提交事务失败: %v", err)
	}
	fmt.Printf("账户 %d 向账户 %d 转账 %.2f 元成功\n", fromAccountId, toAccountId, amount)
	return nil
}

func curdDemo() {
	// 插入数据
	u := &Student{Name: "张三", Age: 20, Grade: "三年级"}
	id, err := insertStudent(u)
	if err != nil {
		fmt.Printf("保存失败：%v\n", err)
	} else {
		fmt.Printf("保存成功，ID：%d\n", id)
	}

	// 查询数据
	student, err := queryStudentById(1)
	if err != nil {
		fmt.Printf("查询失败：%v\n", err)
	} else if student != nil {
		fmt.Printf("查询到学生：%+v\n", *student)
	}

	// 更新用户
	student.Name = "张三"
	student.Grade = "四年级"
	err = updateStudentByName(student)
	if err != nil {
		fmt.Printf("更新失败：%v", err)
	} else {
		fmt.Printf("更新成功")
	}

	// 查询年龄大于18岁的学生
	students, err := queryStudentsByAge(18)
	if err != nil {
		fmt.Printf("查询失败：%v\n", err)
	} else {
		fmt.Printf("查询到大于18岁的学生：%+v\n", students)
	}

	// 删除年龄小于15岁的学生
	affect, err := deleteStudentAge(15)
	if err != nil {
		fmt.Printf("删除失败：%v", err)
	} else {
		fmt.Printf("删除成功，影响行数：%d", affect)
	}
}

// 插入数据
func insertStudent(u *Student) (int64, error) {
	sqlStr := "insert into student(name, age, grade, dr) values(?, ?, ?, ?)"
	result, err := db.Exec(sqlStr, u.Name, u.Age, u.Grade, 0)
	if err != nil {
		return 0, err
	}
	// 获取插入的自增id
	id, err := result.LastInsertId()
	return id, err
}

// 根据id查询单条数据
func queryStudentById(id int) (*Student, error) {
	sqlStr := "select id, name, age, grade, dr from student where id=? and dr=0"
	// QueryRow用于查询单条数据
	row := db.QueryRow(sqlStr, id)
	var u Student
	// 必须调用Scan方法扫描结果到结构体，否则持有的数据库连接不会被释放（字段顺序与SQL一致）
	err := row.Scan(&u.Id, &u.Name, &u.Age, &u.Grade, &u.Dr)
	if err == sql.ErrNoRows { // 没有查询到结果
		return nil, nil
	}
	return &u, err
}

// 查询多条数据
func queryStudentsByAge(age int) ([]Student, error) {
	sqlStr := "select id, name, age, grade, dr from student where age > ? and dr=0"
	// Query 用于查询多条数据
	rows, err := db.Query(sqlStr, age)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 关闭rows，释放持有的数据库连接
	var students []Student
	for rows.Next() {
		var u Student
		err := rows.Scan(&u.Id, &u.Name, &u.Age, &u.Grade, &u.Dr)
		if err != nil {
			return nil, err
		}
		students = append(students, u)
	}
	// 检查遍历时是否出错
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

// 更新数据
func updateStudentByName(u *Student) error {
	sqlStr := "update student set age=?, grade=? where name=?"
	_, err := db.Exec(sqlStr, u.Age, u.Grade, u.Name)
	return err
}

// 删除年龄小于xx数据（逻辑删除）
func deleteStudentAge(age int) (int64, error) {
	sqlStr := "update student set dr=1 where age < ?"
	result, err := db.Exec(sqlStr, age)
	if err != nil {
		return 0, err
	}
	// 受影响的行数
	return result.RowsAffected()
}
