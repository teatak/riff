package riff

type Config struct {
	Id         string     `yaml:"id"`          //id
	IP         string     `yaml:"ip"`          //server ip
	Name       string     `yaml:"name"`        //server random name
	DataCenter string     `yaml:"data_center"` //server data center
	Join       string     `yaml:"join"`        //join address
	AutoPilot  bool       `yaml:"auto_pilot"`  //auto join node
	Addresses  *Addresses `yaml:"addresses"`
	Ports      *Ports     `yaml:"ports"`
	Fanout     int        `yaml:"fan_out"`
}
type Addresses struct {
	Http string `yaml:"http"` //http address
	Dns  string `yaml:"dns"`  //dns address
	Rpc  string `yaml:"rpc"`  //rpc address
}

type Ports struct {
	Http int `yaml:"http"` //http port default 8610
	Dns  int `yaml:"dns"`  //dns port default 8620
	Rpc  int `yaml:"rpc"`  //rpc port defalut 8630
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
