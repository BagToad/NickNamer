package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"nicknamer/internal/commands"
	"nicknamer/internal/config"
	"nicknamer/internal/nicknames"
)

// Bot represents the Discord bot.
type Bot struct {
	session  *discordgo.Session
	config   *config.Config
	commands *commands.Handler
}

// New creates a new Bot instance.
func New(cfg *config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	nickManager, err := nicknames.NewManager(cfg.DataFile)
	if err != nil {
		return nil, fmt.Errorf("error creating nickname manager: %w", err)
	}

	cmdHandler := commands.NewHandler(cfg, nickManager)

	bot := &Bot{
		session:  session,
		config:   cfg,
		commands: cmdHandler,
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	session.AddHandler(bot.onReady)
	session.AddHandler(bot.onInteractionCreate)

	return bot, nil
}

// Start starts the bot and blocks until interrupted.
func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("error opening Discord connection: %w", err)
	}
	defer b.session.Close()

	if err := b.commands.RegisterCommands(b.session); err != nil {
		return fmt.Errorf("error registering commands: %w", err)
	}

	b.config.Logger.Info("NickNamer bot is now running. Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	b.config.Logger.Info("Shutting down...")
	return nil
}

func (b *Bot) onReady(s *discordgo.Session, r *discordgo.Ready) {
	b.config.Logger.Info("Logged in as", "user", r.User.Username)
	s.UpdateGameStatus(0, "Randomizing nicknames!")
}

func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	b.commands.HandleInteraction(s, i)
}
