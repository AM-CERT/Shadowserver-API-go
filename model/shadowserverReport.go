package model

// ShadowserverReport struct
type ShadowserverReport struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	File      string `json:"file"`
	Report    string `json:"report"`
	Url       string `json:"url"`
	Timestamp string `json:"timestamp"`
	FilePath  string `json:"file_path"`
}
