package graphql_test

import (
	"errors"
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/testutil"
)

var directivesTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "TestType",
		Fields: graphql.Fields{
			"a": &graphql.Field{
				Type: graphql.String,
			},
			"b": &graphql.Field{
				Type: graphql.String,
			},
		},
	}),
})

var directivesTestData map[string]interface{} = map[string]interface{}{
	"a": func() interface{} { return "a" },
	"b": func() interface{} { return "b" },
}

func executeDirectivesTestQuery(t *testing.T, doc string) *graphql.Result {
	ast := testutil.TestParse(t, doc)
	ep := graphql.ExecuteParams{
		Schema: directivesTestSchema,
		AST:    ast,
		Root:   directivesTestData,
	}
	return testutil.TestExecute(t, ep)
}

func TestDirectives_DirectivesMustBeNamed(t *testing.T) {
	invalidDirective := graphql.NewDirective(graphql.DirectiveConfig{
		Locations: []string{
			graphql.DirectiveLocationField,
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "TestType",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Directives: []*graphql.Directive{invalidDirective},
	})
	actualErr := gqlerrors.FormatError(err)
	expectedErr := gqlerrors.FormatError(errors.New("Directive must be named."))
	if !testutil.EqualFormattedError(expectedErr, actualErr) {
		t.Fatalf("Expected error to be equal, got: %v", testutil.Diff(expectedErr, actualErr))
	}
}

func TestDirectives_DirectiveNameMustBeValid(t *testing.T) {
	invalidDirective := graphql.NewDirective(graphql.DirectiveConfig{
		Name: "123invalid name",
		Locations: []string{
			graphql.DirectiveLocationField,
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "TestType",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Directives: []*graphql.Directive{invalidDirective},
	})
	actualErr := gqlerrors.FormatError(err)
	expectedErr := gqlerrors.FormatError(errors.New(`Names must match /^[_a-zA-Z][_a-zA-Z0-9]*$/ but "123invalid name" does not.`))
	if !testutil.EqualFormattedError(expectedErr, actualErr) {
		t.Fatalf("Expected error to be equal, got: %v", testutil.Diff(expectedErr, actualErr))
	}
}

func TestDirectives_DirectiveNameMustProvideLocations(t *testing.T) {
	invalidDirective := graphql.NewDirective(graphql.DirectiveConfig{
		Name: "skip",
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "TestType",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Directives: []*graphql.Directive{invalidDirective},
	})
	actualErr := gqlerrors.FormatError(err)
	expectedErr := gqlerrors.FormatError(errors.New(`Must provide locations for directive.`))
	if !testutil.EqualFormattedError(expectedErr, actualErr) {
		t.Fatalf("Expected error to be equal, got: %v", testutil.Diff(expectedErr, actualErr))
	}
}

func TestDirectives_DirectiveArgNamesMustBeValid(t *testing.T) {
	invalidDirective := graphql.NewDirective(graphql.DirectiveConfig{
		Name: "skip",
		Description: "Directs the executor to skip this field or fragment when the `if` " +
			"argument is true.",
		Args: graphql.FieldConfigArgument{
			&graphql.ArgumentConfig{
				Name:        "123if",
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Skipped when true.",
			},
		},
		Locations: []string{
			graphql.DirectiveLocationField,
			graphql.DirectiveLocationFragmentSpread,
			graphql.DirectiveLocationInlineFragment,
		},
	})
	_, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "TestType",
			Fields: graphql.Fields{
				"a": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
		Directives: []*graphql.Directive{invalidDirective},
	})
	actualErr := gqlerrors.FormatError(err)
	expectedErr := gqlerrors.FormatError(errors.New(`Names must match /^[_a-zA-Z][_a-zA-Z0-9]*$/ but "123if" does not.`))
	if !testutil.EqualFormattedError(expectedErr, actualErr) {
		t.Fatalf("Expected error to be equal, got: %v", testutil.Diff(expectedErr, actualErr))
	}
}

func TestDirectivesWorksWithoutDirectives(t *testing.T) {
	query := `{ a, b }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsIfTrueIncludesScalar(t *testing.T) {
	query := `{ a, b @include(if: true) }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsIfFalseOmitsOnScalar(t *testing.T) {
	query := `{ a, b @include(if: false) }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsUnlessFalseIncludesScalar(t *testing.T) {
	query := `{ a, b @skip(if: false) }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnScalarsUnlessTrueOmitsScalar(t *testing.T) {
	query := `{ a, b @skip(if: true) }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentSpreadsIfFalseOmitsFragmentSpread(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag @include(if: false)
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentSpreadsIfTrueIncludesFragmentSpread(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag @include(if: true)
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentSpreadsUnlessFalseIncludesFragmentSpread(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag @skip(if: false)
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnFragmentSpreadsUnlessTrueOmitsFragmentSpread(t *testing.T) {
	query := `
        query Q {
          a
          ...Frag @skip(if: true)
        }
        fragment Frag on TestType {
          b
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnInlineFragmentIfFalseOmitsInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... on TestType @include(if: false) {
            b
          }
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnInlineFragmentIfTrueIncludesInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... on TestType @include(if: true) {
            b
          }
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnInlineFragmentUnlessFalseIncludesInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... on TestType @skip(if: false) {
            b
          }
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnInlineFragmentUnlessTrueIncludesInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... on TestType @skip(if: true) {
            b
          }
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnAnonymousInlineFragmentIfFalseOmitsAnonymousInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... @include(if: false) {
            b
          }
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnAnonymousInlineFragmentIfTrueIncludesAnonymousInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... @include(if: true) {
            b
          }
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnAnonymousInlineFragmentUnlessFalseIncludesAnonymousInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... @skip(if: false) {
            b
          }
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksOnAnonymousInlineFragmentUnlessTrueIncludesAnonymousInlineFragment(t *testing.T) {
	query := `
        query Q {
          a
          ... @skip(if: true) {
            b
          }
        }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksWithSkipAndIncludeDirectives_IncludeAndNoSkip(t *testing.T) {
	query := `{ a, b @include(if: true) @skip(if: false) }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
			"b": "b",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksWithSkipAndIncludeDirectives_IncludeAndSkip(t *testing.T) {
	query := `{ a, b @include(if: true) @skip(if: true) }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksWithSkipAndIncludeDirectives_NoIncludeAndSkip(t *testing.T) {
	query := `{ a, b @include(if: false) @skip(if: true) }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestDirectivesWorksWithSkipAndIncludeDirectives_NoIncludeOrSkip(t *testing.T) {
	query := `{ a, b @include(if: false) @skip(if: false) }`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"a": "a",
		},
	}
	result := executeDirectivesTestQuery(t, query)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
