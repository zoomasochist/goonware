package configurator

import (
	types "goonware/types"

	"fmt"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

func AnnoyancesTab(c *types.Config) g.Layout {
	return Tab("On", &c.Annoyances,
		LabelSliderTooltip("Seconds per tick", &c.TimerDelay, 1, 60, 250,
			"Number of seconds between attempts to spawn an annoyance", FormatSecondSlider),

		Feature("Popups", &c.AnnoyancePopups, 300,
			ChanceSlider(&c.PopupChance, "Chance of a popup opening"),

			LabelSliderTooltip("Popup opacity", &c.PopupOpacity, 1, 100, 150,
				"The opacity of the popup. 100 is opaque, 1 is almost invisible.",
				FormatPercentSlider),

			Setting("Denial mode", &c.PopupDenialMode,
				PercentChanceSlider(&c.PopupDenialChance),
			),

			Setting("Mitosis", &c.PopupMitosis,
				g.SliderInt(&c.PopupMitosisStrength, 2, 10).Size(75),
			),

			Setting("Timeout", &c.PopupTimeout,
				LabelSliderTooltip("Delay", &c.PopupTimeoutDelay, 1, 360, 150, "",
					FormatSecondSlider),
			),
		),

		Feature("Videos", &c.AnnoyanceVideos, 300,
			ChanceSlider(&c.VideoChance, "Chance of a video playing"),

			LabelledSlider("Video volume", &c.VideoVolume, 0, 100),
		),

		Feature("Prompts", &c.AnnoyancePrompts, 300,
			Setting("Allow mistakes", &c.AllowMistakes,
				g.Label("Max mistakes"),
				g.SliderInt(&c.MaxMistakes, 1, 100),
			),
		),

		Feature("Audio", &c.AnnoyanceAudio, 300,
			ChanceSlider(&c.AudioChance, "Chance of audio playing"),

			LabelledSlider("Audio volume", &c.AudioVolume, 0, 100),
		),
	)
}

func LoadPackage(c *types.Config) {
	filename, err := dialog.File().Filter("Edgeware package (.zip)", "zip").Load()

	if err != nil && err != dialog.Cancelled {
		dialog.Message("%s", fmt.Sprintf("Failed to load package; %s", err.Error())).Error()
	} else if err == nil {
		c.LoadedPackage = filename
		//pkg := LoadEdgewarePackage(filename)
	}
}

func SaveAndExit(c *types.Config) {
	fmt.Println("TODO: Shell out")
	SaveConfig(c)
	os.Exit(0)
}
