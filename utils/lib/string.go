package lib

import "strconv"

func Str2Int(s string) (int, error) {
    return strconv.Atoi(s)
}

func Str2Int64(s string) (int64, error) {
    return strconv.ParseInt(s, 10, 64)
}

func Int2Str(i int) string {
    return strconv.Itoa(i)
}

func Int642Str(i int64) string {
    return strconv.FormatInt(i, 10)
}