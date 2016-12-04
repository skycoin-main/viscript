package gfx

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	//"github.com/corpusc/viscript/parser"
	"github.com/corpusc/viscript/tree"
	"github.com/corpusc/viscript/ui"
	"github.com/go-gl/gl/v2.1/gl"
)

type TextPanel struct {
	FractionOfStrip float32 // fraction of the parent PanelStrip (in 1 dimension)
	CursX           int     // current cursor/insert position (in character grid cells/units)
	CursY           int
	MouseX          int // current mouse position in character grid space (units/cells)
	MouseY          int
	IsEditable      bool // editing is hardwired to TextBodies[0], but we probably never want
	// to edit text unless the whole panel is dedicated to just one TextBody (& no graphical trees)
	Rect       *app.Rectangle
	Selection  *ui.SelectionRange
	BarHori    *ui.ScrollBar // horizontal
	BarVert    *ui.ScrollBar // vertical
	TextBodies [][]string
	Trees      []*tree.Tree
}

func (tp *TextPanel) Init() {
	fmt.Printf("TextPanel.Init()\n")

	tp.TextBodies = append(tp.TextBodies, []string{})

	tp.Selection = &ui.SelectionRange{}
	tp.Selection.Init()

	// scrollbars
	tp.BarHori = &ui.ScrollBar{IsHorizontal: true}
	tp.BarVert = &ui.ScrollBar{}
	tp.BarHori.Rect = &app.Rectangle{}
	tp.BarVert.Rect = &app.Rectangle{}

	tp.SetSize()
}

func (tp *TextPanel) SetSize() {
	fmt.Printf("TextPanel.SetSize()\n")

	tp.Rect = &app.Rectangle{
		Rend.ClientExtentY - Rend.CharHei,
		Rend.ClientExtentX,
		-Rend.ClientExtentY,
		-Rend.ClientExtentX}

	if tp.FractionOfStrip == Rend.RunPanelHeiPerc { // FIXME: this is hardwired for one use case for now
		tp.Rect.Top = tp.Rect.Bottom + tp.Rect.Height()*tp.FractionOfStrip
	} else {
		tp.Rect.Bottom = tp.Rect.Bottom + tp.Rect.Height()*Rend.RunPanelHeiPerc
	}

	// set scrollbars' upper left corners
	tp.BarHori.Rect.Left = tp.Rect.Left
	tp.BarHori.Rect.Top = tp.Rect.Bottom + ui.ScrollBarThickness
	tp.BarVert.Rect.Left = tp.Rect.Right - ui.ScrollBarThickness
	tp.BarVert.Rect.Top = tp.Rect.Top
}

func (tp *TextPanel) RespondToMouseClick() {
	Rend.Focused = tp

	// diffs/deltas from home position of panel (top left corner)
	glDeltaXFromHome := Curs.MouseGlX - tp.Rect.Left
	glDeltaYFromHome := Curs.MouseGlY - tp.Rect.Top
	tp.MouseX = int((glDeltaXFromHome + tp.BarHori.ScrollDelta) / Rend.CharWid)
	tp.MouseY = int(-(glDeltaYFromHome + tp.BarVert.ScrollDelta) / Rend.CharHei)

	if tp.MouseY < 0 {
		tp.MouseY = 0
	}

	if tp.MouseY >= len(tp.TextBodies[0]) {
		tp.MouseY = len(tp.TextBodies[0]) - 1
	}
}

func (tp *TextPanel) GoToTopEdge() {
	Rend.CurrY = tp.Rect.Top - tp.BarVert.ScrollDelta
}
func (tp *TextPanel) GoToLeftEdge() float32 {
	Rend.CurrX = tp.Rect.Left - tp.BarHori.ScrollDelta
	return Rend.CurrX
}
func (tp *TextPanel) GoToTopLeftCorner() {
	tp.GoToTopEdge()
	tp.GoToLeftEdge()
}

func (tp *TextPanel) Draw() {
	tp.GoToTopLeftCorner()
	tp.DrawBackground(11, 13)

	cX := Rend.CurrX // current (internal/logic cursor) drawing position
	cY := Rend.CurrY
	cW := Rend.CharWid
	cH := Rend.CharHei
	b := tp.BarHori.Rect.Top // bottom of text area

	// body of text
	for y, line := range tp.TextBodies[0] {
		// if line visible
		if cY <= tp.Rect.Top+cH && cY >= b {
			r := &app.Rectangle{cY, cX + cW, cY - cH, cX} // t, r, b, l

			// if line needs vertical adjustment
			if cY > tp.Rect.Top {
				r.Top = tp.Rect.Top
			}
			if cY-cH < b {
				r.Bottom = b
			}

			//parser.ParseLine(y, line, true)
			SetColor(Gray)

			// process line of text
			for x, c := range line {
				// if char visible
				if cX >= tp.Rect.Left-cW && cX < tp.BarVert.Rect.Left {
					app.ClampLeftAndRightOf(r, tp.Rect.Left, tp.BarVert.Rect.Left)
					Rend.DrawCharAtRect(c, r)

					if tp.IsEditable { //&& Curs.Visible == true {
						if x == tp.CursX && y == tp.CursY {
							SetColor(White)
							//Rend.DrawCharAtRect('_', r)
							Rend.DrawStretchableRect(11, 13, Curs.GetAnimationModifiedRect(*r))
							SetColor(PrevColor)
						}
					}
				}

				cX += cW
				r.Left = cX
				r.Right = cX + cW
			}

			// draw cursor at the end of line if needed
			if cX < tp.BarVert.Rect.Left && y == tp.CursY && tp.CursX == len(line) {
				if tp.IsEditable { //&& Curs.Visible == true {
					SetColor(White)
					app.ClampLeftAndRightOf(r, tp.Rect.Left, tp.BarVert.Rect.Left)
					//Rend.DrawCharAtRect('_', r)
					Rend.DrawStretchableRect(11, 13, Curs.GetAnimationModifiedRect(*r))
				}
			}

			cX = tp.GoToLeftEdge()
		}

		cY -= cH // go down a line height
	}

	SetColor(GrayDark)
	tp.DrawScrollbarChrome(10, 11, tp.Rect.Right-ui.ScrollBarThickness, tp.Rect.Top)                          // vertical bar background
	tp.DrawScrollbarChrome(13, 12, tp.Rect.Left, tp.Rect.Bottom+ui.ScrollBarThickness)                        // horizontal bar background
	tp.DrawScrollbarChrome(12, 11, tp.Rect.Right-ui.ScrollBarThickness, tp.Rect.Bottom+ui.ScrollBarThickness) // corner elbow piece
	SetColor(Gray)
	tp.BarHori.SetSize(tp.Rect, tp.TextBodies[0], cW, cH) // FIXME to consider multiple bodies & multiple trees
	tp.BarVert.SetSize(tp.Rect, tp.TextBodies[0], cW, cH)
	Rend.DrawStretchableRect(11, 13, tp.BarHori.Rect) // 2,11 (pixel checkerboard)    // 14, 15 (square in the middle)
	Rend.DrawStretchableRect(11, 13, tp.BarVert.Rect) // 13, 12 (double horizontal lines)    // 10, 11 (double vertical lines)
	SetColor(White)
}

