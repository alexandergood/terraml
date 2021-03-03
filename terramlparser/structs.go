package terramlparser

type Module struct {
	Src      string				    `yaml:"src"`
	Name     string   			    `yaml:"name"`
	Var      map[string]interface{} `yaml:"var"`
}

type Resource struct {
	ResourceType string 			    `yaml:"type"`
	Name         string 			    `yaml:"name"`
	Config       map[string]interface{} `yaml:"config"`
}

type Task struct {
	Module   Module   `yaml:"module"`
	Resource Resource `yaml:"resource"`
}

type Provider struct {
	Name   string 	                `yaml:"name"`
	Config []map[string]interface{} `yaml:"config"`
}

type RemoteState struct {
	BackendType string 		   `yaml:"backend"`
	Config      map[string]interface{} `yaml:"config"`
}

type TerraformConf struct {
	RemoteState  RemoteState `yaml:"remote_state"`
}

type TerramlFileContent struct {
	TerraformConf  TerraformConf `yaml:"terraform_conf"`
	Providers      []Provider    `yaml:"providers"`
	Provision      []Task	     `yaml:"provision"`
}

func (m *Module) IsEmpty() bool {
	if m.Src == "" && m.Name == ""  {
		return true
	} else {
		return false
	}
}

func (r *Resource) IsEmpty() bool {
	if r.ResourceType == "" && r.Name == "" {
		return true
	} else {
		return false
	}
}
