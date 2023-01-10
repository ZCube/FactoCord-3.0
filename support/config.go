package support

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/flynn/json5"
)

var FactoCordVersion string

var ConfigPath = "./config.json"

var GuildID string

// Config is a config interface.
var Config configT

type configT struct {
	Executable       string   `json:"executable"`
	LaunchParameters []string `json:"launch_parameters"`
	Autolaunch       bool     `json:"autolaunch"`

	DiscordToken            string `json:"discord_token"`
	GameName                string `json:"game_name"`
	FactorioChannelID       string `json:"factorio_channel_id"`
	Prefix                  string `json:"prefix"`
	HaveServerEssentials    bool   `json:"have_server_essentials"`
	IngameDiscordUserColors bool   `json:"ingame_discord_user_colors"`

	AllowPingingEveryone bool `json:"allow_pinging_everyone"`

	EnableConsoleChannel  bool   `json:"enable_console_channel"`
	FactorioConsoleChatID string `json:"factorio_console_chat_id"`

	AdminIDs     []string          `json:"admin_ids"`
	CommandRoles map[string]string `json:"command_roles"`

	ModListLocation string `json:"mod_list_location"`
	Username        string `json:"username"`
	ModPortalToken  string `json:"mod_portal_token"`

	Messages struct {
		BotStart          string `json:"bot_start"`
		BotStop           string `json:"bot_stop"`
		ServerStart       string `json:"server_start"`
		ServerStop        string `json:"server_stop"`
		ServerFail        string `json:"server_fail"`
		ServerSave        string `json:"server_save"`
		PlayerJoin        string `json:"player_join"`
		PlayerLeave       string `json:"player_leave"`
		DownloadStart     string `json:"download_start"`
		DownloadProgress  string `json:"download_progress"`
		DownloadComplete  string `json:"download_complete"`
		Unpacking         string `json:"unpacking"`
		UnpackingComplete string `json:"unpacking_complete"`
	} `json:"messages"`
}

func (conf *configT) MustLoad() {
	if !FileExists(ConfigPath) {
		fmt.Println("Error: config.json not found.")
		fmt.Println("Make sure that you copied 'config-example.json' and current working directory is correct")
		Exit(7)
	}
	contents, err := ioutil.ReadFile(ConfigPath)
	Critical(err, "... when reading config.json")

	conf.defaults()
	err = json5.Unmarshal(contents, &conf)
	if err != nil {
		Critical(err, "... when parsing config.json")
	}

	if len(conf.Executable) == 0 {
		step := -1
		args := []string{}
		for _, arg := range os.Args {
			if arg == "--" {
				step = 0
			}
			if step == 1 {
				conf.Executable = arg
			} else if step > 1 {
				args = append(args, arg)
			}
			if step >= 0 {
				step++
			}
		}
		if step > 1 {
			conf.LaunchParameters = args
		}
	}
	if len(conf.DiscordToken) == 0 {
		conf.DiscordToken = os.Getenv("DISCORD_TOKEN")
	}
	if len(conf.FactorioChannelID) == 0 {
		conf.FactorioChannelID = os.Getenv("FACTORIO_CHANNEL_ID")
	}
	if len(conf.Username) == 0 {
		conf.Username = os.Getenv("USERNAME")
	}
	if len(conf.ModPortalToken) == 0 {
		conf.ModPortalToken = os.Getenv("TOKEN")
	}

	if len(conf.DiscordToken) == 0 {
		if data, err := os.ReadFile("/discord/token"); err == nil {
			str := strings.TrimSpace(string(data))
			if len(str) > 0 {
				conf.DiscordToken = str
			}
		}
	}
	if len(conf.FactorioChannelID) == 0 {
		if data, err := os.ReadFile("/discord/factorio_channel_id"); err == nil {
			str := strings.TrimSpace(string(data))
			if len(str) > 0 {
				conf.FactorioChannelID = str
			}
		}
	}
	if len(conf.Username) == 0 {
		if data, err := os.ReadFile("/account/username"); err == nil {
			str := strings.TrimSpace(string(data))
			if len(str) > 0 {
				conf.Username = str
			}
		}
	}
	if len(conf.ModPortalToken) == 0 {
		if data, err := os.ReadFile("/account/token"); err == nil {
			str := strings.TrimSpace(string(data))
			if len(str) > 0 {
				conf.ModPortalToken = str
			}
		}
	}
}

func (conf *configT) Load() error {
	if !FileExists(ConfigPath) {
		return fmt.Errorf("config.json not found")
	}
	contents, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return fmt.Errorf("error reading config.json: %s", err)
	}

	test := configT{}
	err = json5.Unmarshal(contents, &test)
	if err != nil {
		return fmt.Errorf("error parsing config.json: %s", err)
	}

	conf.defaults()
	err = json5.Unmarshal(contents, &conf)
	Critical(err, "wtf?? error parsing config.json 2nd time")
	return nil
}

func (conf *configT) defaults() {
	conf.Autolaunch = true
	conf.GameName = "Factorio"
	conf.Prefix = "$"
	// conf.HaveServerEssentials = false
	// conf.IngameDiscordUserColors = false
	conf.Messages.BotStart = "**:white_check_mark: Bot started! Launching server...**"
	conf.Messages.BotStop = ":v:"
	conf.Messages.ServerStart = "**:white_check_mark: The server has started!**"
	conf.Messages.ServerStop = "**:octagonal_sign: The server has stopped!**"
	conf.Messages.ServerFail = "**:skull: The server has crashed!**"
	conf.Messages.ServerSave = "**:floppy_disk: Game saved!**"
	conf.Messages.PlayerJoin = "**:arrow_up: {username}**"
	conf.Messages.PlayerLeave = "**:arrow_down: {username}s**"
	conf.Messages.DownloadStart = ":arrow_down: Downloading {file}..."
	conf.Messages.DownloadProgress = ":arrow_down: Downloading {file}: {percent}%"
	conf.Messages.DownloadComplete = ":white_check_mark: Downloaded {file}"
	conf.Messages.Unpacking = ":pinching_hand: Unpacking {file}..."
	conf.Messages.UnpackingComplete = ":ok_hand: Server updated to {version}"
}
