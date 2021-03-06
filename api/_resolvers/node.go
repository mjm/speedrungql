package resolvers

import (
	"context"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"
)

type nodeResolver interface {
	ID() graphql.ID
}

func (r *Resolvers) Node(ctx context.Context, args struct{ ID graphql.ID }) (*Node, error) {
	kind := relay.UnmarshalKind(args.ID)
	var id string
	if err := relay.UnmarshalSpec(args.ID, &id); err != nil {
		return nil, err
	}

	var n nodeResolver

	switch kind {
	case "category":
		cat, err := r.client.GetCategory(ctx, id)
		if err != nil {
			return nil, err
		}
		if cat != nil {
			n = &Category{*cat, r.client}
		}
	case "engine":
		eng, err := r.client.GetEngine(ctx, id)
		if err != nil {
			return nil, err
		}
		if eng != nil {
			n = &Engine{*eng, r.client}
		}
	case "game":
		game, err := r.client.GetGame(ctx, id)
		if err != nil {
			return nil, err
		}
		if game != nil {
			n = &Game{*game, r.client}
		}
	case "genre":
		genre, err := r.client.GetGenre(ctx, id)
		if err != nil {
			return nil, err
		}
		if genre != nil {
			n = &Genre{*genre, r.client}
		}
	case "level":
		level, err := r.client.GetLevel(ctx, id)
		if err != nil {
			return nil, err
		}
		if level != nil {
			n = &Level{*level, r.client}
		}
	case "platform":
		plat, err := r.client.GetPlatform(ctx, id)
		if err != nil {
			return nil, err
		}
		if plat != nil {
			n = &Platform{*plat, r.client}
		}
	case "region":
		reg, err := r.client.GetRegion(ctx, id)
		if err != nil {
			return nil, err
		}
		if reg != nil {
			n = &Region{*reg, r.client}
		}
	case "run":
		run, err := r.client.GetRun(ctx, id)
		if err != nil {
			return nil, err
		}
		if run != nil {
			n = &Run{*run, r.client}
		}
	case "user":
		user, err := r.client.GetUser(ctx, id)
		if err != nil {
			return nil, err
		}
		if user != nil {
			n = &User{*user, r.client}
		}
	case "variable":
		v, err := r.client.GetVariable(ctx, id)
		if err != nil {
			return nil, err
		}
		if v != nil {
			n = &Variable{*v, r.client}
		}
	}

	if n == nil {
		return nil, nil
	}

	return &Node{n}, nil
}

type Node struct {
	nodeResolver
}

func (n *Node) ToCategory() (*Category, bool) {
	c, ok := n.nodeResolver.(*Category)
	return c, ok
}

func (n *Node) ToEngine() (*Engine, bool) {
	e, ok := n.nodeResolver.(*Engine)
	return e, ok
}

func (n *Node) ToGame() (*Game, bool) {
	g, ok := n.nodeResolver.(*Game)
	return g, ok
}

func (n *Node) ToGenre() (*Genre, bool) {
	g, ok := n.nodeResolver.(*Genre)
	return g, ok
}

func (n *Node) ToLevel() (*Level, bool) {
	l, ok := n.nodeResolver.(*Level)
	return l, ok
}

func (n *Node) ToPlatform() (*Platform, bool) {
	p, ok := n.nodeResolver.(*Platform)
	return p, ok
}

func (n *Node) ToRegion() (*Region, bool) {
	r, ok := n.nodeResolver.(*Region)
	return r, ok
}

func (n *Node) ToRun() (*Run, bool) {
	r, ok := n.nodeResolver.(*Run)
	return r, ok
}

func (n *Node) ToUser() (*User, bool) {
	u, ok := n.nodeResolver.(*User)
	return u, ok
}

func (n *Node) ToVariable() (*Variable, bool) {
	v, ok := n.nodeResolver.(*Variable)
	return v, ok
}
