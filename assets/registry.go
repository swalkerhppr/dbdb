package assets

import (
	"embed"
	"encoding/xml"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tinne26/etxt"
	"github.com/yohamta/ganim8/v2"
)

func InitRegistry() {
	metaMap := MetadataMap{}
	if err := xml.Unmarshal(metadata, &metaMap); err != nil {
		log.Printf("Could not initialize assets! %v", err)
		return
	}
	Registry = &registry{
		imageMap  : make(map[string]*ebiten.Image),
		spriteMap : make(map[string]*ganim8.Sprite),
		fontLib   : etxt.NewFontLibrary(),
	}

	Registry.loadImages(metaMap.ImageFiles, assetFS)

	Registry.fontLib.ParseEmbedDirFonts("game_assets/fonts", assetFS)
	log.Printf("Loaded Fonts: %+v", Registry.fontLib)
}

var Registry *registry

type registry struct {
	imageMap  map[string]*ebiten.Image
	spriteMap map[string]*ganim8.Sprite
	fontLib   *etxt.FontLibrary
}

func (r *registry) Image(name string) *ebiten.Image {
	if img, ok := r.imageMap[name]; ok {
		return img
	}
	panic("Could not find image: " + name)
}

func (r *registry) Sprite(name string) *ganim8.Sprite {
	if sprite, ok := r.spriteMap[name]; ok {
		return sprite.Clone()
	}
	panic("Could not find sprite: " + name)
}

func (r *registry) UnsafeSprite(name string) *ganim8.Sprite {
	return r.spriteMap[name]
}

// Gets the given font renderer
func (r *registry) TextRenderer(name string, size int, c color.Color) *etxt.Renderer {
	if ! r.fontLib.HasFont(name) {
		panic("Requested font that does not exist: " + name)
	}
	render := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	render.SetCacheHandler(glyphsCache.NewHandler())
	render.SetFont(r.fontLib.GetFont(name))
	render.SetColor(c)
	render.SetSizePx(size)
	return render
}

// Gets the default font renderer (Pixelify Sans Regular) in the given size
func (r *registry) DefaultTextRenderer(size int) *etxt.Renderer {
	return r.TextRenderer("Pixelify Sans Regular", size, color.RGBA{R: 56, G: 44, B: 53, A: 255} )
}

func (r *registry) loadImages(fileinfos []FileInfo, fs embed.FS) {

	for _, fi := range fileinfos {
		img, _, err := ebitenutil.NewImageFromFileSystem(fs, "game_assets/" + fi.Filename)
		if err != nil {
			log.Printf("Could not read image file: %s. %v", fi.Filename, err)
			continue
		}
		r.imageMap[fi.Filename] = img
		log.Printf("Registered Image %s.", fi.Filename)
		for _, meta := range fi.FileMetadata {
			grid := ganim8.NewGrid(meta.Width, meta.Height, img.Bounds().Dx(), img.Bounds().Dy(), meta.Left, meta.Top)
			sprite := ganim8.NewSprite(img, grid.Frames(convertFrames(meta.Frames)...))
			r.spriteMap[meta.Name] = sprite
			log.Printf("Registered Sprite %s.", meta.Name)
		}
	}
}

func convertFrames(specs []FrameSpec) []interface{} {
	res := make([]interface{}, len(specs) * 2)
	for i, s := range specs {
		if s.ColRange != "" {
			res[i * 2] = s.ColRange
		} else {
			res[i * 2] = s.ColNum
		}
		if s.RowRange != "" {
			res[(i*2)+1] = s.RowRange
		} else {
			res[(i*2)+1] = s.RowNum
		}
	}
	return res
}
