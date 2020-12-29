package style

import "fmt"

func QuestionElement(str string, a ...interface{}) string {
	// return color.HiBlueString(str, a...)
	return fmt.Sprintf(str, a...)
}
