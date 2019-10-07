package schema

import "errors"

var (
	gitHub = &Schema{
		Name: "GitHub",
		Blueprint: []*Table{
			&Table{
				Name: "users",
				Columns: []*Column{
					&Column{Name: "asd", Validator: isString},
				},
			},
		},
	}
)

func isString(val interface{}) error {
	if _, ok := val.(string); ok {
		return nil
	}
	return errors.New("invalid value type")
}
