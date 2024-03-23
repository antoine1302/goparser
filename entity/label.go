package entity

type Label struct {
	Id          int      `xml:"id"`
	Name        string   `xml:"name"`
	ContactInfo string   `xml:"contactinfo"`
	Profile     string   `xml:"profile"`
	Urls        []string `xml:"urls>url"`
}
