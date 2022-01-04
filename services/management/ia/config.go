package ia

type Config struct {
	Enabled       bool    `yaml:"enabled"`
	TemplatesDir  *string `yaml:"templates_dir"`
	RulesDir      *string `yaml:"rules_dir"`
	DirOwner      *string `yaml:"dir_owner"`
	DirOwnerGroup *string `yaml:"dir_owner_group"`
}
