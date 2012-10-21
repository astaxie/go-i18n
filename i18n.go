package i18n

import (
	"fmt"
)

var DefaultLocale string

func init() {
	DefaultLocale = "zh-CN"
}

func SetLocale(lc string) {
	DefaultLocale = lc
}

/**
* Translation
**/

func T(key string) string {

}
