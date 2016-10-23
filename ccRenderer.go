/* TODO:

* KEY-BASED NAVIGATION (CTRL-HOME/END - PGUP/DN)
* ??DONE?? BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
* HORIZONTAL SCROLLBARS
	horizontal could be charHei thickness
	vertical could easily be a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar
*/

package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

var rend = CcRenderer{}
var black = []float32{0, 0, 0, 1}
var blue = []float32{0, 0, 1, 1}
var cyan = []float32{0, 0.5, 1, 1}
var fuschia = []float32{0.6, 0.2, 0.3, 1}
var gray = []float32{0.25, 0.25, 0.25, 1}
var grayDark = []float32{0.15, 0.15, 0.15, 1}
var grayLight = []float32{0.4, 0.4, 0.4, 1}
var green = []float32{0, 1, 0, 1}
var magenta = []float32{1, 0, 1, 1}
var orange = []float32{0.8, 0.35, 0, 1}
var purple = []float32{0.6, 0, 0.8, 1}
var red = []float32{1, 0, 0, 1}
var tan = []float32{0.55, 0.47, 0.37, 1}
var violet = []float32{0.4, 0.2, 1, 1}
var white = []float32{1, 1, 1, 1}
var yellow = []float32{1, 1, 0, 1}

// dimensions (in pixel units)
var prevAppWidth int = 800
var prevAppHeight int = 600
var currAppWidth int = 800
var currAppHeight int = 600

type CcRenderer struct {
	DistanceFromOrigin float32
	RunPanelHeiPerc    float32 // FIXME: hardwired value for a specific use case
	ClientExtentX      float32 // distance from the center to an edge of the app's root/client area
	ClientExtentY      float32
	// ....in the cardinal directions from the center, corners would be farther away)
	PixelWid        float32
	PixelHei        float32
	CharWid         float32
	CharHei         float32
	CharWidInPixels int
	CharHeiInPixels int
	UvSpan          float32 // looking into 16/16 atlas/grid of character tiles
	// FIXME: below is no longer a maximum of what fits on a max-sized panel (taking up the whole app window) anymore.
	// 		but is still used as a guide for sizes
	MaxCharsX int // this is used to give us proportions like an 80x25 text console screen, ....
	MaxCharsY int // ....from a cr.DistanceFromOrigin*2-by-cr.DistanceFromOrigin*2 gl space
	// current position renderer draws to
	CurrX     float32
	CurrY     float32
	PrevColor []float32 // previous
	CurrColor []float32
	Focused   *TextPanel
	Panels    []*TextPanel
}

func (cr *CcRenderer) Init() {
	if cr.ClientExtentX == 0.0 || cr.ClientExtentY == 0.0 {
		fmt.Println("CcRenderer.Init(): FIRST TIME")
		cr.DistanceFromOrigin = 3
		cr.UvSpan = float32(1.0) / 16 // how much uv a pixel spans
		cr.RunPanelHeiPerc = 0.4
		cr.PrevColor = grayDark
		cr.CurrColor = grayDark
		cr.ClientExtentX = cr.DistanceFromOrigin
		cr.ClientExtentY = cr.DistanceFromOrigin
		cr.PixelWid = cr.ClientExtentX * 2 / float32(currAppWidth)
		cr.PixelHei = cr.ClientExtentY * 2 / float32(currAppHeight)
		cr.MaxCharsX = 80
		cr.MaxCharsY = 25
		cr.CharWid = float32(6) / float32(cr.MaxCharsX)
		cr.CharHei = float32(6) / float32(cr.MaxCharsY)
		cr.CharWidInPixels = int(float32(currAppWidth) / float32(cr.MaxCharsX))
		cr.CharHeiInPixels = int(float32(currAppHeight) / float32(cr.MaxCharsY))
	} else {
		cr.ClientExtentX = float32(currAppWidth) * cr.PixelWid / 2
		cr.ClientExtentY = float32(currAppHeight) * cr.PixelHei / 2
		fmt.Printf("CcRenderer.Init(): for resize changes - ClientExtentX: %.2f\n", cr.ClientExtentX)

		for _, pan := range cr.Panels {
			pan.Rect.Right = -cr.DistanceFromOrigin + 2*cr.ClientExtentX
			pan.BarVert.PosX = pan.Rect.Right - pan.BarVert.Thickness

			//pan.Init()
			/*
				if i == 0 {
					pan.Rect.Top = -cr.DistanceFromOrigin + 2*cr.ClientExtentY - rend.CharHei
					pan.Rect.Bottom = pan.Rect.Top - cr.ClientExtentY*2*pan.BandPercent
					pan.BarHori.PosY = pan.Rect.Bottom + pan.BarHori.Thickness
				} else {

				}
			*/
		}
	}

	if len(cr.Panels) == 0 {
		cr.Panels = append(cr.Panels, &TextPanel{BandPercent: 1 - cr.RunPanelHeiPerc, IsEditable: true})
		cr.Panels = append(cr.Panels, &TextPanel{BandPercent: cr.RunPanelHeiPerc, IsEditable: true}) // console (runtime feedback log)	// FIXME so its not editable once we're done debugging some things
		cr.Focused = cr.Panels[0]

		cr.Panels[0].Init()
		cr.Panels[0].SetupDemoProgram()
		cr.Panels[1].Init()
	}
}

func (cr *CcRenderer) Color(newColor []float32) {
	cr.PrevColor = cr.CurrColor
	cr.CurrColor = newColor
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &newColor[0])
}

func (cr *CcRenderer) DrawAll() {
	curs.Update()
	menu.Draw()

	for _, pan := range cr.Panels {
		pan.Draw()
	}
}

func (cr *CcRenderer) ScrollPanelThatIsHoveredOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	for _, pan := range cr.Panels {
		pan.ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY)
	}
}

func (cr *CcRenderer) DrawCharAtRect(char rune, r *Rectangle) {
	u := float32(int(char) % 16)
	v := float32(int(char) / 16)
	sp := rend.UvSpan

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(u*sp, v*sp+sp)
	gl.Vertex3f(r.Left, r.Bottom, 0)

	gl.TexCoord2f(u*sp+sp, v*sp+sp)
	gl.Vertex3f(r.Right, r.Bottom, 0)

	gl.TexCoord2f(u*sp+sp, v*sp)
	gl.Vertex3f(r.Right, r.Top, 0)

	gl.TexCoord2f(u*sp, v*sp)
	gl.Vertex3f(r.Left, r.Top, 0)
}

func (cr *CcRenderer) DrawQuad(atlasX, atlasY float32, r *Rectangle) {
	sp /* span */ := rend.UvSpan
	u := float32(atlasX) * sp
	v := float32(atlasY) * sp

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(r.Left, r.Bottom, 0)

	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(r.Right, r.Bottom, 0)

	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(r.Right, r.Top, 0)

	gl.TexCoord2f(u, v)
	gl.Vertex3f(r.Left, r.Top, 0)
}