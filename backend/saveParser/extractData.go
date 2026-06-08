package saveparser

import (
	"log"

	"github.com/tidwall/gjson"
)

// saveData is a json string
func GetRunHistoryEntries(saveData string) []string {
	runs := gjson.Get(saveData, "root.properties.RunHistory_0.Entries_0")
	if runs.Exists() {
		log.Println("Found run history")
		retval := make([]string, 0)
		runs.ForEach(func(key, value gjson.Result) bool {
			retval = append(retval, value.String())
			return true
		})
		return retval
	} else {
		return nil
	}
}
