package configurator

import (
	"fmt"
	"strings"

	types "goonware/types"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

var newTag string
var booru int32
var boorus []string = []string{"https://e621.net/", "https://rule34.xxx/"}

func DriveFillerTab(c *types.Config) g.Layout {
	return Tab("On", &c.DriveFiller,
		g.PopupModal("Add New Tag").Layout(g.Layout{
			g.InputText(&newTag).Size(300),
			g.Button("Ok").OnClick(func() {
				if !(strings.TrimSpace(newTag) == "") {
					c.DriveFillerTags = append(c.DriveFillerTags, newTag)
					newTag = ""
				}
				g.CloseCurrentPopup()
			}),
		}),

		LabelSliderTooltip("Fill delay", &c.DriveFillerDelay, 10, 3000, 200,
			"Delay between each image save", FormatMillisecondSlider),

		g.Row(
			g.Button("Select root").OnClick(func() { SelectBase(c) }),
			g.Label("("+c.DriveFillerBase+")"),
		),
		g.Tooltip("Goonware won't save images in any directory above this one"),

		g.Row(
			g.Label("Image source"),
			TooltipRadio("Package", "Fill the drive with images taken from the loaded package",
				&c.DriveFillerImageSource, types.DriveFillerImageSourcePackage),
			TooltipRadio("Download", "Fill the drive with images downloaded from a booru",
				&c.DriveFillerImageSource, types.DriveFillerImageSourceBooru),
		),

		ShowIf(c.DriveFillerImageSource == types.DriveFillerImageSourceBooru,
			Setting("Minimum score", &c.DriveFillerDownloadMinimumScoreToggle,
				g.SliderInt(&c.DriveFillerDownloadMinimumScoreThreshold, -50, 100).Size(150),
			),

			g.Combo("Booru", c.DriveFillerBooru, boorus, &booru).
				Size(250).
				OnChange(func() { c.DriveFillerBooru = boorus[booru] }),

			Setting("Search specific tags", &c.DriveFillerImageUseTags,
				g.ListBox("Drive Filler Tags", c.DriveFillerTags).
					Border(false).
					Size(300, 300).
					ContextMenu([]string{"Remove"}).
					OnMenu(func(i int, m string) {
						c.DriveFillerTags = RemoveElement(c.DriveFillerTags, int32(i))
					},
					),
				g.Button("Add").OnClick(func() { g.OpenPopup("Add New Tag") }),
			),
		),
	)
}

func SelectBase(c *types.Config) {
	folder, err := dialog.Directory().Browse()

	if err != nil && err != dialog.Cancelled {
		dialog.Message("%s", fmt.Sprintf("Couldn't select directory: %s", err.Error())).Error()
	} else if err == nil {
		c.DriveFillerBase = folder
	}
}

func RemoveElement(slice []string, idx int32) []string {
	if len(slice) == 1 {
		return []string{}
	}

	return append(slice[:idx], slice[idx+1:]...)
}
