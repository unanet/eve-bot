package commands

import (
	"regexp"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
)

func ExtractServiceArtifactsOpt(opts CommandOptions) eveapimodels.ArtifactDefinitions {
	if val, ok := opts[args.ServicesName]; ok {
		if artifactDefs, ok := val.(eveapimodels.ArtifactDefinitions); ok {
			return artifactDefs
		}
		return nil

	}
	return nil
}

func ExtractDatabaseArtifactsOpt(opts CommandOptions) eveapimodels.ArtifactDefinitions {
	if val, ok := opts[args.DatabasesName]; ok {
		if artifactDefs, ok := val.(eveapimodels.ArtifactDefinitions); ok {
			return artifactDefs
		}
		return nil

	}
	return nil
}

func ExtractForceDeployOpt(opts CommandOptions) bool {
	if val, ok := opts[args.ForceDeployName]; ok {
		if forceDepVal, ok := val.(bool); ok {
			return forceDepVal
		}
		return false
	}
	return false
}

func ExtractDryrunOpt(opts CommandOptions) bool {
	if val, ok := opts[args.DryrunName]; ok {
		if dryRunVal, ok := val.(bool); ok {
			return dryRunVal
		}
		return false
	}
	return false
}

func ExtractEnvironmentOpt(opts CommandOptions) string {
	if val, ok := opts[params.EnvironmentName]; ok {
		if envVal, ok := val.(string); ok {
			return envVal
		}
		return ""
	}
	return ""
}

func ExtractNSOpt(opts CommandOptions) eveapimodels.StringList {
	if val, ok := opts[params.NamespaceName]; ok {
		if nsVal, ok := val.(string); ok {
			return eveapimodels.StringList{nsVal}
		}
		return nil
	}
	return nil
}

func CleanUrls(input string) string {
	matcher := regexp.MustCompile(`<[a-zA-Z]+://[a-zA-Z._\-:\d/|]+>`)
	matchIndexes := matcher.FindAllStringIndex(input, -1)
	matchCount := len(matchIndexes)

	if matchCount == 0 {
		return input
	}

	firstMatchIndex := matchIndexes[0][0]
	lastMatchIndex := matchIndexes[matchCount-1][1]

	firstPart := input[0:firstMatchIndex]
	lastPart := input[lastMatchIndex:]

	cleanPart := firstPart
	for i, v := range matchIndexes {
		if i > 0 {
			previousMatchLastIndex := matchIndexes[i-1][1]
			currentMatchFirstIndex := matchIndexes[i][0]
			middleMatch := input[previousMatchLastIndex:currentMatchFirstIndex]
			cleanPart = cleanPart + middleMatch
		}

		matchedVal := input[v[0]:v[1]]
		cleanVal := ""

		if strings.Contains(matchedVal, "|") {
			vals := strings.Split(matchedVal, "|")
			cleanVal = vals[1][:len(vals[1])-len(">")]
		} else {
			cleanVal = strings.ReplaceAll(matchedVal, "<", "")
			cleanVal = strings.ReplaceAll(cleanVal, ">", "")
		}

		cleanPart = cleanPart + cleanVal
	}
	return cleanPart + lastPart
}

func nonHelpCmd() []EvebotCommand {
	var cmds []EvebotCommand

	for k, v := range CommandInitializerMap {
		if k != "help" {
			cmds = append(cmds, v.(func([]string, string, string) EvebotCommand)([]string{}, "", ""))
		}
	}
	return cmds
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

func ackMsg(msg string, valid bool) ChatAckMsgFn {
	return func() (string, bool) {
		return msg, valid
	}
}

func hydrateMetadataMap(keyvals []string) params.MetadataMap {
	result := make(params.MetadataMap, 0)
	if len(keyvals) == 0 {
		return nil
	}
	for _, s := range keyvals {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			key := CleanUrls(argKV[0])
			value := CleanUrls(argKV[1])
			result[key] = value
		}
	}
	return result
}
