package daemon

import (
	"fmt"
	types "goonware/types"

	"math/rand"

	g "github.com/AllenDang/giu"
)

var showingPrompt bool
var showingPopup bool

func DoAnnoyances(c types.Config, pkg types.EdgewarePackage) {

	willDoPopup := (rand.Intn(101) + 1) < int(c.PopupChance)
	if c.AnnoyancePopups && willDoPopup {
		image := pkg.ImageFiles[rand.Intn(len(pkg.ImageFiles))]

		go MakePopup(image)
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
	showingPopup = true
	wnd := g.NewMasterWindow("Goonware"+fmt.Sprint(rand.Int()), 500, 600,
		g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|
			g.MasterWindowFlagsFloating)
	wnd.Run(func() {
		g.SingleWindow().Layout(
			g.ImageWithFile(imagePath).Size(g.Auto, g.Auto),
			g.Button("Submit <3").OnClick(func() { wnd.Close(); showingPopup = false }),
		)
	})
}

func MakePrompt(text string) {
	var input string

	showingPrompt = true
	wnd := g.NewMasterWindow("Goonware"+string(rand.Int()), 500, 300,
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
