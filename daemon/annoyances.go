package daemon

import (
	"fmt"
	types "goonware/types"
	"image"
	_ "image/jpeg"
	"os"
	"strings"
	"time"

	"math/rand"

	g "github.com/AllenDang/giu"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var showingPrompt bool
var showingPopups int32

// There's something deeply awry with something this code is doing.
// I think the root cause is Giu, but I'm not sure precisely what I'm doing to trigger it.
// I've gotten segfaults and panics from GLFW, imgui, Giu, and Win32, images not showing up at
// all, or having Giu draw totally random things (like a preview of a bunch of random characters?).
// Clearly something is very wrong somewhere, but I've found a few things that seems to mitigate
// whatever the problem is. There's a few comments below elaborating on this.
func DoAnnoyances(c *types.Config, pkg *types.EdgewarePackage) {

	willDoPopup := (rand.Intn(101) + 1) < int(c.PopupChance)
	if c.AnnoyancePopups && willDoPopup && showingPopups < 5 {
		imagePath := pkg.ImageFiles[rand.Intn(len(pkg.ImageFiles))]

		showingPopups++
		// This fixed most (all but one) of the issues I found, with GLFW/Giu/Windows/Imgui usually
		// being confused about where its window context is - I assume in high-popup configurations
		// windows are being closed and opened so fast that one of these libraries tries to reuse
		// a context just as its being deleted? Idk. But this works.
		time.Sleep(600 * time.Millisecond)
		go MakePopup(imagePath)
	}

	willDoVideo := (rand.Intn(101) + 1) < int(c.VideoChance)
	if c.AnnoyanceVideos && willDoVideo {
		// Todo
	}

	willDoPrompt := (rand.Intn(101) + 1) < int(c.PromptChance)
	if c.AnnoyancePrompts && !showingPrompt && willDoPrompt {
		prompts := pkg.Prompts[rand.Intn(len(pkg.Prompts))].Prompts
		prompt := prompts[rand.Intn(len(prompts))]

		go MakePrompt(prompt)
	}

	willDoAudio := (rand.Intn(101) + 1) < int(c.AudioChance)
	if c.AnnoyanceAudio && willDoAudio {
		// Todo
	}
}

func MakePopup(imagePath string) {
	// TODO: This is an awfully hacky solution to a bug in Giu. Apparently it doesn't properly
	// clean up after itself or something (probably because wnd.Close is a very recent addition);
	// when an ImageWithFile is created its id (which apparently is used to cache information about
	// the image) is based entirely off of the image path (treated as a string). When a window with
	// an image is created, destroyed, and then a new window is created pointing to the same image
	// is created, it (i suppose) still uses that old id it found despite the image data itself
	// having been destroyed, resulting in image popups appearing more than once appearing garbled.
	// This fixes that. It's also bad. Should be reported to github.com/AllenDang/giu.
	hotfixPath := imagePath + strings.Repeat(" ", rand.Intn(150))
	fmt.Println("'" + hotfixPath + "'")
	w, h := ImageDimensions(imagePath)
	for w > 800 || h > 600 {
		w = w / 2
		h = h / 2
	}

	x, y := GetNextWindowPos()

	wnd := g.NewMasterWindow("Goonware Popup"+fmt.Sprint(rand.Int()), w/2, h/2,
		g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|
			g.MasterWindowFlagsFloating)
	wnd.SetPos(x, y)
	wnd.Run(func() {
		g.SingleWindow().Layout(
			g.ImageWithFile(hotfixPath).Size(g.Auto, g.Auto),
			g.Button("Submit <3").OnClick(func() { wnd.Close(); showingPopups-- }),
		)
	})
}

func MakePrompt(text string) {
	var input string

	showingPrompt = true
	wnd := g.NewMasterWindow("Goonware Prompt"+fmt.Sprint(rand.Int()), 500, 300,
		g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|
			g.MasterWindowFlagsFloating)
	largerFont := g.GetDefaultFonts()[0].SetSize(20)

	wnd.Run(func() {
		g.SingleWindow().Layout(
			g.Label("Repeat.").Font(largerFont),
			g.Dummy(g.Auto/3, 1),
			g.Label(text).Wrapped(true),
			g.Dummy(1, 20),
			g.InputTextMultiline(&input).Size(g.Auto, g.Auto).OnChange(func() {
				if text == input {
					wnd.Close()
					showingPrompt = false
				}
			}),
		)
	})
}

func ImageDimensions(imagePath string) (int, int) {
	if reader, err := os.Open(imagePath); err == nil {
		defer reader.Close()
		image, _, err := image.DecodeConfig(reader)
		if err != nil {
			panic(err)
		}
		return image.Width, image.Height
	}
	panic("Couldn't get dimensions of image")
}

func GetNextWindowPos() (int, int) {
	//glfw.Init()
	//defer glfw.Terminate()

	monitors := glfw.GetMonitors()
	monitor := monitors[rand.Intn(len(monitors))]
	videoMode := monitor.GetVideoMode()

	x := rand.Intn(videoMode.Height)
	y := rand.Intn(videoMode.Width)

	return x, y
}
