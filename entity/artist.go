package entity

type Artist struct {
	Id             int      `xml:"id"`
	Name           string   `xml:"name"`
	Realname       string   `xml:"realname"`
	Profile        string   `xml:"profile"`
	Urls           []string `xml:"urls>url"`
	NameVariations []string `xml:"namevariations>name"`
}

// type Artist struct {
// 	Id             int             `xml:"id"`
// 	Name           string          `xml:"name"`
// 	Realname       string          `xml:"realname"`
// 	Profile        string          `xml:"profile"`
// 	Urls           []string        `xml:"urls>url"`
// 	NameVariations []string        `xml:"namevariations>name"`
// 	Aliases        []ArtistPartial `xml:"aliases>name"`
// 	Members        []ArtistPartial `xml:"members>name"`
// }

// type ArtistPartial struct {
// 	Id   int    `xml:"id,attr"`
// 	Name string `xml:",chardata"`
// }
