package syncsign

import config "github.com/tommzn/go-config"

// NewDisplayConfig extracts list of display ids from passed config and returns a
// DisplayCondig which can be used to ensure valid display ids.
func NewDisplayConfig(conf config.Config) *DisplayConfig {

	displays := make(map[string]struct{})
	displaysCfg := conf.GetAsSliceOfMaps("hdb.displays")
	for _, displayCfg := range displaysCfg {
		if displayId, ok := displayCfg["id"]; ok {
			displays[displayId] = struct{}{}
		}
	}
	return &DisplayConfig{displays: displays}
}

// Exists returns true if passed display id is available in internal display list.
func (cfg *DisplayConfig) Exists(displayId string) bool {
	_, ok := cfg.displays[displayId]
	return ok
}

// All returns the list of all available display ids.
func (cfg *DisplayConfig) All() []string {
	displayIds := []string{}
	for displayId, _ := range cfg.displays {
		displayIds = append(displayIds, displayId)
	}
	return displayIds
}
