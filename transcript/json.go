package transcript

import (
	"encoding/json"
	"fmt"
)

// JSONFileToJSONFile validates that origin is a PodcastIndex transcript JSON
// document, merges optional metadata, and writes it to destination. It returns
// an *InvalidJSONError (wrapping the decode error, if any) when the input is not
// a valid transcript, and writes nothing in that case.
func JSONFileToJSONFile(originFile, destinationFile string, metadata map[string]string) error {
	text, err := ReadTextRobust(originFile)
	if err != nil {
		return err
	}
	var data any
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return fmt.Errorf("%w: %s: %w", newInvalidJSONError(), originFile, err)
	}

	object, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("%w: %s", newInvalidJSONError(), originFile)
	}
	if _, ok := object["version"].(string); !ok {
		return fmt.Errorf("%w: %s", newInvalidJSONError(), originFile)
	}
	if _, ok := object["segments"].([]any); !ok {
		return fmt.Errorf("%w: %s", newInvalidJSONError(), originFile)
	}

	if len(metadata) > 0 {
		object["metadata"] = metadata
	}
	out, err := encodeJSON(object, jsonIndent)
	if err != nil {
		return err
	}
	return WriteTextUTF8(destinationFile, string(out))
}
