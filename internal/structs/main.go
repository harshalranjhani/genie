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

type Subheading struct {
	LineNum int
	Content string
}

type Heading struct {
	FilePath    string
	LineNum     int
	Content     string
	Subheadings []Subheading
}

type ReplicateMusicResponse struct {
	Status string `json:"status"`
	Logs   string `json:"logs"`
	URLs   struct {
		Get string `json:"get"`
	} `json:"urls"`
	Output string `json:"output,omitempty"`
}

type MailObj struct {
	Email    string    `json:"email"`
	Headings []Heading `json:"headings"`
}

type MailRequest struct {
	MailObj MailObj `json:"mailObj"`
}

type UserStatus struct {
	Email  string `json:"email"`
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}

type ReadmeTemplateResponse struct {
	Template string `json:"template"`
}

// Engine represents an AI engine configuration
type Engine struct {
	Name         string
	Models       []string
	DefaultModel string
	Features     EngineFeatures
}

// EngineFeatures represents supported features for an engine
type EngineFeatures struct {
	SupportsImageGen      bool
	SupportsChat          bool
	SupportsSafeMode      bool
	SupportsReasoning     bool
	SupportsDocumentation bool
}
