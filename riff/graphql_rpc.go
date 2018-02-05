package riff

import (
	"net"
	"time"
	"github.com/gimke/riff/api"
	"net/rpc"
)

func mutationService(name,address string,cmd api.CmdType) error {
	conn, err := net.DialTimeout("tcp", address, time.Second*10)
	if err != nil {
		server.Logger.Printf(errorServerPrefix+"%v\n",err)
		return err
	}
	codec := api.NewGobClientCodec(conn)
	client := rpc.NewClientWithCodec(codec)
	defer client.Close()

	var result bool
	err = client.Call("Mutation.Service", api.ParamServiceMutation{
		Name: name,
		Cmd:  cmd,
	}, &result)
	if err != nil {
		server.Logger.Printf(errorServerPrefix+"%v\n",err)
		return err
	}
	return nil
}
