package string_utils

import "math/rand"

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func Truncate(s string, desiredLength int) string {
    if len(s) <= desiredLength {
        return s
    }

    return s[0:desiredLength]
}

func RandomShuffle(s string, seed int64) string {
    l := len(s)
    bytes := []byte(s)

    rand.Seed(seed)
    rand.Shuffle(l, func(i, j int) { bytes[i], bytes[j] = bytes[j], bytes[i] })

    return string(bytes)
}

