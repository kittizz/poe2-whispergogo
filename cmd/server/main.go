package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/gofiber/fiber/v2"
)

var telegramBot *bot.Bot

func main() {
	var err error
	telegramBot, err = bot.New(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	go startBot()

	app := fiber.New()

	app.Post("/send", sendNotification)
	app.Get("/verify/:chatId", verifyChatId)

	log.Fatal(app.Listen(":80"))
}
func startBot() {

	handler := func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}
		message := fmt.Sprintf(
			"Welcome\\! Here's your Chat ID:\n\n`%d`\n\nYou can use this ID to receive notifications through our system\\. Tap the ID to copy",
			update.Message.Chat.ID,
		)

		_, err := telegramBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      message,
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}

	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handler)

	telegramBot.Start(context.Background())

}
func sendNotification(c *fiber.Ctx) error {
	var req struct {
		ChatID  string `json:"chat_id"`
		Message string `json:"message"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	fmt.Println(req)

	_, err := telegramBot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: req.ChatID,
		Text:   req.Message,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send message"})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func verifyChatId(c *fiber.Ctx) error {
	chatId := c.Params("chatId", "")

	fmt.Println("Verifying chat ID:", chatId)

	// ทดสอบส่งข้อความไปยัง chat ID
	m, err := telegramBot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatId,
		Text: `🎮 Successfully connected to WhisperGoGo!
Your trade notifications from Path of Exile 2 will be sent to this chat.
Happy trading, Exile! 💎`,
	})

	if err != nil {
		return c.JSON(fiber.Map{
			"valid": false,
			"error": "Chat ID not found or bot cannot send message to this chat",
		})
	}

	return c.JSON(fiber.Map{
		"valid":   true,
		"chat_id": fmt.Sprintf("%d", m.Chat.ID),
	})
}
