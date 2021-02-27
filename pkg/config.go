package pkg

type Config struct {
	Name                         string `yaml:"name"`
	Password                     string `yaml:"password"`
	PollingSeconds               int    `yaml:"pollingSeconds"`
	ScoreServerURL               string `yaml:"scoreServerURL"`
	RedeploymentApiServerAddress string `yaml:"redeploymentApiServerAddress"`
	RedeploymentApiServerPort    int    `yaml:"redeploymentApiServerPort"`
	Group                        string `yaml:"group"`
	PenaltyMinutes               int    `yaml:"penaltyMinutes"`
}
