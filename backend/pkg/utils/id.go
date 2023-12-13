package utils

import (
	"fmt"

	"github.com/904566722/myutils/strings"
)

const (
	idRandStrLen = 8
	spiltStr     = "-"
)

func GenId(prefix string) string {
	return fmt.Sprintf("%s%s%s", prefix, spiltStr, strings.RandStr(idRandStrLen))
}
