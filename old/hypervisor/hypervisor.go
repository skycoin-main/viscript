package hypervisor

import (
	"fmt"

	"github.com/corpusc/viscript/app"
	//"github.com/corpusc/viscript/gfx"
	//"github.com/corpusc/viscript/msg"
	//"github.com/corpusc/viscript/script"
	//"github.com/go-gl/gl/v2.1/gl"
	//"github.com/go-gl/glfw/v3.2/glfw"
	//"log"
	"runtime"

	igl "github.com/corpusc/viscript/gl" //internal gl
)

//glfw
//glfw.PollEvents()
//only remaining

var CloseWindow chan int // write to channel to close

func init() {
	fmt.Println("hypervisor.init()")
	CloseWindow = make(chan int) //write to close windows
	initTerminals()
}

func HypervisorInit() {
	fmt.Println("hypervisor.HypervisorInit()")
	app.MakeHighlyVisibleLogHeader(app.Name, 15)
	// GLFW event handling must run on the main OS thread
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func HypervisorScreenTeardown() {
	igl.ScreenTeardown()
}

func HypervisorScreenInit() {
	fmt.Printf("Hypervisor: screen init\n")
	igl.WindowInit()
	igl.LoadTextures()
	igl.InitRenderer()
}

func HypervisorInitInputEvents() {
	fmt.Printf("Hypervisor: init InitInputEvents \n")
	igl.InitInputEvents(igl.GlfwWindow)
	igl.InitMiscEvents(igl.GlfwWindow)
}

func PollUiInputEvents() {
	igl.PollEvents() //move to gl
}

//could be in messages
func DispatchInputEvents(ch chan []byte) []byte {
	message := []byte{}

	for len(ch) > 0 { //if channel has elements
		v := <-ch //read from channel
		message = ProcessInputEvents(v)

	}

	return message
}

func UpdateDrawBuffer() {
	igl.DrawScene()

	// for _, term := range Terms {
	// 	//fmt.Println("drawing a terminal", term.FractionOfStrip)
	// 	term.Draw()
	// }
}

func SwapDrawBuffer() {
	igl.SwapDrawBuffer()
}