package configurator

import (
	types "goonware/types"

	g "github.com/AllenDang/giu"
)

func PassiveTab(c *types.Config) g.Layout {
	return Tab("On", &c.Passive,
		Setting("Autotype", &c.PassiveAutoType,
			g.Checkbox("Automatically press enter", &c.PassiveAutoTypeAutoEnter),
			g.Tooltip("Automatically press enter after typing a phrase"),
		),
	)
}
