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
			Type: graphql.String,
		},
		"gitBranch": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var serviceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Service",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"ip": &graphql.Field{
			Type: graphql.String,
		},
		"port": &graphql.Field{
			Type: graphql.Int,
		},
		"state": &graphql.Field{
			Type: graphql.String,
		},
		"config": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var nodeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Node",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"dataCenter": &graphql.Field{
			Type: graphql.String,
		},
		"ip": &graphql.Field{
			Type: graphql.String,
		},
		"port": &graphql.Field{
			Type: graphql.Int,
		},
		"state": &graphql.Field{
			Type: graphql.String,
		},
		"version": &graphql.Field{
			Type: graphql.Int,
		},
		"snapShot": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func init() {
	nodeType.AddFieldConfig("services", &graphql.Field{
		Type:        graphql.NewList(serviceType),
		Description: "The services of the node, or an empty list if they have none.",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if n, ok := p.Source.(*api.Node); ok {
				name := n.Name
				return server.api.Node(name).Services, nil
			} else {
				return nil, nil
			}
		},
	})
	serviceType.AddFieldConfig("nodes", &graphql.Field{
		Type:        graphql.NewList(nodeType),
		Description: "The nodes of the service, or an empty list if they have none.",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if s, ok := p.Source.(*api.Service); ok {
				name := s.Name
				return server.api.Service(name, api.StateAll).Nodes, nil
			} else {
				return nil, nil
			}
		},
	})
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"riff": &graphql.Field{
			Type:        riffType,
			Description: "riff",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return map[string]interface{}{
					"version":   common.Version,
					"gitSha":    common.GitSha,
					"gitBranch": common.GitBranch,
				}, nil
			},
		},
		"server": &graphql.Field{
			Type: nodeType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name := server.Self.Name
				return server.api.Node(name), nil
			},
		},
		"node": &graphql.Field{
			Type: nodeType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "node",
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
		"nodes": &graphql.Field{
			Type:        graphql.NewList(nodeType),
			Description: "List of node",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return server.api.Nodes(), nil
			},
		},
		"services": &graphql.Field{
			Type:        graphql.NewList(serviceType),
			Description: "List of service",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return server.api.Services(), nil
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
