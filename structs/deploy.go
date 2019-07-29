package structs

//FullConfig .
type FullConfig struct {
	ID           string                       `json:"_id"`
	AWSConfig    map[string]DesignationConfig `json:"aws-stages,omitempty"`
	CampusConfig map[string]DesignationConfig `json:"campus-stages,omitempty"`
}

//DesignationConfig .
type DesignationConfig struct {
	Task                 string                 `json:"task,omitempty"`
	Port                 string                 `json:"port"`
	EnvironmentVariables map[string]string      `json:"environment-values,omitempty"`
	DockerInfo           map[string]interface{} `json:"docker-info,omitempty"`
}

//DeviceDeploymentConfig .
type DeviceDeploymentConfig struct {
	ID           string                                 `json:"_id"`
	Designations map[string]DesignationDeploymentConfig `json:"designations"`
}

//DesignationDeploymentConfig .
type DesignationDeploymentConfig struct {
	EnvironmentVariables map[string]string      `json:"environment-values,omitempty"`
	DockerInfo           map[string]interface{} `json:"docker-info,omitempty"`
	DockerServices       []string               `json:"docker-services,omitempty"`
	Services             []string               `json:"services,omitempty"`
}

// ServiceConfigWrapper .
type ServiceConfigWrapper struct {
	ID           string                   `json:"_id"`
	Designations map[string]ServiceConfig `json:"designations,omitempty"`
}

// ServiceConfig .
type ServiceConfig struct {
	Data map[string]map[string]string `json:"data"`
}
