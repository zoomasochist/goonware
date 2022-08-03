package configurator

import (
	"fmt"

	g "github.com/AllenDang/giu"
)

func FormatMillisecondSlider(value int32) string {
	if value < 1000 {
		return fmt.Sprintf("%d milliseconds", value)
	}

	if value%1000 == 0 {
		if value == 1000 {
			return "1 second"
		}

		return fmt.Sprintf("%d seconds", value/1000)
	}

	return fmt.Sprintf("%d seconds %d milliseconds", value/1000, value%1000)
}

func FormatSecondSlider(value int32) string {
	if value < 60 {
		if value == 1 {
			return fmt.Sprintf("1 second")
		}

		return fmt.Sprintf("%d seconds", value)
	}

	if value%60 == 0 {
		if value == 60 {
			return fmt.Sprintf("1 minute")
		}

		return fmt.Sprintf("%d minutes", value/60)
	}

	if value < 120 {
		return fmt.Sprintf("1 minute %d seconds", value%60)
	}

	return fmt.Sprintf("%d minutes %d seconds", value/60, value%60)
}

func FormatMinuteSlider(value int32) string {
	if value < 60 {
		if value == 1 {
			return fmt.Sprintf("1 minute")
		}

		return fmt.Sprintf("%d minutes", value)
	}

	if value%60 == 0 {
		if value == 60 {
			return fmt.Sprintf("1 hour")
		}

		return fmt.Sprintf("%d hours", value/60)
	}

	if value < 120 {
		return fmt.Sprintf("1 hour %d minutes", value%60)
	}

	return fmt.Sprintf("%d hours, %d minutes", value/60, value%60)
}

func FormatPercentSlider(value int32) string {
	// I think having to do this is a consequence of me not exactly using the g.Slider* .Format
	// method correctly. It expects a string that it can pass to C's sprintf along with the value of
	// the slider. With FormatSecondSlider it so happens that the format it gets is along the lines
	// of sprintf("5 seconds");, so all is well. But here, %% would usually be used to escape the
	// first % into a literal %, but because giu passes it into sprintf, it ends up with
	// sprintf("5%");, which is invalid. so %%%% results in sprintf("5%%");, which is okay again.
	return fmt.Sprintf("%d%%%%", value)
}

func ShowIf(condition bool, layout ...g.Widget) g.Layout {
	if condition {
		return layout
	}

	return g.Layout{}
}

func LabelSliderTooltip(label string, value *int32, min, max int32, size float32, tooltip string,
	format func(int32) string) g.Widget {
	return g.Row(
		g.Label(label),
		g.SliderInt(value, min, max).Format(format(*value)).Size(size),
		g.Tooltip(tooltip),
	)
}

func LabelledSlider(label string, value *int32, min, max int32) g.Widget {
	return g.Row(
		g.Label(label),
		g.SliderInt(value, min, max).Size(150),
	)
}

func PercentChanceSlider(value *int32) g.Widget {
	return LabelSliderTooltip("Chance", value, 1, 100, 150, "Chance of occurance per tick",
		FormatSecondSlider)
}

func Setting(name string, selected *bool, widgets ...g.Widget) g.Layout {
	return g.Layout{
		g.Row(g.Checkbox(name, selected)),
		ShowIf(*selected, append(widgets, g.Dummy(g.Auto, 10))...),
	}
}

func Feature(name string, selected *bool, height float32, widgets ...g.Widget) g.Layout {
	return g.Layout{
		g.Row(
			g.Checkbox(name, selected),
		),

		ShowIf(*selected, g.Child().Layout(widgets...).Size(g.Auto, height)),
	}
}

func Tab(name string, selected *bool, widgets ...g.Widget) g.Layout {
	return g.Layout{
		g.Row(
			g.Checkbox(name, selected),
		),
		g.Separator(),
		g.Dummy(1, 10),

		ShowIf(*selected, widgets...),
	}
}

func ChanceSlider(value *int32, tooltip string) g.Widget {
	return g.Row(
		g.Label("Chance"),
		g.SliderInt(value, 1, 100).Format(FormatPercentSlider(*value)).Size(150),
		g.Tooltip(tooltip),
	)
}

func TooltipRadio(label, tooltip string, store *int32, selectedValue int32) g.Layout {
	return g.Layout{
		g.RadioButton(label, *store == selectedValue).OnChange(func() { *store = selectedValue }),
		g.Tooltip(tooltip),
	}
}
