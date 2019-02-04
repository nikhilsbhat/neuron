package commoncloud

type Clouds struct {
	Cloud   string `json:"cloud"`
	Region  string `json:"region"`
	Profile string `json:"profile"`
	GetRaw  bool   `json:"getraw"`
}
