package controllers

import "myapp/internal/logic"

type HandlerLayer struct {
	*logic.LogicLayer
}

func NewHandlerLayer(ll *logic.LogicLayer) *HandlerLayer {
	return &HandlerLayer{
		LogicLayer: ll,
	}
}
