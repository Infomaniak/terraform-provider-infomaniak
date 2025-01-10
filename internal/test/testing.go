package test

import (
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func MustGetTestFile(group, subgroup, file string) string {
	body, err := os.ReadFile(path.Join("tests", group, subgroup, file))
	if err != nil {
		panic(err)
	}

	return string(body)
}

func TestCheckGetFromState(name string, key string, out *string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		ms := s.RootModule()

		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", name, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", name, ms.Path)
		}

		val, ok := is.Attributes[key]

		if ok && val != "" {
			*out = val
			return nil
		}

		if _, ok := is.Attributes[key+".#"]; ok {
			return fmt.Errorf(
				"%s: list or set attribute '%s' must be checked by element count key (%s) or element value keys (e.g. %s). Set element value checks should use TestCheckTypeSet functions instead.",
				name,
				key,
				key+".#",
				key+".0",
			)
		}

		if _, ok := is.Attributes[key+".%"]; ok {
			return fmt.Errorf(
				"%s: map attribute '%s' must be checked by element count key (%s) or element value keys (e.g. %s).",
				name,
				key,
				key+".%",
				key+".examplekey",
			)
		}

		return fmt.Errorf("%s: Attribute '%s' expected to be set", name, key)
	}
}
