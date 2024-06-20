package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-telegram/ui/keyboard/reply"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/gocolly/colly"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New("", opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)

}

func onReplyKeyboardSelect(ctx context.Context, b *bot.Bot, update *models.Update) {
	switch update.Message.Text {
	case "/alive":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Всё хорошо!",
		})
	case "/time":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   time.Now().Format("2006-01-02 15:04:05"),
		})
	case "/book":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   parse(),
		})
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := reply.New(
		b,
		reply.WithPrefix("reply_keyboard"),
		reply.IsSelective(),
		reply.IsOneTimeKeyboard(),
		reply.ResizableKeyboard(),
	).
		Button("/book", b, bot.MatchTypeExact, onReplyKeyboardSelect).
		Button("/time", b, bot.MatchTypeExact, onReplyKeyboardSelect).
		Button("/alive", b, bot.MatchTypeExact, onReplyKeyboardSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Выберите команду",
		ReplyMarkup: kb,
	})

	switch update.Message.Text {
	case "/alive":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Всё хорошо",
		})
	case "/time":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   time.Now().Format("2006-01-02 15:04:05"),
		})
	case "/book":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   parse(),
		})
	}

}

func parse() string {
	c := colly.NewCollector()

	buff := bytes.Buffer{}
	c.OnHTML("div[class=blog-content]", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, element *colly.HTMLElement) {
			if (i > 10 && i < 20) || i == 21 {
				buff.WriteString(fmt.Sprintf("%s\n", element.Text))
			}

		})
	})

	c.Visit("https://it-vacancies.ru/blog/10-knig-po-razrabotke-na-go-kotorye-pomogut-vam-stat-ekspertom/?ysclid=lxnaykrfyy923951004")

	return buff.String()
}