// ATM the only different between the 2 funcs below is the top left corner (involving 3 vertices)
func (tp *TextPanel) DrawScrollbarChrome(atlasCellX, atlasCellY, l, t float32) { // left, top
	sp := Rend.UvSpan
	u := float32(atlasCellX) * sp
	v := float32(atlasCellY) * sp

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(l, tp.Rect.Bottom, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(tp.Rect.Right, tp.Rect.Bottom, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(tp.Rect.Right, t, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(l, t, 0)
}

func (tp *TextPanel) DrawBackground(atlasCellX, atlasCellY float32) {
	SetColor(GrayDark)
	Rend.DrawStretchableRect(atlasCellX, atlasCellY,
		&app.Rectangle{
			tp.Rect.Top,
			tp.Rect.Right - ui.ScrollBarThickness,
			tp.Rect.Bottom + ui.ScrollBarThickness,
			tp.Rect.Left})
}

func (tp *TextPanel) ScrollIfMouseOver(mousePixelDeltaX, mousePixelDeltaY float64) {
	if tp.ContainsMouseCursor() {
		// position increments in gl space
		xInc := float32(mousePixelDeltaX) * Rend.PixelWid
		yInc := float32(mousePixelDeltaY) * Rend.PixelHei
		tp.BarHori.Scroll(xInc)
		tp.BarVert.Scroll(yInc)
	}
}

func (tp *TextPanel) ContainsMouseCursor() bool {
	return MouseCursorIsInside(tp.Rect)
}

func (tp *TextPanel) ContainsMouseCursorInsideOfScrollBars() bool {
	return MouseCursorIsInside(&app.Rectangle{
		tp.Rect.Top, tp.Rect.Right - ui.ScrollBarThickness, tp.Rect.Bottom + ui.ScrollBarThickness, tp.Rect.Left})
}

func (tp *TextPanel) RemoveCharacter(fromUnderCursor bool) {
	txt := tp.TextBodies[0]

	if fromUnderCursor {
		if len(txt[tp.CursY]) > tp.CursX {
			txt[tp.CursY] = txt[tp.CursY][:tp.CursX] + txt[tp.CursY][tp.CursX+1:len(txt[tp.CursY])]
		}
	} else {
		if tp.CursX > 0 {
			txt[tp.CursY] = txt[tp.CursY][:tp.CursX-1] + txt[tp.CursY][tp.CursX:len(txt[tp.CursY])]
			tp.CursX--
		}
	}
}

func (tp *TextPanel) SetupDemoProgram() {
	txt := []string{}

	txt = append(txt, "// ------- variable declarations ------- -------")
	//txt = append(txt, "var myVar int32")
	txt = append(txt, "var a int32 = 42")
	txt = append(txt, "var b int32 = 58")
	txt = append(txt, "")
	txt = append(txt, "// ------- builtin function calls ------- ------- ------- ------- ------- ------- ------- end")
	txt = append(txt, "//    sub32(7, 9)")
	//txt = append(txt, "sub32(4,8)")
	//txt = append(txt, "mult32(7, 7)")
	//txt = append(txt, "mult32(3,5)")
	//txt = append(txt, "div32(8,2)")
	//txt = append(txt, "div32(15,  3)")
	//txt = append(txt, "add32(2,3)")
	//txt = append(txt, "add32(a, b)")
	txt = append(txt, "")
	txt = append(txt, "// ------- user function calls -------")
	txt = append(txt, "myFunc(a, b)")
	txt = append(txt, "")
	txt = append(txt, "// ------- function declarations -------")
	txt = append(txt, "func myFunc(a int32, b int32){")
	txt = append(txt, "")
	txt = append(txt, "        div32(6, 2)")
	txt = append(txt, "        innerFunc(a,b)")
	txt = append(txt, "}")
	txt = append(txt, "")
	txt = append(txt, "func innerFunc (a, b int32) {")
	txt = append(txt, "        var locA int32 = 71")
	txt = append(txt, "        var locB int32 = 29")
	txt = append(txt, "        sub32(locA, locB)")
	txt = append(txt, "}")

	/*
		for i := 0; i < 22; i++ {
			txt = append(txt, fmt.Sprintf("%d: put lots of text on screen", i))
		}
	*/

	tp.TextBodies[0] = txt
}