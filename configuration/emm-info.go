package configuration

type Stream struct {
	Collectors   []string `json:"collectors"`
	Distributors []string `json:"distributors"`
	Name         string   `json:"name"`
}

type LogicalServer struct {
	Name         string `json:"name"`
	Ip           string `json:"ip"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Port         string `json:"port"`
	ActiveStream string `json:"active_stream"`
}

type Cluster struct {
	Name           string          `json:"name"`
	LogicalServers []LogicalServer `json:"logical_servers"`
}

type EMMFileConfig struct {
	Streams  []Stream  `json:"streams"`
	Clusters []Cluster `json:"clusters"`
}
