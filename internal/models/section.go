package models

// Section represents a content section in the portfolio
type Section struct {
	Title   string
	Content []string
}

// Sections contains the content for all portfolio sections
var Sections = []Section{
	{"about", []string{
		"i am a software engineer and full-time observer and tinkerer.",
		"i love all kinds of engineering and development. i love free software, freedom in general.",
	}},
	{"experience", []string{
		"💼 software engineer @ akgun technology (2025 - Present)",
		"🧪 software developer intern @ comp. (2020 - 2022)",
		"🧑‍🎓 software developer intern @ comp. (2020 - 2021)",
		"📚 software engineering student @ canakkale onsekiz mart university -- turkey (2017 - 2023)",
	}},
	{"projects", []string{
		"🔧 ssh tui portfolio",
		"",
		"",
		"",
	}},
	{"links", []string{
		"github: https://github.com/cankurttekin",
		"linkedin: https://linkedin.com/in/cankurttekin",
	}},
}
