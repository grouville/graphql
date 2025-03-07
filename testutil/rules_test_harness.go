package testutil

import (
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/language/location"
	"github.com/dagger/graphql/language/parser"
	"github.com/dagger/graphql/language/source"
)

var TestSchema *graphql.Schema

func init() {

	var beingInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Being",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "surname",
						Type: graphql.Boolean,
					},
				},
			},
		},
	})
	var petInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Pet",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "surname",
						Type: graphql.Boolean,
					},
				},
			},
		},
	})
	var canineInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Canine",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "surname",
						Type: graphql.Boolean,
					},
				},
			},
		},
	})
	var dogCommandEnum = graphql.NewEnum(graphql.EnumConfig{
		Name: "DogCommand",
		Values: graphql.EnumValueConfigMap{
			"SIT": &graphql.EnumValueConfig{
				Value: 0,
			},
			"HEEL": &graphql.EnumValueConfig{
				Value: 1,
			},
			"DOWN": &graphql.EnumValueConfig{
				Value: 2,
			},
		},
	})
	var dogType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Dog",
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			return true
		},
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "surname",
						Type: graphql.Boolean,
					},
				},
			},
			"nickname": &graphql.Field{
				Type: graphql.String,
			},
			"barkVolume": &graphql.Field{
				Type: graphql.Int,
			},
			"barks": &graphql.Field{
				Type: graphql.Boolean,
			},
			"doesKnowCommand": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "dogCommand",
						Type: dogCommandEnum,
					},
					&graphql.ArgumentConfig{
						Name: "nextDogCommand",
						Type: dogCommandEnum,
					},
				},
			},
			"isHousetrained": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name:         "atOtherHomes",
						Type:         graphql.Boolean,
						DefaultValue: true,
					},
				},
			},
			"isAtLocation": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "x",
						Type: graphql.Int,
					},
					&graphql.ArgumentConfig{
						Name: "y",
						Type: graphql.Int,
					},
				},
			},
		},
		Interfaces: []*graphql.Interface{
			beingInterface,
			petInterface,
			canineInterface,
		},
	})
	var furColorEnum = graphql.NewEnum(graphql.EnumConfig{
		Name: "FurColor",
		Values: graphql.EnumValueConfigMap{
			"BROWN": &graphql.EnumValueConfig{
				Value: 0,
			},
			"BLACK": &graphql.EnumValueConfig{
				Value: 1,
			},
			"TAN": &graphql.EnumValueConfig{
				Value: 2,
			},
			"SPOTTED": &graphql.EnumValueConfig{
				Value: 3,
			},
		},
	})

	var catType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Cat",
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			return true
		},
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "surname",
						Type: graphql.Boolean,
					},
				},
			},
			"nickname": &graphql.Field{
				Type: graphql.String,
			},
			"meowVolume": &graphql.Field{
				Type: graphql.Int,
			},
			"meows": &graphql.Field{
				Type: graphql.Boolean,
			},
			"furColor": &graphql.Field{
				Type: furColorEnum,
			},
		},
		Interfaces: []*graphql.Interface{
			beingInterface,
			petInterface,
		},
	})
	var catOrDogUnion = graphql.NewUnion(graphql.UnionConfig{
		Name: "CatOrDog",
		Types: []*graphql.Object{
			dogType,
			catType,
		},
	})
	var intelligentInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Intelligent",
		Fields: graphql.Fields{
			"iq": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	var humanType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Human",
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			return true
		},
		Interfaces: []*graphql.Interface{
			beingInterface,
			intelligentInterface,
		},
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "surname",
						Type: graphql.Boolean,
					},
				},
			},
			"pets": &graphql.Field{
				Type: graphql.NewList(petInterface),
			},
			"iq": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	humanType.AddFieldConfig("relatives", &graphql.Field{
		Type: graphql.NewList(humanType),
	})

	var alienType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Alien",
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			return true
		},
		Interfaces: []*graphql.Interface{
			beingInterface,
			intelligentInterface,
		},
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "surname",
						Type: graphql.Boolean,
					},
				},
			},
			"iq": &graphql.Field{
				Type: graphql.Int,
			},
			"numEyes": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})
	var dogOrHumanUnion = graphql.NewUnion(graphql.UnionConfig{
		Name: "DogOrHuman",
		Types: []*graphql.Object{
			dogType,
			humanType,
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			// not used for validation
			return nil
		},
	})
	var humanOrAlienUnion = graphql.NewUnion(graphql.UnionConfig{
		Name: "HumanOrAlien",
		Types: []*graphql.Object{
			alienType,
			humanType,
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			// not used for validation
			return nil
		},
	})

	var complexInputObject = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "ComplexInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"requiredField": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"intField": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"stringField": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"booleanField": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"stringListField": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.String),
			},
		},
	})
	var complicatedArgs = graphql.NewObject(graphql.ObjectConfig{
		Name: "ComplicatedArgs",
		// TODO List
		// TODO Coercion
		// TODO NotNulls
		Fields: graphql.Fields{
			"intArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "intArg",
						Type: graphql.Int,
					},
				},
			},
			"nonNullIntArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "nonNullIntArg",
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
			},
			"stringArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "stringArg",
						Type: graphql.String,
					},
				},
			},
			"booleanArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "booleanArg",
						Type: graphql.Boolean,
					},
				},
			},
			"enumArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "enumArg",
						Type: furColorEnum,
					},
				},
			},
			"floatArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "floatArg",
						Type: graphql.Float,
					},
				},
			},
			"idArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "idArg",
						Type: graphql.ID,
					},
				},
			},
			"stringListArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "stringListArg",
						Type: graphql.NewList(graphql.String),
					},
				},
			},
			"complexArgField": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "complexArg",
						Type: complexInputObject,
					},
				},
			},
			"multipleReqs": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "req1",
						Type: graphql.NewNonNull(graphql.Int),
					},
					&graphql.ArgumentConfig{
						Name: "req2",
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
			},
			"multipleOpts": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name:         "opt1",
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					&graphql.ArgumentConfig{
						Name:         "opt2",
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
			},
			"multipleOptAndReq": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "req1",
						Type: graphql.NewNonNull(graphql.Int),
					},
					&graphql.ArgumentConfig{
						Name: "req2",
						Type: graphql.NewNonNull(graphql.Int),
					},
					&graphql.ArgumentConfig{
						Name:         "opt1",
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					&graphql.ArgumentConfig{
						Name:         "opt2",
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
			},
		},
	})
	queryRoot := graphql.NewObject(graphql.ObjectConfig{
		Name: "QueryRoot",
		Fields: graphql.Fields{
			"human": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "id",
						Type: graphql.ID,
					},
				},
				Type: humanType,
			},
			"alien": &graphql.Field{
				Type: alienType,
			},
			"dog": &graphql.Field{
				Type: dogType,
			},
			"cat": &graphql.Field{
				Type: catType,
			},
			"pet": &graphql.Field{
				Type: petInterface,
			},
			"catOrDog": &graphql.Field{
				Type: catOrDogUnion,
			},
			"dogOrHuman": &graphql.Field{
				Type: dogOrHumanUnion,
			},
			"humanOrAlien": &graphql.Field{
				Type: humanOrAlienUnion,
			},
			"complicatedArgs": &graphql.Field{
				Type: complicatedArgs,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryRoot,
		Directives: []*graphql.Directive{
			graphql.IncludeDirective,
			graphql.SkipDirective,
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onQuery",
				Locations: []string{graphql.DirectiveLocationQuery},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onMutation",
				Locations: []string{graphql.DirectiveLocationMutation},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onSubscription",
				Locations: []string{graphql.DirectiveLocationSubscription},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onField",
				Locations: []string{graphql.DirectiveLocationField},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onFragmentDefinition",
				Locations: []string{graphql.DirectiveLocationFragmentDefinition},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onFragmentSpread",
				Locations: []string{graphql.DirectiveLocationFragmentSpread},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onInlineFragment",
				Locations: []string{graphql.DirectiveLocationInlineFragment},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onSchema",
				Locations: []string{graphql.DirectiveLocationSchema},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onScalar",
				Locations: []string{graphql.DirectiveLocationScalar},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onObject",
				Locations: []string{graphql.DirectiveLocationObject},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onFieldDefinition",
				Locations: []string{graphql.DirectiveLocationFieldDefinition},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onArgumentDefinition",
				Locations: []string{graphql.DirectiveLocationArgumentDefinition},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onInterface",
				Locations: []string{graphql.DirectiveLocationInterface},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onUnion",
				Locations: []string{graphql.DirectiveLocationUnion},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onEnum",
				Locations: []string{graphql.DirectiveLocationEnum},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onEnumValue",
				Locations: []string{graphql.DirectiveLocationEnumValue},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onInputObject",
				Locations: []string{graphql.DirectiveLocationInputObject},
			}),
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "onInputFieldDefinition",
				Locations: []string{graphql.DirectiveLocationInputFieldDefinition},
			}),
		},
		Types: []graphql.Type{
			catType,
			dogType,
			humanType,
			alienType,
		},
	})
	if err != nil {
		panic(err)
	}
	TestSchema = &schema

}
func expectValidRule(t *testing.T, schema *graphql.Schema, rules []graphql.ValidationRuleFn, queryString string) {
	source := source.NewSource(&source.Source{
		Body: []byte(queryString),
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.ValidateDocument(schema, AST, rules)
	if len(result.Errors) > 0 {
		t.Fatalf("Should validate, got %v", result.Errors)
	}
	if result.IsValid != true {
		t.Fatalf("IsValid should be true, got %v", result.IsValid)
	}

}
func expectInvalidRule(t *testing.T, schema *graphql.Schema, rules []graphql.ValidationRuleFn, queryString string, expectedErrors []gqlerrors.FormattedError) {
	source := source.NewSource(&source.Source{
		Body: []byte(queryString),
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		t.Fatal(err)
	}
	result := graphql.ValidateDocument(schema, AST, rules)
	if len(result.Errors) != len(expectedErrors) {
		t.Fatalf("Should have %v errors, got %v", len(expectedErrors), len(result.Errors))
	}
	if result.IsValid != false {
		t.Fatalf("IsValid should be false, got %v", result.IsValid)
	}
	for _, expectedErr := range expectedErrors {
		found := false
		for _, err := range result.Errors {
			if EqualFormattedError(expectedErr, err) {
				found = true
				break
			}
		}
		if found == false {
			t.Fatalf("Unexpected result, Diff: %v", Diff(expectedErrors, result.Errors))
		}
	}

}
func ExpectPassesRule(t *testing.T, rule graphql.ValidationRuleFn, queryString string) {
	expectValidRule(t, TestSchema, []graphql.ValidationRuleFn{rule}, queryString)
}
func ExpectFailsRule(t *testing.T, rule graphql.ValidationRuleFn, queryString string, expectedErrors []gqlerrors.FormattedError) {
	expectInvalidRule(t, TestSchema, []graphql.ValidationRuleFn{rule}, queryString, expectedErrors)
}
func ExpectFailsRuleWithSchema(t *testing.T, schema *graphql.Schema, rule graphql.ValidationRuleFn, queryString string, expectedErrors []gqlerrors.FormattedError) {
	expectInvalidRule(t, schema, []graphql.ValidationRuleFn{rule}, queryString, expectedErrors)
}
func ExpectPassesRuleWithSchema(t *testing.T, schema *graphql.Schema, rule graphql.ValidationRuleFn, queryString string) {
	expectValidRule(t, schema, []graphql.ValidationRuleFn{rule}, queryString)
}
func RuleError(message string, locs ...int) gqlerrors.FormattedError {
	locations := []location.SourceLocation{}
	for i := 0; i < len(locs); i += 2 {
		line := locs[i]
		col := 0
		if i+1 < len(locs) {
			col = locs[i+1]
		}
		locations = append(locations, location.SourceLocation{
			Line:   line,
			Column: col,
		})
	}
	return gqlerrors.FormattedError{
		Message:   message,
		Locations: locations,
	}
}
