package commands

import (
	"regexp"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

// ExtractArtifactsDefinition extracts the ArtifactDefinitions from the CommandOptions
func ExtractArtifactsDefinition(defType string, opts CommandOptions) eve.ArtifactDefinitions {
	if val, ok := opts[defType]; ok {
		if artifactDefs, ok := val.(eve.ArtifactDefinitions); ok {
			return artifactDefs
		}
		return nil

	}
	return nil
}

// ExtractBoolOpt extracts a bool key/val from the opts
func ExtractBoolOpt(defType string, opts CommandOptions) bool {
	if val, ok := opts[defType]; ok {
		if forceDepVal, ok := val.(bool); ok {
			return forceDepVal
		}
		return false
	}
	return false
}

// ExtractStringOpt extracts a string key/val from the options
func ExtractStringOpt(defType string, opts CommandOptions) string {
	if val, ok := opts[defType]; ok {
		if envVal, ok := val.(string); ok {
			return envVal
		}
		return ""
	}
	return ""
}

// ExtractStringListOpt extracts a string slice key value from the options
func ExtractStringListOpt(defType string, opts CommandOptions) eve.StringList {
	if val, ok := opts[defType]; ok {
		if nsVal, ok := val.(string); ok {
			return eve.StringList{nsVal}
		}
		return nil
	}
	return nil
}

func cleanEncoding(input string) string {
	input = strings.ReplaceAll(input, "&lt;", "<")
	input = strings.ReplaceAll(input, "&gt;", ">")
	return input
}

// CleanUrls cleans the incoming URLs
// this iterates the incoming command and removes any encoding slack adds to URLs
func CleanUrls(input string) string {
	matcher := regexp.MustCompile(`<([^>]*)>`)
	matchIndexes := matcher.FindAllStringIndex(input, -1)
	matchCount := len(matchIndexes)

	if matchCount == 0 {
		log.Logger.Info("CleanUrls no match", zap.String("input", input))
		cleanedEnc := cleanEncoding(input)
		log.Logger.Info("CleanUrls no match cleaned", zap.String("cleanedEnc", cleanedEnc))
		return cleanedEnc
	}

	cleanPart := input[0:matchIndexes[0][0]]
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

		cleanPart += cleanVal
	}
	result := cleanPart + input[matchIndexes[matchCount-1][1]:]
	log.Logger.Info("CleanUrls result", zap.String("result", result))
	cleanedEnc := cleanEncoding(result)
	log.Logger.Info("CleanUrls result cleaned", zap.String("cleanedEnc", cleanedEnc))
	return cleanedEnc
}

func hydrateMetadataMap(keyvals []string) params.MetadataMap {
	log.Logger.Info("TROY keyvals", zap.Strings("key vals", keyvals))
	result := make(params.MetadataMap, 0)
	if len(keyvals) == 0 {
		return nil
	}
	for _, s := range keyvals {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			//result[CleanUrls(argKV[0])] = CleanUrls(strings.Join(argKV[1:], "="))
			result[argKV[0]] = strings.Join(argKV[1:], "=")
		}
	}

	log.Logger.Info("TROY keyvals result", zap.Any("result", result))
	return result
}
