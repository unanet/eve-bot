package commands

import (
	"regexp"
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve/pkg/eve"
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

func ExtractMetadataField(opts CommandOptions) eve.MetadataField {
	if metadataMap, metaDataOK := opts[params.MetadataName].(params.MetadataMap); metaDataOK {
		return metadataMap.ToMetadataField()
	}
	return nil
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
		return cleanEncoding(input)
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
	return cleanEncoding(result)
}

func hydrateMetadataMap(keyvals []string) params.MetadataMap {
	result := make(params.MetadataMap)
	if len(keyvals) == 0 {
		return nil
	}

	for i, s := range keyvals {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")

			result[argKV[0]] = extractMetadataValue(argKV[1], keyvals[i+1:])
		}
	}

	return result
}

func extractMetadataValue(initialValue string, keyvals []string) string {
	var builder strings.Builder

	// Let's not forget our initial value for a properly formatted metadata value of `FOO=BAR`.
	// If this has a space like `FOO= BAR`, we will trim the spaces out later
	builder.WriteString(initialValue + " ")

	for _, s := range keyvals {
		// While we iterate over our values, we want to stop if we see a `=`. This will indicate the start of a new key
		if strings.Contains(s, "=") {
			break
		}

		builder.WriteString(s + " ")
	}
	return strings.TrimSpace(builder.String())
}