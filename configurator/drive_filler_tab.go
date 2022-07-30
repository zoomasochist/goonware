package configurator

import (
	"fmt"
	"strings"

	types "goonware/types"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

var selectedTag int32
var newTag string
var boorus []string = []string{"https://e621.net/", "https://rule34.xxx/"}

func DriveFillterTab(c *types.Config) []g.Widget {
	largerFont := g.GetDefaultFonts()[0].SetSize(20)

	return []g.Widget{
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

		g.Checkbox("On", &c.DriveFiller),

		ConditionOrNothing(c.DriveFiller, g.Layout{
			g.Row(
				LabelSliderTooltip("Fill delay", &c.DriveFillerDelay, 10, 3000, 200,
					"How many milliseconds to wait before writing another image",
					FormatMillisecondSlider),
			),
			g.Row(
				g.Button("Select base").OnClick(func() { SelectBase(c) }),
				g.Label("(" + c.DriveFillerBase + ")"),
			),
			StandardSeparation(),

			g.Row(
				g.Label("Image source"),
				g.RadioButton("Use Package", c.DriveFillerImageSource == 0).
					OnChange(func() { c.DriveFillerImageSource = 0 }),
				g.Tooltip("Fill the drive with files in the currently loaded package's img directory"),

				g.RadioButton("Download", c.DriveFillerImageSource == 1).
					OnChange(func() { c.DriveFillerImageSource = 1 }),
				g.Tooltip("Fill the drive with images downloaded from a booru of your choosing"),
			),
			StandardSeparation(),

			ConditionOrNothing(c.DriveFillerImageSource == 1, g.Layout{
				g.Row(g.Label("Image Downloader").Font(largerFont)),
				StandardSeparation(),

				g.Row(
					g.Combo("Booru", boorus[c.DriveFillerBooru], boorus, &c.DriveFillerBooru).
						Size(250),
				),

				g.Row(
					g.RadioButton("Anything", !c.DriveFillerImageUseTags).
						OnChange(func() { c.DriveFillerImageUseTags = false }),
					g.RadioButton("Specific Tags", c.DriveFillerImageUseTags).
						OnChange(func() { c.DriveFillerImageUseTags = true}),
				),
				ConditionOrNothing(c.DriveFillerImageUseTags, g.Layout{
					g.Child().Layout(g.Layout{
						g.ListBox("Drive Filler Tags", c.DriveFillerTags).
							Border(false).
							Size(g.Auto, g.Auto).
							ContextMenu([]string{"Remove"}).
								OnMenu(func(i int, m string) {
									c.DriveFillerTags = RemoveElement(c.DriveFillerTags, int32(i))
								}),
					}).Size(300, 300),	
					g.Row(
						g.Button("Add").OnClick(func() { g.OpenPopup("Add New Tag") }),
					),
				}),

				g.Row(
					g.Checkbox("Minimum score", &c.DriveFillerDownloadMinimumScoreToggle),
					ConditionOrNothing(c.DriveFillerDownloadMinimumScoreToggle,
						g.Layout{g.SliderInt(&c.DriveFillerDownloadMinimumScoreThreshold, -50, 100).
							Size(150)},
					),
				),
				StandardSeparation(),
			}),
		}),
	}
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