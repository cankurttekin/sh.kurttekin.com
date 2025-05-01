package models

type Portfolio struct {
	Title    string    // User's name or title
	Sections []Section // Content sections
	Theme    Theme     // Color theme settings
}

type Theme struct {
	Primary   string // Primary color (for highlights, borders)
	Accent    string // Accent color (for selected items)
	Text      string // Main text color
	Subtle    string // Subtle text color (for secondary information)
	Links     string // Color for links
	Selection string // Color for selected links
}

func DefaultPortfolio() Portfolio {
	return Portfolio{
		Title: "cankurttekin",
		Sections: []Section{
			{
				Title: "about",
				Content: []string{
					"i am a software engineer and full-time observer and tinkerer.",
					"i love all kinds of engineering and development. i love free software, freedom in general.\n",
					"github: https://github.com/cankurttekin",
					"linkedin: https://linkedin.com/in/cankurttekin",
					"website: https://can.kurttekin.com",
					"blog: https://blog.kurttekin.com",
					"gpg: https://pgp.mit.edu/pks/lookup?op=get&search=0xAC9A980E2",
				},
			},
			/*
			{
				Title: "experience",
				Content: []string{
					"software engineer @ akgun technology (2025 - Present)",
					"software developer intern @ comp. (2020 - 2022)",
					"software developer intern @ comp. (2020 - 2021)",
					"computer engineering @ canakkale onsekiz mart university -- turkey (2017 - 2023)",
				},
			},
			*/
			{
				Title: "projects",
				Content: []string{
					"",
					"",
					"",
					"",
				},
			},
			{
				Title: "my setup",
				Content: []string{
					"fedora with swaywm no ricing",
					"text editor: neovim",
					"terminal: foot",
					"browser: fennec on android, firefox on desktop with vimium",
					"ad blocking: old android phone running debian(chroot) with pi-hole dns sinkhole and ublock origin on browsers",
					"dotfiles: https://github.com/cankurttekin/dotfiles",
				},
			},
			{
				Title: "bookmarks",
				Content: []string{
					"brodierobertson: https://www.youtube.com/@BrodieRobertson",
					"theprimeagen: https://www.youtube.com/channel/UC8ENHE5xdFSwx71u3fDH5Xw",
					"technology connections: https://www.youtube.com/@TechnologyConnections",
					"bigclivedotcom: https://www.youtube.com/@bigclivedotcom",
					"computerphile: https://www.youtube.com/@Computerphile",
					"low level: https://www.youtube.com/@LowLevel",

				},
			},
		},
		Theme: Theme{
			Primary:   "#5f87ff", // Vibrant blue
			Accent:    "#ff6ac1", // Pink
			Text:      "#abb2bf", // Light gray
			Subtle:    "#565c64", // Dark gray
			Links:     "#61afef", // Light blue for links
			Selection: "#c678dd", // Purple for selections
		},
	}
}

// GetPortfolio returns the portfolio information
// This is the function you would modify if you want to load from a file instead
func GetPortfolio() Portfolio {
	return DefaultPortfolio()
}
