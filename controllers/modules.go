package controllers

import "github.com/goava/di"

var Options = di.Options(
	di.ProvideValue(PeopleApiController{}, di.As(new(ApiControler))),
)
