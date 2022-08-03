package configurator

import (
	types "goonware/types"

	"encoding/json"
	"errors"
	"os"

	g "github.com/AllenDang/giu"
)

var workingDirectory = Expect(os.UserConfigDir()) + "/goonware"
var packageExtractDirectory = workingDirectory + "/package"
var configFileLocation = workingDirectory + "/goonware.json"

func ConfiguratorUI() error {
	c, err := NewOrLoadConfig()
	if err != nil {
		return err
	}

	wnd := g.NewMasterWindow("Goonware", 700, 700, g.MasterWindowFlagsNotResizable)
	wnd.Run(func() {
		g.SingleWindow().Layout(
			g.TabBar().TabItems(
				g.TabItem("Annoyances").Layout(TabWrapper(AnnoyancesTab(&c))),
				g.TabItem("Drive Filler").Layout(TabWrapper(DriveFillerTab(&c))),
				g.TabItem("Passive").Layout(TabWrapper(PassiveTab(&c)...)),
				// Comes last
				g.TabItem("About").Layout(TabWrapper(AboutTab()...)),
			),

			g.Row(
				g.Button("Load package").OnClick(func() { LoadPackage(&c) }),
				g.Condition(len(c.LoadedPackage) != 0,
					g.Layout{g.Label("(Using package " + c.LoadedPackage + ")")},
					g.Layout{g.Label("(No package loaded)")},
				),
			),

			g.Row(
				g.Button("Save").OnClick(func() { SaveConfig(&c) }),
				g.Button("Save and Exit").OnClick(func() { SaveAndExit(&c) }),
				g.Row(
					g.Checkbox("Launch on startup", &c.StartOnBoot),
					g.Tooltip("Silently start Goonware every time your computer starts, executing"+
						" whatever package was running last time."),
					g.Checkbox("Run Goonware on save and exit", &c.RunOnExit),
				),
			),

			g.Row(
				g.Label("Mode"),
				g.RadioButton("Normal", c.Mode == types.ModeNormal).OnChange(func() {
					c.Mode = types.ModeNormal
				}),
				g.Tooltip("As soon as Goonware starts, it will start running payloads."),

				g.RadioButton("Hibernate", c.Mode == types.ModeHibernate).OnChange(func() {
					c.Mode = types.ModeHibernate
				}),
				g.Tooltip("Goonware will wait a random amount of time (within given limits) before"+
					"\nspamming payloads, then stop and start waiting again."),

				ShowIf(c.Mode == types.ModeHibernate,
					LabelSliderTooltip("Min. wait", &c.HibernateMinWaitMinutes, 0, 120, 50,
						"The minimum amount of time Goonware will hibernate", FormatMinuteSlider),
					LabelSliderTooltip("Max. wait", &c.HibernateMaxWaitMinutes, 0, 120, 50,
						"The maximum amount of time Goonware will hibernate", FormatMinuteSlider),
					LabelSliderTooltip("Wake for", &c.HibernateActivityLength, 1,
						600, 50, "How long Goonware will spam payloads", FormatSecondSlider),
				),
			),
		)
	})

	return nil
}

func TabWrapper(w ...g.Widget) g.Widget {
	return g.Child().Layout(w...).Size(g.Auto, 580).Flags(g.WindowFlagsNoScrollbar)
}

func NewOrLoadConfig() (types.Config, error) {
	if _, err := os.Stat(configFileLocation); errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(packageExtractDirectory, os.ModePerm)

		return types.Config{
			WorkingDirectory: workingDirectory,
			// General
			Mode:                    types.ModeNormal,
			HibernateMinWaitMinutes: 120,
			HibernateMaxWaitMinutes: 3600,
			HibernateActivityLength: 20, // Sec
			StartOnBoot:             false,
			RunOnExit:               false,
			LoadedPackage:           "",
			LoadedPackageData:       nil,

			// Annoyances
			Annoyances:      true,
			TimerDelay:      10,
			AnnoyancePopups: true,

			PopupChance:          50,
			PopupOpacity:         85,
			PopupDenialMode:      false,
			PopupDenialChance:    50,
			PopupMitosis:         true,
			PopupMitosisStrength: 4,
			PopupTimeout:         false,
			PopupTimeoutDelay:    30,

			AnnoyanceVideos: true,
			VideoChance:     25,
			VideoVolume:     25,

			AnnoyancePrompts: true,
			PromptChance:     25,
			AllowMistakes:    true,
			MaxMistakes:      1,

			AnnoyanceAudio: true,
			AudioChance:    25,
			AudioVolume:    25,

			// Drive Filler
			DriveFiller:                              false,
			DriveFillerDelay:                         1000, // Ms
			DriveFillerBase:                          Expect(os.UserHomeDir()),
			DriveFillerBooru:                         "https://e621.net/",
			DriveFillerTags:                          []string{"feral+paws", "feral+rimming"},
			DriveFillerImageSource:                   types.DriveFillerImageSourceBooru,
			DriveFillerImageUseTags:                  true,
			DriveFillerDownloadMinimumScoreToggle:    false,
			DriveFillerDownloadMinimumScoreThreshold: 0,
		}, nil
	}

	return LoadConfig()
}

func SaveConfig(c *types.Config) error {
	structBytes, err := json.Marshal(*c)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFileLocation, structBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() (types.Config, error) {
	structBytes, err := os.ReadFile(configFileLocation)
	if err != nil {
		return types.Config{}, err
	}

	var c types.Config
	err = json.Unmarshal(structBytes, &c)
	if err != nil {
		return types.Config{}, err
	}

	return c, nil
}

func Expect(s string, err error) string {
	if err != nil {
		panic(err)
	}

	return s
}
