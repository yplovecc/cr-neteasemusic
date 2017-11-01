package main

var La = map[string]int{
	"ar": 1, "en": 0, "th": 9, "ja": 3, "pt": 2, "zh": 4,
	"de": 5, "es": 6, "fr": 7, "tr": 8, "id": 10, "ru": 11,
	"ko": 12, "in": 13, "we": 0, "un": 0}

var Alphabets = map[int]string{
	65: "A", 66: "B", 67: "C", 68: "D", 69: "E", 70: "F", 71: "G",
	72: "H", 73: "I", 74: "J", 75: "K", 76: "L", 77: "M", 78: "N",
	79: "O", 80: "P", 81: "Q", 82: "R", 83: "S", 84: "T", 85: "U",
	86: "V", 87: "W", 88: "X", 89: "Y", 90: "Z", 0: "Other"}

var Gender = map[string]int{
	"man": 1, "woman": 2, "combine": 3}

type Crawlinfo struct {
	lo    string
	gen   string
	catid int
}

var Crawlinfos = []Crawlinfo{
	{"zh", "man", 1001}, {"zh", "woman", 1002}, {"zh", "combine", 1003},
	{"ja", "man", 6001}, {"ja", "woman", 6002}, {"ja", "combine", 6003},
	{"ko", "man", 7001}, {"ko", "woman", 7002}, {"ko", "combine", 7003},
	{"we", "man", 2001}, {"we", "woman", 2002}, {"we", "combine", 2003},
	{"un", "man", 4001}, {"un", "woman", 4002}, {"un", "combine", 4003},
}
