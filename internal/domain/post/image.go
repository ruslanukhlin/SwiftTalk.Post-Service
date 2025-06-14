package post

type Image struct {
	UUID string
	URL  string
}

func NewImage(uuid, url string) *Image {
	return &Image{
		UUID: uuid,
		URL:  url,
	}
}
