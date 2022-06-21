package components

import "github.com/tomknightdev/dwarven-fortresses/enums"

type Resource struct {
	enums.ResourceTypeEnum
}

func NewResource(rt enums.ResourceTypeEnum) Resource {
	return Resource{
		ResourceTypeEnum: rt,
	}
}
