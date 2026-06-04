package domain

type KategoriPayload struct {
	UUID     string  `json:"uuid"`
	UUIDSub  *string `json:"uuidSub"`
	FullText string  `json:"full_text"`
}
