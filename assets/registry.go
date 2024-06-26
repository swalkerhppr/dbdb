package assets

import (
	"embed"
	"encoding/xml"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
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
		soundMap  : make(map[string]*audio.Player),
		audioCtx  : audio.NewContext(44100),
	}

	Registry.loadImages(metaMap.ImageFiles, assetFS)

	Registry.fontLib.ParseEmbedDirFonts("game_assets/fonts", assetFS)
	log.Printf("Loaded Fonts: %+v", Registry.fontLib)

	Registry.loadSounds(metaMap.AudioFiles, assetFS)
}

var Registry *registry

type registry struct {
	imageMap  map[string]*ebiten.Image
	spriteMap map[string]*ganim8.Sprite
	fontLib   *etxt.FontLibrary
	soundMap  map[string]*audio.Player
	audioCtx  *audio.Context
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

func (r *registry) Sound(name string) *audio.Player {
	player, exist := r.soundMap[name]
	if !exist {
		log.Printf("Could not find sound %s!", name)
		return r.soundMap["not-found.ogg"]
	}
	player.SetPosition(0)
	return player
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

// Gets the default font renderer in the given size
func (r *registry) DefaultTextRenderer(size int) *etxt.Renderer {
	return r.TextRenderer("Kode Mono Regular", size - 3, color.RGBA{R: 56, G: 44, B: 53, A: 255} )
}
func (r *registry) BoldTextRenderer(size int) *etxt.Renderer {
	return r.TextRenderer("Kode Mono Bold", size - 3, color.RGBA{R: 56, G: 44, B: 53, A: 255} )
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

func (r *registry) loadSounds(fileinfos []FileInfo, fs embed.FS) {
	for _, fi := range fileinfos {
		f, err := fs.Open("game_assets/" + fi.Filename)
		if err != nil {
			log.Printf("Could not read sound file %s", fi.Filename)
			continue
		}
		stream, err := vorbis.DecodeWithoutResampling(f)
		if err != nil {
			log.Printf("Could not decode sound file %s", fi.Filename)
			f.Close()
			continue
		}
		player, err := r.audioCtx.NewPlayer(stream)
		if err != nil {
			log.Printf("Could not create player for sound file %s", fi.Filename)
			f.Close()
			continue
		}
		if fi.Volume > 0 {
			player.SetVolume(fi.Volume)
		}
		r.soundMap[fi.Filename] = player
		log.Printf("Loaded sound file: %s", fi.Filename)
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
