package storage

const (
	StorePhoto = "Storage.StorePhoto"
)

type StorePhotoArgs struct {
	UserId   int64  `json:"user-id"`
	FileType string `json:"file-type"`
	FileSize int64  `json:"file-size"`
	Photo    []byte `json:"photo"`
}

type StorePhotoReply struct {
	Filekey  string
	Url      string
	Suffixes []string
	Width    int
	Height   int
}
