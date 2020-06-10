package botcommands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Name() string
	Help() *bothelp.Help
	User() string
	Channel() string
	IsValid() bool
	IsHelpRequest() bool
	AckMsg() (string, bool)
	ErrMsg() string
	APIOptions() map[string]interface{}
}

//
// @evebot 				(show toplevel evebot help/welcome message)
// @evebot help 		(show toplevel help with full command list and help usage)
// @evebot cmd			(show specific command help)
// @evebot cmd help		(show specific command help)
// @evebot help cmd		(show specific sub/command help)
func isHelpRequest(inputCmd []string, cmdName string) bool {
	if len(inputCmd) == 0 || inputCmd[0] == "help" || inputCmd[len(inputCmd)-1] == "help" || (len(inputCmd) == 1 && inputCmd[0] == cmdName) {
		return true
	}
	return false
}

func baseIsValid(input []string) bool {
	if input == nil || len(input) == 0 {
		return false
	}
	return true
}

func baseAckMsg(cmd EvebotCommand, cmdInput []string) (msg string, cont bool) {
	if cmd.IsHelpRequest() {
		return fmt.Sprintf("<@%s>...\n\n%s", cmd.User(), cmd.Help().String()), false
	}
	if cmd.IsValid() == false {
		return fmt.Sprintf("Yo <@%s>, one of us goofed up...¯\\_(ツ)_/¯...I don't know what to do with: `%s`\n\nTry running: ```@evebot %s help```\n\n", cmd.User(), cmdInput, cmd.Name()), false
	}
	if len(cmd.ErrMsg()) > 0 {
		return fmt.Sprintf("Whoops <@%s>! I detected some command *errors:*\n\n ```%v```", cmd.User(), cmd.ErrMsg()), false
	}
	// Happy Path
	return fmt.Sprintf("Sure <@%s>, I'll `%s` that right away. BRB!", cmd.User(), cmd.Name()), true
}

func baseErrMsg(errs []error) string {
	msg := ""
	if len(errs) > 0 {
		for _, v := range errs {
			if len(msg) == 0 {
				msg = v.Error()
			} else {
				msg = msg + "\n" + v.Error()
			}
		}
	}
	return msg
}

func ExtractDatabaseArtifactsOpt(opts map[string]interface{}) eveapi.ArtifactDefinitions {
	log.Logger.Debug("ExtractDatabaseArtifactsOpt", zap.Any("opts", opts))
	if val, ok := opts[botargs.DatabasesName]; ok {
		log.Logger.Debug("ExtractDatabaseArtifactsOpt databases", zap.Any("val", val))
		if artifactDefs, ok := val.(eveapi.ArtifactDefinitions); ok {
			log.Logger.Debug("ExtractDatabaseArtifactsOpt databases artifactDefs", zap.Any("artifactDefs", artifactDefs))
			return artifactDefs
		} else {
			return nil
		}
	}
	return nil
}

func ExtractServiceArtifactsOpt(opts map[string]interface{}) eveapi.ArtifactDefinitions {
	if val, ok := opts[botargs.ServicesName]; ok {
		if artifactDefs, ok := val.(eveapi.ArtifactDefinitions); ok {
			return artifactDefs
		} else {
			return nil
		}
	}
	return nil
}

func ExtractForceDeployOpt(opts map[string]interface{}) bool {
	if val, ok := opts[botargs.ForceDeployName]; ok {
		if forceDepVal, ok := val.(bool); ok {
			return forceDepVal
		} else {
			return false
		}
	}
	return false
}

func ExtractDryrunOpt(opts map[string]interface{}) bool {
	if val, ok := opts[botargs.DryrunName]; ok {
		if dryRunVal, ok := val.(bool); ok {
			return dryRunVal
		} else {
			return false
		}
	}
	return false
}

func ExtractEnvironmentOpt(opts map[string]interface{}) string {
	if val, ok := opts[botparams.EnvironmentName]; ok {
		if envVal, ok := val.(string); ok {
			return envVal
		} else {
			return ""
		}
	}
	return ""
}

func ExtractNSOpt(opts map[string]interface{}) eveapi.StringList {
	if val, ok := opts[botparams.NamespaceName]; ok {
		if nsVal, ok := val.(string); ok {
			return eveapi.StringList{nsVal}
		} else {
			return nil
		}
	}
	return nil
}

type baseCommand struct {
	input               []string
	requiredInputLength int
	name, channel, user string
	valid               bool
	errs                []error
	summary             bothelp.Summary
	usage               bothelp.Usage
	examples            bothelp.Examples
	optionalArgs        botargs.Args
	requiredParams      botparams.Params
	apiOptions          map[string]interface{} // when we resolve the optionalArgs and requiredParams we hydrate this map for fast lookup
}
