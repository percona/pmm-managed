package ia

type Config struct {
	Enabled       bool    `yaml:"enabled"`
	TemplatesDir  *string `yaml:"templates_dir"`
	RulesDir      *string `yaml:"rules_dir"`
	DirOwner      *string `yaml:"dir_owner"`
	DirOwnerGroup *string `yaml:"dir_owner_group"`
}

func (c *Config) Init() {
	if c.TemplatesDir == nil {
		*c.TemplatesDir = "/srv/ia/templates"
	}
	if c.RulesDir == nil {
		*c.RulesDir = "/etc/ia/rules"
	}
	if c.DirOwner == nil {
		*c.DirOwner = "pmm"
	}
	if c.DirOwnerGroup == nil {
		*c.DirOwnerGroup = "pmm"
	}
}
