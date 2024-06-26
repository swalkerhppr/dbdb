package assets

type MetadataMap struct {
	ImageFiles   []FileInfo `xml:"ImageFile"`
	AudioFiles   []FileInfo `xml:"AudioFile"`
}

type FileInfo struct {
	Filename     string     `xml:"filename,attr"`
	FileMetadata []Metadata `xml:"Metadata"`
	Volume       float64    `xml:"volume,attr"`
}

type Metadata struct {
	// Name in the map it will be called
	Name         string      `xml:"name,attr"`
	// Dimensions of sprite cells
	Height       int         `xml:"height,attr"`
	Width        int         `xml:"width,attr"`
	// Starting coordinates
	Top          int         `xml:"top,attr"`
	Left         int         `xml:"left,attr"`
	Frames       []FrameSpec `xml:"Frames"`
}

type FrameSpec struct {
	// Arguments to the grid.Frames function. See https://github.com/yohamta/ganim8 for reference.
	// Col (AKA x) is first
	ColNum   int    `xml:"col-num,attr"`
	ColRange string `xml:"col-range,attr"`
	// Row (AKA y) is second
	RowNum   int    `xml:"row-num,attr"`
	RowRange string `xml:"row-range,attr"`
}
