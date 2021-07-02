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

type ui struct {
	winWidth          int //= 1280
	winHeight         int //= 720
	renderer          *sdl.Renderer
	window            *sdl.Window
	textureAtlas      *sdl.Texture
	textureIndex      map[game.Tile][]*sdl.Rect
	prevKeyboardState []uint8
	keyboardState     []uint8
	centerX           int
	centerY           int
	r                 *rand.Rand
	levelChan         chan *game.Level
	inputChan         chan *game.Input
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {

	winWidth, winHeight := int32(1280), int32(720)

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
	renderer, err := sdl.CreateRenderer(window,
		-1,
		sdl.RENDERER_ACCELERATED, // USE GPU!
	)
	if err != nil {
		log.Panicf("renderer: %v", err)
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	keyboardState := sdl.GetKeyboardState()
	prevKeyboardState := make([]uint8, len(keyboardState))
	for i, v := range keyboardState {
		prevKeyboardState[i] = v
	}

	ui := &ui{
		winWidth:          int(winWidth),
		winHeight:         int(winHeight),
		renderer:          renderer,
		window:            window,
		prevKeyboardState: prevKeyboardState,
		keyboardState:     keyboardState,
		centerX:           -1,
		centerY:           -1,
		r:                 rand.New(rand.NewSource(1)),
		levelChan:         levelChan,
		inputChan:         inputChan,
	}
	ui.textureAtlas = ui.imgFileToTexture("cmd/rpg/ui2d/assets/tiles.png")
	ui.loadTextureIndex()

	return ui
}

func (ui *ui) loadTextureIndex() {
	fd, err := os.Open("cmd/rpg/ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	ui.textureIndex = make(map[game.Tile][]*sdl.Rect)

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
		ui.textureIndex[tileRune] = rects
	}
}

func (ui *ui) imgFileToTexture(filename string) *sdl.Texture {
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
	tex, err := ui.renderer.CreateTexture(
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
	// init sdl
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Panicf("init sdl: %v", err)
	}
}

func (ui *ui) Draw(level *game.Level) {
	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	}

	limit := 5
	if level.Player.X > ui.centerX+limit {
		ui.centerX++
	} else if level.Player.X < ui.centerX-limit {
		ui.centerX--
	} else if level.Player.Y > ui.centerY+limit {
		ui.centerY++
	} else if level.Player.Y < ui.centerY+limit {
		ui.centerY--
	}

	offsetX := int32(ui.winWidth/2 - ui.centerX*32)
	offsetY := int32(ui.winHeight/2 - ui.centerY*32)

	ui.renderer.Clear()
	ui.r.Seed(1)

	// draw tiles
	for y, row := range level.Map {
		for x, tile := range row {
			if tile == game.Blank {
				continue
			}
			srcRects := ui.textureIndex[tile]
			srcRect := srcRects[ui.r.Intn(len(srcRects))]
			dstRect := &sdl.Rect{int32(x*32) + offsetX, int32(y*32) + offsetY, 32, 32}

			pos := game.Pos{x, y}
			if _, ok := level.Debug[pos]; ok {
				ui.textureAtlas.SetColorMod(128, 0, 0)
			} else {
				ui.textureAtlas.SetColorMod(255, 255, 255)
			}

			ui.renderer.Copy(ui.textureAtlas, srcRect, dstRect)
		}
	}

	// draw monsters
	for pos, monster := range level.Monsters {
		monsterSrcRect := ui.textureIndex[game.Tile(monster.Rune)][0]
		ui.renderer.Copy(ui.textureAtlas,
			//&sdl.Rect{int32(21 * 32), int32(59 * 32), 32, 32},
			monsterSrcRect,
			&sdl.Rect{int32(pos.X)*32 + offsetX, int32(pos.Y)*32 + offsetY, 32, 32},
		)
	}

	// draw player
	playerSrcRect := ui.textureIndex['@'][0]
	ui.renderer.Copy(ui.textureAtlas,
		//&sdl.Rect{int32(21 * 32), int32(59 * 32), 32, 32},
		playerSrcRect,
		&sdl.Rect{int32(level.Player.X)*32 + offsetX, int32(level.Player.Y)*32 + offsetY, 32, 32},
	)

	// present
	ui.renderer.Present()
}

func (ui *ui) keyPress(scancode uint8) bool {
	return ui.keyboardState[scancode] == 0 && ui.prevKeyboardState[scancode] != 0
}

func (ui *ui) Run() {

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChan <- &game.Input{
					Typ: game.QuitGame,
				}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{
						Typ:          game.CloseWindow,
						LevelChannel: ui.levelChan,
					}
				}
			}
		}

		// check to see if we need to draw
		select {
		case newLevel, ok := <-ui.levelChan:
			if ok {
				ui.Draw(newLevel)
			}
		default:
		}

		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
			var input game.Input
			if ui.keyPress(sdl.SCANCODE_UP) {
				input.Typ = game.Up
			}
			if ui.keyPress(sdl.SCANCODE_DOWN) {
				input.Typ = game.Down
			}
			if ui.keyPress(sdl.SCANCODE_LEFT) {
				input.Typ = game.Left
			}
			if ui.keyPress(sdl.SCANCODE_RIGHT) {
				input.Typ = game.Right
			}
			if ui.keyPress(sdl.SCANCODE_S) {
				input.Typ = game.Search
			}
			if ui.keyPress(sdl.SCANCODE_Q) {
				input.Typ = game.QuitGame
			}

			for i, v := range ui.keyboardState {
				ui.prevKeyboardState[i] = v
			}

			if input.Typ != game.None {
				ui.inputChan <- &input
			}
		}
	}
}
