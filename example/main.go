package main

import (
	"github.com/0x16F/qparser"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Request struct {
	Name  string `query:"name"`
	Email string `query:"email"`
}

type User struct {
	Id    int
	Name  string
	Email string
}

func main() {
	db, err := gorm.Open(&gorm.Config{})
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		var req Request

		if err := c.QueryParser(&req); err != nil {
			return err
		}

		options, err := qparser.ParseStruct(&req)
		if err != nil {
			return err
		}

		var users []User

		if err := options.Apply(db.WithContext(c.Context()).Model(&User{})).Find(&users).Error; err != nil {
			return err
		}

		return nil
	})

	app.Listen(":3000")
}
