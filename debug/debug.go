package debug

import (
	"encoding/json"
	"log"
)

// PrintJSON prints message followed by the json representation of the passed object.
// If any error occurs, it just prints a debug message.
func PrintJSON(message string, v interface{}) {
	js, err := json.MarshalIndent(v, "", " ")

	if err != nil {
		log.Printf("Can't marshal the object to json: $s", err.Error())
		return
	}

	log.Printf("%s %s", message, js)
}
