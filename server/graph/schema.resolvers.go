package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.39

import (
	"context"
	"graphql_json_go/graph/gentypes"
	"graphql_json_go/graph/model"
)

// CreateModel is the resolver for the createModel field.
func (r *mutationResolver) CreateModel(ctx context.Context, name string, fields []*model.ModelFieldInput) (*model.Model, error) {
	return CreateModel(ctx, name, fields)
}

// GetModels is the resolver for the getModels field.
func (r *queryResolver) GetModels(ctx context.Context) ([]*model.Model, error) {
	return GetModels(ctx)
}

// SendInformaton is the resolver for the sendInformaton field.
func (r *queryResolver) SendInformaton(ctx context.Context, info map[string]interface{}, modelName string) (bool, error) {
	return SendInformaton(ctx, info, modelName)
}

// Mutation returns gentypes.MutationResolver implementation.
func (r *Resolver) Mutation() gentypes.MutationResolver { return &mutationResolver{r} }

// Query returns gentypes.QueryResolver implementation.
func (r *Resolver) Query() gentypes.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
