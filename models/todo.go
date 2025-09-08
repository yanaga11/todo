package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() (err error) {
	dsn := "root:Password@tcp(127.0.0.1:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&Todo{})
	return nil
}

type Todo struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"size:255" validate:"required"`
}

func CreateTodo(content string) (err error) {
	todo := Todo{Content: content}
	err = db.Create(&todo).Error
	return err
}

func DeleteTodo(id int) (err error) {
	err = db.Delete(&Todo{}, id).Error
	return err
}

func ListTodo(id int) (todo Todo, err error) {
	err = db.First(&todo, id).Error
	return todo, err
}

func GetAllTodos() (todos []Todo, err error) {
    err = db.Find(&todos).Error
    return todos, err
}

func SearchTodos(content string) (todos []Todo, err error) {
    if content == "" {
        return GetAllTodos()
    }
    err = db.Where("content LIKE ?", "%"+content+"%").Find(&todos).Error
    return todos, err
}

func UpdateTodo(t Todo) (err error) {
	err = db.Save(&t).Error
	return err
}
