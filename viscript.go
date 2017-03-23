/*

------- NEXT THINGS TODO: -------

* limit resizing to require at least 16 char columns

* make command line dynamic based on terminal's .GridSize
	should MaxCommandSize be dynamic?  i don't like the idea
	of reserving more than 2 lines at the bottom of a terminal
	(what happens after we've autoscrolled down a whole page)

* Resize terminals
	* i believe we'll change the actual grid size, then get enough data
			from the terminal task to fill the backscroll
	* change flow of text with wrapping, so for example,
			squeezing horizontally would cause more lines

* back buffer scrolling
	* pgup/pgdn hotkeys
	* 1-3 lines with scrollwheel

* Fix getting a resizing pointer outside of focused terminal.
		When you click outside terminal it can land on a background
		terminal which then pops in front.  Blocking the resize




------- OLDER TODO: ------- (everything below was for the text editor)

* KEY-BASED NAVIGATION
	* CTRL-HOME/END - PGUP/DN
* BACKSPACE/DELETE at the ends of lines
	pulls us up to prev line, or pulls up next line
* when there is no scrollbar, should be able to see/interact with text in that area
* when auto appending to the end of a terminal, scroll all the way down
		(manual activity in the middle could increase size, so do this only when appending to body)


------- LOWER PRIORITY POLISH: -------

* if typing goes past right of screen, auto-horizontal-scroll as you type
* same for when newlines/enters/returns push cursor past the bottom of visible space
* scrollbars should have a bottom right corner, and a thickness sized background
		for void space, reserved for only that, so the bar never covers up the rightmost
		character/cursor
* when pressing delete at/after the end of a line, should pull up the line below
* vertical scrollbars could have a smaller rendering of the first ~40 chars?
		however not if we map the whole vertical space (when scrollspace is taller than screen),
		because this requires scaling the text.  and keeping the aspect ratio means ~40 (max)
		would alter the width of the scrollbar

*/

package main

import (
	"github.com/corpusc/viscript/hypervisor"
	"github.com/corpusc/viscript/rpc/terminalmanager"
	"github.com/corpusc/viscript/viewport"
)

func main() {
	println("Starting...")

	hypervisor.Init()

	viewport.DebugPrintInputEvents = true
	viewport.Init() //runtime.LockOSThread(), InitCanvas()

	// rpc
	go func() {
		rpcInstance := terminalmanager.NewRPC()
		rpcInstance.Serve()
	}()

	println("Start Loop;")
	for viewport.CloseWindow == false {
		viewport.DispatchEvents() //event channel

		hypervisor.TickTasks()

		viewport.PollUiInputEvents()
		viewport.Tick()
		viewport.UpdateDrawBuffer()
		viewport.SwapDrawBuffer() //with new frame
	}

	println("Closing down viewport")
	viewport.TeardownScreen()
	hypervisor.HypervisorTeardown()
}
