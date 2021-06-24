package ui2d

import (
	"bufio"
	"github.com/scottcagno/gng/cmd/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	winWidth  = 1280
	winHeight = 720
)

var renderer *sdl.Renderer
var textureAtlas *sdl.Texture
var textureIndex map[game.Tile][]*sdl.Rect

func loadTextureIndex() {
	fd, err := os.Open("cmd/rpg/ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	textureIndex = make(map[game.Tile][]*sdl.Rect)

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		lines := strings.Split(strings.TrimSpace(line), ",")
		tileRune := game.Tile(lines[0][0])
		x, err := strconv.Atoi(lines[1])
		if err != nil {
			log.Panicf("atoi: %v", err)
		}
		y, err := strconv.Atoi(lines[2])
		if err != nil {
			log.Panicf("atoi: %v", err)
		}
		c, err := strconv.Atoi(lines[3])
		if err != nil {
			log.Panicf("atoi: %v", err)
		}
		var rects []*sdl.Rect
		for i := 0; i < c; i++ {
			rects = append(rects, &sdl.Rect{int32(x * 32), int32(y * 32), 32, 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		textureIndex[tileRune] = rects
	}
}

func imgFileToTexture(filename string) *sdl.Texture {
	fd, err := os.Open(filename)
	if err != nil {
		log.Panicf("opening file: %v", err)
	}
	defer fd.Close()

	img, err := png.Decode(fd)
	if err != nil {
		log.Panicf("png decode: %v", err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	// create texture
	tex, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STATIC,
		int32(w),
		int32(h),
	)
	if err != nil {
		log.Panicf("texture: %v", err)
	}
	tex.Update(nil, pixels, w*4)

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		log.Panicf("setting blend mode: %v", err)
	}

	return tex
}

func init() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)

	// init sdl
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Panicf("init sdl: %v", err)
	}

	// create window
	window, err := sdl.CreateWindow(
		"RPG",
		200,
		200,
		winWidth,
		winHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Panicf("window: %v", err)
	}

	// create renderer
	renderer, err = sdl.CreateRenderer(window,
		-1,
		sdl.RENDERER_ACCELERATED, // USE GPU!
	)
	if err != nil {
		log.Panicf("renderer: %v", err)
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	textureAtlas = imgFileToTexture("cmd/rpg/ui2d/assets/tiles.png")
}

type UI2d struct {
	// TODO: stuff...
}

func (ui *UI2d) Draw(level *game.Level) {
	loadTextureIndex()
	rand.Seed(1)
	for y, row := range level.Map {
		for x, tile := range row {
			if tile == game.Blank {
				continue
			}
			srcRects := textureIndex[tile]
			srcRect := srcRects[rand.Intn(len(srcRects))]
			dstRect := &sdl.Rect{int32(x * 32), int32(y * 32), 32, 32}
			renderer.Copy(textureAtlas, srcRect, dstRect)
		}
	}
	renderer.Present()
	sdl.Delay(5000)
}
