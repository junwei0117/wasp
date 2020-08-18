package fairroulette

import (
	"github.com/iotaledger/wasp/tools/wasp-client/config"
	"github.com/spf13/pflag"
)

var scConfig *config.SCConfig

func HookFlags() *pflag.FlagSet {
	scConfig = config.NewSC("fairroulette", "fr")
	return scConfig.Flags
}

var commands = map[string]func([]string){
	"set":    scConfig.HandleSetCmd,
	"admin":  adminCmd,
	"status": statusCmd,
	"bet":    betCmd,
}

func Cmd(args []string) {
	scConfig.HandleCmd(args, commands)
}
