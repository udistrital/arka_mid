package models

type PlantillaArchivoResponse struct {
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	File     string `json:"file"`
	Version  string `json:"version"`
}
