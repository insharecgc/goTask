package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用户
type User struct {
	gorm.Model        // 嵌入默认字段：ID（主键）、CreatedAt、UpdatedAt、DeletedAt（软删除）
	Name       string `gorm:"type:varchar(50);not null;comment:用户名"`
	Age        int    `gorm:"type:int;comment:年龄"`
	PostNum    int    `gorm:"type:int;default:0;column:post_num;comment:发布文章数"`
}

// 文章
type Post struct {
	gorm.Model
	UserId        int64     `gorm:"column:user_id;not null;comment:文章发布用户ID"`
	Title         string    `gorm:"type:varchar(100);not null;comment:文章标题"`
	Content       string    `gorm:"type:varchar(1000);not null;comment:文章内容"`
	CommentStatus string    `gorm:"type:varchar(255);comment:文章评论状态"`
	Comments      []Comment // 关联评论字段：一对多
}

// 评论
type Comment struct {
	gorm.Model
	PostId  int    `gorm:"column:post_id;not null;comment:文章ID"`
	Content string `gorm:"type:varchar(255);not null;comment:评论内容"`
}

// 定义查询结果结构体
type MaxCountPost struct {
	Post     Post `gorm:"embedded"`        // 这个映射一定要要，否则无法赋值，字段首字母也必须大写，否则无法导出
	MaxCount int  `gorm:"column:maxCount"` // 最大评论数，这个映射一定要要，否则无法赋值
}

func initDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}
	return db
}

func main() {
	fmt.Println("gorm连接mysql操作CRUD")
	dsn := "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	db := initDB(dsn)

	fmt.Println("创建表，如果已存在会跳过")
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// 初始化数据
	// db.Create(&User{Name: "大曾子", Age: 28})
	// db.Create(&User{Name: "王老板", Age: 32})
	// db.Create(&Post{UserId: 1, Title: "标题1", Content: "内容1"})
	// db.Create(&Post{UserId: 2, Title: "标题2", Content: "内容2"})
	// db.Create(&Post{UserId: 1, Title: "标题3", Content: "内容3"})
	// db.Create(&Post{UserId: 1, Title: "标题4", Content: "内容4"})
	// db.Create(&Comment{Content: "评论1", PostId: 1})
	// db.Create(&Comment{Content: "评论2", PostId: 1})
	// db.Create(&Comment{Content: "评论3", PostId: 2})
	// db.Create(&Comment{Content: "评论4", PostId: 3})
	// db.Create(&Comment{Content: "评论5", PostId: 4})

	db.Delete(&Comment{Model: gorm.Model{ID: 5}})

	//查询用户1的所有文章及评论
	var posts []Post
	db.Preload("Comments").Where("user_id=?", 1).Find(&posts)
	for _, post := range posts {
		fmt.Printf("用户1对应的文章标题:%s,内容:%s,评论:[", post.Title, post.Content)
		for _, comment := range post.Comments {
			fmt.Printf("%s;", comment.Content)
		}
		fmt.Printf("]\n")
	}

	// 查询评论数量最多的文章信息
	var maxCountPost MaxCountPost
	db.Table("posts").
		Select("posts.*, count(comments.id) as maxCount").
		Joins("left join comments on posts.id = comments.post_id").
		Group("posts.id").
		Order("maxCount desc").
		Limit(1).
		Scan(&maxCountPost)
	fmt.Printf("评论最多的文章标题:%s,内容:%s,评论数量:%d\n", maxCountPost.Post.Title, maxCountPost.Post.Content, maxCountPost.MaxCount)
}

// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段
func (p *Post) AfterCreate(db *gorm.DB) error {
	log.Println("文章创建完成后，触发钩子函数")
	res := db.Model(&User{}).Where("id =?", p.UserId).Update("post_num", gorm.Expr("post_num + 1"))
	if res.Error != nil {
		log.Printf("更新用户文章数失败:%v\n", res.Error)
		return res.Error
	}
	log.Printf("用户%d的文章数更新成功\n", p.UserId)
	return nil
}

/*
为 Comment 模型添加一个钩子函数，
在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"
*/
func (c *Comment) AfterDelete(db *gorm.DB) error {
	log.Println("评论删除后，触发钩子函数")
	postId := c.PostId
	// 如果是根据ID删除评论，需要先根据评论ID，找到文章ID
	if postId == 0 {
		var comment Comment
		res := db.Unscoped().Select("post_id").First(&comment, c.ID)
		if res.Error != nil {
			log.Printf("评论记录不存在：%v\n", res.Error)
			return res.Error
		}
		postId = comment.PostId
	}

	// 检查当前文章对应的评论数是否为0
	var count int64
	res := db.Model(&Comment{}).Where("post_id = ?", postId).Count(&count)
	if res.Error != nil {
		log.Printf("查询文章的评论数失败: %v\n", res.Error)
		return res.Error
	}

	if count == 0 {
		log.Println("当前文章的评论数为0，需要更新文章评论状态字段为【无评论】")
		res = db.Model(&Post{}).Where("id = ?", postId).Update("comment_status", "无评论")
		if res.Error != nil {
			log.Printf("更新文章评论状态字段失败: %v\n", res.Error)
		}
	}
	return nil
}
