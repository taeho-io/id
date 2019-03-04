package server

type Settings struct {}

func NewSettings() Settings {
	return Settings{}
}

func MockSettings() Settings {
	return NewSettings()
}

