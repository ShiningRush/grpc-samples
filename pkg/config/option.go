package config

func ConfigDirectory(dir string) OptAction {
	return func(o *configOption) {
		if dir == "" {
			panic("config.ConfigDirectory: config directory is empty")
		}

		o.dir = dir
	}
}

func ConfigName(name string) OptAction {
	return func(o *configOption) {
		if name == "" {
			panic("config.ConfigDirectory: config name is empty")
		}

		o.name = name
	}
}
