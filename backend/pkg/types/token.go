package types

// TokenInput models a token input.
type TokenInput struct {
	Tokenname                string `json:"tokenname"`
	Bearer                   string `json:"bearer"`
	Displayname              string `json:"displayname"`
	Description              string `json:"description"`
	FileFromIPFS             string `json:"fileFromIPFS"`
	FileFromBase64           string `json:"fileFromBase64"`
	MetadataPlaceholderName  string `json:"metadataPlaceholderName"`
	MetadataPlaceholderValue string `json:"metadataPlaceholderValue"`
}
