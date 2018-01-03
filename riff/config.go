package riff

type Config struct {
	Id         string     `yaml:"id"`
	IP         string     `yaml:"ip"`
	Name       string     `yaml:"name"`
	DataCenter string     `yaml:"data_center"`
	Addresses  *Addresses `yaml:"addresses"`
	Ports      *Ports     `yaml:"ports"`
}
type Addresses struct {
	Http string `yaml:"http"`
	Dns  string `yaml:"dns"`
	Rpc  string `yaml:"rpc"`
}

type Ports struct {
	Http int `yaml:"http"`
	Dns  int `yaml:"dns"`
	Rpc  int `yaml:"rpc"`
}

//func NewConfig(rpc, name, dataCenter string) (*Config, error) {
//	//make default config
//	addresses := &Addresses{
//		Http: "127.0.0.1",
//		Dns:  "127.0.0.1",
//		Rpc:  rpc,
//	}
//	ports := &Ports{
//		Http: 8610,
//		Dns:  8620,
//		Rpc:  8630,
//	}
//	return &Config{
//		Id:         "",
//		IP:         "",
//		Addresses:  addresses,
//		Ports:      ports,
//		Name:       name,
//		DataCenter: dataCenter,
//	}, nil
//}
