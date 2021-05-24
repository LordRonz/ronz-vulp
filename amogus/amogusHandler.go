package amogus

import (
	"strings"
)

func RemoveDups(s string) string {
    var buf strings.Builder
    var last rune
    for i, r := range s {
        if r != last || i == 0 {
            buf.WriteRune(r)
            last = r
        }
    }
    return buf.String()
}