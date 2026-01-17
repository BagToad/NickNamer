package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"nicknamer/internal/config"
	"nicknamer/internal/nicknames"
)

// Handler manages slash command registration and execution.
type Handler struct {
	config      *config.Config
	nickManager *nicknames.Manager
	commands    []*discordgo.ApplicationCommand
}

// NewHandler creates a new command handler.
func NewHandler(cfg *config.Config, nm *nicknames.Manager) *Handler {
	return &Handler{
		config:      cfg,
		nickManager: nm,
		commands:    buildCommands(),
	}
}

func buildCommands() []*discordgo.ApplicationCommand {
	minWords := float64(1)
	maxWords := float64(10)

	return []*discordgo.ApplicationCommand{
		{
			Name:        "remember",
			Description: "Add words to the nickname pool",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "words",
					Description: "Space-separated words to remember",
					Required:    true,
				},
			},
		},
		{
			Name:        "forget",
			Description: "Remove a word from the nickname pool",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "word",
					Description: "The word to forget",
					Required:    true,
				},
			},
		},
		{
			Name:        "forgetall",
			Description: "Clear all words from the nickname pool",
		},
		{
			Name:        "names",
			Description: "List all words in the nickname pool",
		},
		{
			Name:        "randomizeme",
			Description: "Randomize your own nickname",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "words",
					Description: "Number of words to use (default: 2)",
					Required:    false,
					MinValue:    &minWords,
					MaxValue:    maxWords,
				},
			},
		},
		{
			Name:        "randomize",
			Description: "Randomize another user's nickname",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to randomize",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "words",
					Description: "Number of words to use (default: 2)",
					Required:    false,
					MinValue:    &minWords,
					MaxValue:    maxWords,
				},
			},
		},
		{
			Name:        "randomizeall",
			Description: "Randomize nicknames for all eligible members",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "words",
					Description: "Number of words to use (default: 2)",
					Required:    false,
					MinValue:    &minWords,
					MaxValue:    maxWords,
				},
			},
		},
		{
			Name:        "flip",
			Description: "Reverse the word order in a user's nickname",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to flip (defaults to yourself)",
					Required:    false,
				},
			},
		},
		{
			Name:        "rolename",
			Description: "View or set the role required for nickname changes",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "role",
					Description: "The new role name (leave empty to view current)",
					Required:    false,
				},
			},
		},
		{
			Name:        "reloadnames",
			Description: "Reload nickname data from disk",
		},
	}
}

// RegisterCommands registers all slash commands with Discord.
func (h *Handler) RegisterCommands(s *discordgo.Session) error {
	for _, cmd := range h.commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("failed to register command %s: %w", cmd.Name, err)
		}
		h.config.Logger.Info("Registered command", "name", cmd.Name)
	}
	return nil
}

// HandleInteraction routes incoming interactions to the appropriate handler.
func (h *Handler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	guildID := i.GuildID

	switch data.Name {
	case "remember":
		h.handleRemember(s, i, guildID, data.Options)
	case "forget":
		h.handleForget(s, i, guildID, data.Options)
	case "forgetall":
		h.handleForgetAll(s, i, guildID)
	case "names":
		h.handleNames(s, i, guildID)
	case "randomizeme":
		h.handleRandomizeMe(s, i, guildID, data.Options)
	case "randomize":
		h.handleRandomize(s, i, guildID, data.Options)
	case "randomizeall":
		h.handleRandomizeAll(s, i, guildID, data.Options)
	case "flip":
		h.handleFlip(s, i, guildID, data.Options)
	case "rolename":
		h.handleRoleName(s, i, guildID, data.Options)
	case "reloadnames":
		h.handleReloadNames(s, i)
	}
}

func (h *Handler) respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	})
}

func (h *Handler) handleRemember(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	wordsStr := opts[0].StringValue()
	words := strings.Fields(wordsStr)
	if len(words) == 0 {
		h.respond(s, i, "You must provide at least one word to remember.")
		return
	}

	var failed []string
	for _, w := range words {
		if !h.nickManager.Remember(w, guildID) {
			failed = append(failed, w)
		}
	}

	if len(failed) > 0 {
		h.respond(s, i, fmt.Sprintf("Could not remember these words (already in the list): %s", strings.Join(failed, ", ")))
		return
	}
	h.respond(s, i, fmt.Sprintf("Okay! I will remember: %s", strings.Join(words, ", ")))
}

func (h *Handler) handleForget(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	word := opts[0].StringValue()
	h.nickManager.Forget(word, guildID)
	h.respond(s, i, fmt.Sprintf("Okay! I will forget: %s", word))
}

func (h *Handler) handleForgetAll(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string) {
	h.nickManager.ForgetAll(guildID)
	h.respond(s, i, "Okay! I forgot everything. Now I'm useless 🙃")
}

func (h *Handler) handleNames(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string) {
	words := h.nickManager.GetWords(guildID)
	if len(words) == 0 {
		h.respond(s, i, "🤔 I remember... nothing :(")
		return
	}
	h.respond(s, i, fmt.Sprintf("🤔 I remember:\n```%s```", strings.Join(words, " ")))
}

