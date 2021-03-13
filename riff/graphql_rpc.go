package riff

import (
	"github.com/teatak/riff/api"
)

func mutationService(name, address string, cmd api.CmdType) error {
	client, err := api.NewClient(address)
	if err != nil {
		return err
	}
	defer client.Close()

	var result bool
	err = client.Call("Mutation.Service", api.ParamServiceMutation{
		Name: name,
		Cmd:  cmd,
	}, &result)
	if err != nil {
		server.Logger.Printf(errorServerPrefix+"%v\n", err)
		return err
	}
	return nil
}

func registeService(address string, config string) error {
	client, err := api.NewClient(address)
	if err != nil {
		return err
	}
	defer client.Close()

	var result bool
	err = client.Call("Mutation.RegisteService", config, &result)
	if err != nil {
		server.Logger.Printf(errorServerPrefix+"%v\n", err)
		return err
	}
	return nil
}

func unregisteService(address string, name string) error {
	client, err := api.NewClient(address)
	if err != nil {
		return err
	}
	defer client.Close()

	var result bool
	err = client.Call("Mutation.UnregisteService", name, &result)
	if err != nil {
		server.Logger.Printf(errorServerPrefix+"%v\n", err)
		return err
	}
	return nil
}
