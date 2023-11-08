package telebot

import (
	"time"

	"github.com/rtsoy/software-design-patterns-project/config"
	"github.com/rtsoy/software-design-patterns-project/internal/repository"
	tele "gopkg.in/telebot.v3"
)

const _defaultLongPollerTimeout = 10 * time.Second

// Bot represents a Telegram bot instance and its associated functionality.
type Bot struct {
	bot             *tele.Bot
	repo            *repository.Repository
	lastSentMessage *lastSentMessage

	cardAttachmentHandler CardHandler
}

// New creates a new Bot instance with the provided repository and starts it.
// It initializes the Telegram bot with the provided token and long poller settings.
func New(repo *repository.Repository) (*Bot, error) {
	// Initialize the Telegram bot settings.
	pref := tele.Settings{
		Token:  config.Get().TelegramToken,
		Poller: &tele.LongPoller{Timeout: _defaultLongPollerTimeout},
	}

	// Create a new Bot instance with the Telegram bot and the repository.
	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	// Create a new Bot instance with the Telegram bot and the repository.
	bot := &Bot{
		bot:             b,
		repo:            repo,
		lastSentMessage: newLastSentMessage(),
	}

	// Create a chain of responsibility for card attachment verification.
	bot.cardAttachmentHandler = &ConcreteCardHandler{
		next: &CardNumberHandler{
			next: &CardExpDateHandler{
				next: &CardHolderHandler{
					next: &CardCVVHandler{
						next: nil, // The last handler in the chain.
						b:    bot,
					},
					b: bot,
				},
				b: bot,
			},
			b: bot,
		},
		b: bot,
	}

	// Initialize message handlers.
	bot.initHandlers()

	// Start the Telegram bot.
	bot.start()

	return bot, nil
}

// start starts the Telegram bot instance.
func (b *Bot) start() {
	b.bot.Start()
}

// Stop stops the Telegram bot instance.
func (b *Bot) Stop() {
	b.bot.Stop()
}
