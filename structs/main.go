package structs

type GenResponse struct {
	Candidates []struct {
		Content struct {
			Parts []string `json:"Parts"`
		} `json:"Content"`
	} `json:"Candidates"`
}

type File struct {
	Name string
	Size int64
}

type Directory struct {
	Name     string
	Files    []File
	Children []Directory
}
