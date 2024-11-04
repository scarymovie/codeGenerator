package internal

type OpenAPI struct {
	Paths map[string]map[string]struct {
		OperationId string `yaml:"operationId"`
		Parameters  []struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			In          string `yaml:"in"`
			Schema      struct {
				Type string `yaml:"type"`
			} `yaml:"schema"`
			Required bool `yaml:"required"`
		} `yaml:"parameters"`
		Responses map[string]struct {
			Description string `yaml:"description"`
			Content     map[string]struct {
				Schema struct {
					Type  string `yaml:"type"`
					Items struct {
						Ref string `yaml:"$ref"`
					} `yaml:"items"`
					Ref string `yaml:"$ref"`
				} `yaml:"schema"`
			} `yaml:"content"`
		} `yaml:"responses"`
	} `yaml:"paths"`
	Components struct {
		Schemas map[string]struct {
			Type        string   `yaml:"type"`
			Description string   `yaml:"description"`
			Required    []string `yaml:"required"`
			Properties  map[string]struct {
				Type        string `yaml:"type"`
				Description string `yaml:"description"`
				Format      string `yaml:"format"`
			} `yaml:"properties"`
		} `yaml:"schemas"`
	} `yaml:"components"`
}
