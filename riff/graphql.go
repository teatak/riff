package riff

import (
	"fmt"
	"github.com/gimke/riff/api"
	"github.com/gimke/riff/common"
	"github.com/graphql-go/graphql"
)

var riffType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Riff",
	Fields: graphql.Fields{
		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "current version",
		},
		"gitSha": &graphql.Field{
			Type:        graphql.String,
			Description: "current git sha",
		},
		"gitBranch": &graphql.Field{
			Type:        graphql.String,
			Description: "current git branch",
		},
	},
})

var serviceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Service",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "name of service",
		},
	},
})

var nestServiceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NestService",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "name of service",
		},
		"ip": &graphql.Field{
			Type:        graphql.String,
			Description: "ip of service",
		},
		"port": &graphql.Field{
			Type:        graphql.Int,
			Description: "port of service",
		},
		"state": &graphql.Field{
			Type:        graphql.String,
			Description: "state of service (Alive,Suspect,Dead)",
		},
		"config": &graphql.Field{
			Type:        graphql.String,
			Description: "config of service",
		},
	},
})

var nodeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Node",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "name of node",
		},
		"dataCenter": &graphql.Field{
			Type:        graphql.String,
			Description: "data center of node",
		},
		"ip": &graphql.Field{
			Type:        graphql.String,
			Description: "ip of node",
		},
		"port": &graphql.Field{
			Type:        graphql.Int,
			Description: "port of node",
		},
		"state": &graphql.Field{
			Type:        graphql.String,
			Description: "state of node (Alive,Suspect,Dead)",
		},
		"version": &graphql.Field{
			Type:        graphql.Int,
			Description: "version of node",
		},
		"snapShot": &graphql.Field{
			Type:        graphql.String,
			Description: "snapshot of node",
		},
	},
})

var nestNodeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NestNode",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "name of node",
		},
		"dataCenter": &graphql.Field{
			Type:        graphql.String,
			Description: "data center of node",
		},
		"ip": &graphql.Field{
			Type:        graphql.String,
			Description: "ip of node",
		},
		"port": &graphql.Field{
			Type:        graphql.Int,
			Description: "port of node",
		},
		"state": &graphql.Field{
			Type:        graphql.String,
			Description: "state of node (Alive,Suspect,Dead)",
		},
		"version": &graphql.Field{
			Type:        graphql.Int,
			Description: "version of node",
		},
		"snapShot": &graphql.Field{
			Type:        graphql.String,
			Description: "snapshot of node",
		},
		"config": &graphql.Field{
			Type:        graphql.String,
			Description: "config of node",
		},
	},
})

func init() {
	nodeType.AddFieldConfig("services", &graphql.Field{
		Type:        graphql.NewList(nestServiceType),
		Description: "the services of the node, or an empty list if they have none.",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if n, ok := p.Source.(*api.Node); ok {
				name := n.Name
				if n.NestServices != nil {
					return n.NestServices, nil
				} else {
					return server.api.Node(name).NestServices, nil
				}
			} else {
				return nil, nil
			}
		},
	})
	serviceType.AddFieldConfig("nodes", &graphql.Field{
		Type:        graphql.NewList(nestNodeType),
		Description: "the nodes of the service, or an empty list if they have none.",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if s, ok := p.Source.(*api.Service); ok {
				if s.NestNodes != nil {
					return s.NestNodes, nil
				} else {
					name := s.Name
					return server.api.Service(name, api.StateAll).NestNodes, nil
				}
			} else {
				return nil, nil
			}
		},
	})
}

var enumStateType = graphql.NewEnum(graphql.EnumConfig{
	Name: "State",
	Values: graphql.EnumValueConfigMap{
		"Alive": &graphql.EnumValueConfig{
			Value: api.StateAlive,
		},
		"Suspect": &graphql.EnumValueConfig{
			Value: api.StateSuspect,
		},
		"Dead": &graphql.EnumValueConfig{
			Value: api.StateDead,
		},
		"All": &graphql.EnumValueConfig{
			Value: api.StateAll,
		},
	},
})
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"nestNodeType": &graphql.Field{
			Type:        nestNodeType,
			Description: "only export",
		},
		"nestServiceType": &graphql.Field{
			Type:        nestServiceType,
			Description: "only export",
		},
		"riff": &graphql.Field{
			Type:        riffType,
			Description: "get riff version, git sha or brance",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return map[string]interface{}{
					"version":   common.Version,
					"gitSha":    common.GitSha,
					"gitBranch": common.GitBranch,
				}, nil
			},
		},
		"server": &graphql.Field{
			Type:        nodeType,
			Description: "get current node",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name := server.Self.Name
				return server.api.Node(name), nil
			},
		},
		"nodes": &graphql.Field{
			Type:        graphql.NewList(nodeType),
			Description: "list of node",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return server.api.Nodes(), nil
			},
		},
		"node": &graphql.Field{
			Type: nodeType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "get node",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, ok := p.Args["name"].(string)
				if ok {
					if n := server.api.Node(name); n != nil {
						return n, nil
					} else {
						return nil, fmt.Errorf("NOT_FOUND")
					}
				}
				return nil, nil
			},
		},
		"services": &graphql.Field{
			Type:        graphql.NewList(serviceType),
			Description: "list of service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return server.api.Services(), nil
			},
		},
		"service": &graphql.Field{
			Type: serviceType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"state": &graphql.ArgumentConfig{
					Type: enumStateType,
				},
			},
			Description: "get service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, ok := p.Args["name"].(string)
				state, stateOk := p.Args["state"]
				if ok && stateOk {
					stateType := state.(api.StateType)
					if n := server.api.Service(name, stateType); n != nil {
						return n, nil
					} else {
						return nil, fmt.Errorf("NOT_FOUND")
					}
				}
				return nil, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
	//Mutation: rootMutation,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		server.Logger.Printf(errorServicePrefix+"wrong result, unexpected errors: %v\n", result.Errors)
	}
	return result
}
