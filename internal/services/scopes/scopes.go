package scopes

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type Scope struct {
	Attributes map[string]schema.Attribute
}

func New(attributes map[string]schema.Attribute) Scope {
	return Scope{
		Attributes: attributes,
	}
}

func (s *Scope) Subscope(attributes map[string]schema.Attribute) Scope {
	var out = make(map[string]schema.Attribute)
	maps.Copy(out, s.Attributes)
	maps.Copy(out, attributes)

	return Scope{
		Attributes: out,
	}
}

func (s *Scope) Build() map[string]schema.Attribute {
	return s.Attributes
}