func (h *Handler) handleRandomizeMe(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	if !h.nickManager.HasWords(guildID) {
		h.respond(s, i, "I don't remember any words... 🤔")
		return
	}

	member := i.Member
	if !h.hasRole(s, guildID, member, h.nickManager.GetRoleName(guildID)) {
		h.respond(s, i, fmt.Sprintf("Sorry, you don't have the %s role.", h.nickManager.GetRoleName(guildID)))
		return
	}

	n := 2
	if len(opts) > 0 {
		n = int(opts[0].IntValue())
	}

	name := h.nickManager.GenerateName(guildID, n)
	if name == "" {
		h.respond(s, i, "Not enough words in the pool for that many!")
		return
	}

	if err := s.GuildMemberNickname(guildID, member.User.ID, name); err != nil {
		h.respond(s, i, "Failed to change nickname. I may not have permission.")
		return
	}
	h.respond(s, i, fmt.Sprintf("We will all call %s *%s*!", member.User.Username, name))
}

func (h *Handler) handleRandomize(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	if !h.nickManager.HasWords(guildID) {
		h.respond(s, i, "I don't remember any words... 🤔")
		return
	}

	targetUser := opts[0].UserValue(s)
	targetMember, err := s.GuildMember(guildID, targetUser.ID)
	if err != nil {
		h.respond(s, i, "Could not find that member.")
		return
	}

	if !h.hasRole(s, guildID, targetMember, h.nickManager.GetRoleName(guildID)) {
		h.respond(s, i, fmt.Sprintf("Sorry, %s doesn't have the %s role.", targetUser.Username, h.nickManager.GetRoleName(guildID)))
		return
	}

	n := 2
	if len(opts) > 1 {
		n = int(opts[1].IntValue())
	}

	name := h.nickManager.GenerateName(guildID, n)
	if name == "" {
		h.respond(s, i, "Not enough words in the pool for that many!")
		return
	}

	if err := s.GuildMemberNickname(guildID, targetUser.ID, name); err != nil {
		h.respond(s, i, "Failed to change nickname. I may not have permission.")
		return
	}
	h.respond(s, i, fmt.Sprintf("We will all call %s *%s*!", targetUser.Username, name))
}

func (h *Handler) handleRandomizeAll(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	if !h.nickManager.HasWords(guildID) {
		h.respond(s, i, "I don't remember any words... 🤔")
		return
	}

	n := 2
	if len(opts) > 0 {
		n = int(opts[0].IntValue())
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		h.respond(s, i, "Could not get guild information.")
		return
	}

	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		h.respond(s, i, "Could not get member list.")
		return
	}

	roleName := h.nickManager.GetRoleName(guildID)
	var results []string

	for _, member := range members {
		if member.User.Bot {
			continue
		}
		if member.User.ID == guild.OwnerID {
			continue
		}
		if !h.hasRole(s, guildID, member, roleName) {
			continue
		}

		name := h.nickManager.GenerateName(guildID, n)
		if name == "" {
			continue
		}

		if err := s.GuildMemberNickname(guildID, member.User.ID, name); err != nil {
			results = append(results, fmt.Sprintf("%s: failed", member.User.Username))
			continue
		}
		results = append(results, fmt.Sprintf("%s → *%s*", member.User.Username, name))
	}

	if len(results) == 0 {
		h.respond(s, i, "No eligible members found.")
		return
	}
	h.respond(s, i, fmt.Sprintf("🎉 Party time!\n%s", strings.Join(results, "\n")))
}

func (h *Handler) handleFlip(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	var member *discordgo.Member
	if len(opts) > 0 {
		targetUser := opts[0].UserValue(s)
		m, err := s.GuildMember(guildID, targetUser.ID)
		if err != nil {
			h.respond(s, i, "Could not find that member.")
			return
		}
		member = m
	} else {
		member = i.Member
	}

	if !h.hasRole(s, guildID, member, h.nickManager.GetRoleName(guildID)) {
		h.respond(s, i, "That user doesn't have the required role.")
		return
	}

	oldNick := member.Nick
	if oldNick == "" {
		oldNick = member.User.Username
	}

	words := strings.Fields(oldNick)
	for left, right := 0, len(words)-1; left < right; left, right = left+1, right-1 {
		words[left], words[right] = words[right], words[left]
	}
	newNick := strings.Join(words, " ")

	if err := s.GuildMemberNickname(guildID, member.User.ID, newNick); err != nil {
		h.respond(s, i, "Failed to change nickname. I may not have permission.")
		return
	}
	h.respond(s, i, fmt.Sprintf("Flipped %s from *%s* to *%s*!", member.User.Username, oldNick, newNick))
}

func (h *Handler) handleRoleName(s *discordgo.Session, i *discordgo.InteractionCreate, guildID string, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(opts) == 0 {
		h.respond(s, i, fmt.Sprintf("The current role is: %s", h.nickManager.GetRoleName(guildID)))
		return
	}

	perms := i.Member.Permissions
	if perms&discordgo.PermissionAdministrator == 0 {
		h.respond(s, i, "You must be an administrator to change the role name.")
		return
	}

	newRole := opts[0].StringValue()
	h.nickManager.SetRoleName(newRole, guildID)
	h.respond(s, i, fmt.Sprintf("Role changed to: %s", newRole))
}

func (h *Handler) handleReloadNames(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := h.nickManager.Load(); err != nil {
		h.respond(s, i, "Failed to reload data from disk.")
		return
	}
	h.respond(s, i, "Data reloaded from disk!")
}

func (h *Handler) hasRole(s *discordgo.Session, guildID string, member *discordgo.Member, roleName string) bool {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return false
	}

	roleIDMap := make(map[string]string)
	for _, role := range roles {
		roleIDMap[role.ID] = role.Name
	}

	for _, roleID := range member.Roles {
		if name, exists := roleIDMap[roleID]; exists && name == roleName {
			return true
		}
	}
	return false
}
