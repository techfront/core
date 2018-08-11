package media

type Size struct {
	Width  int64
	Height int64
	Crop   bool
}

type Media struct {
	Sizes         []Size
	ThumbnailsFormat string
	ThumbnailsSize   Size
}

var MediaConfig *Media

func Setup() {
	m := new(Media)

	m.ThumbnailsFormat = ".jpg"
	m.ThumbnailsSize = Size{Width: 800, Height: 0, Crop: false}

	m.Sizes = append(m.Sizes, Size{Width: 50, Height: 0, Crop: false})
	m.Sizes = append(m.Sizes, Size{Width: 100, Height: 0, Crop: false})
	m.Sizes = append(m.Sizes, Size{Width: 500, Height: 0, Crop: false})
	m.Sizes = append(m.Sizes, Size{Width: 1000, Height: 0, Crop: false})

	MediaConfig = m
}