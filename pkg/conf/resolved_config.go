package conf

type ResolvedConfig struct {
	App          *AppConfig
	Repos        map[string]*RepoSpec
	Charts       map[string]*ChartSpec
	Environments map[string]*Environment
}

func NewEmptyResolvedConfig() *ResolvedConfig {
	return &ResolvedConfig{
		App: &AppConfig{
			DebugHelm:       false,
			ChartsDir:       "",
			EnvironmentsDir: "",
		},
		Repos:        make(map[string]*RepoSpec),
		Charts:       make(map[string]*ChartSpec),
		Environments: make(map[string]*Environment),
	}
}
