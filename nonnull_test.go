package graphql_test

import (
	"sort"
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/language/location"
	"github.com/dagger/graphql/testutil"
)

var syncError = "sync"
var nonNullSyncError = "nonNullSync"
var promiseError = "promise"
var nonNullPromiseError = "nonNullPromise"

var throwingData = map[string]interface{}{
	"sync": func() interface{} {
		panic(syncError)
	},
	"nonNullSync": func() interface{} {
		panic(nonNullSyncError)
	},
	"promise": func() interface{} {
		panic(promiseError)
	},
	"nonNullPromise": func() interface{} {
		panic(nonNullPromiseError)
	},
}

var nullingData = map[string]interface{}{
	"sync": func() interface{} {
		return nil
	},
	"nonNullSync": func() interface{} {
		return nil
	},
	"promise": func() interface{} {
		return nil
	},
	"nonNullPromise": func() interface{} {
		return nil
	},
}

var dataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DataType",
	Fields: graphql.Fields{
		"sync": &graphql.Field{
			Type: graphql.String,
		},
		"nonNullSync": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"promise": &graphql.Field{
			Type: graphql.String,
		},
		"nonNullPromise": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var nonNullTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: dataType,
})

func init() {
	throwingData["nest"] = func() interface{} {
		return throwingData
	}
	throwingData["nonNullNest"] = func() interface{} {
		return throwingData
	}
	throwingData["promiseNest"] = func() interface{} {
		return throwingData
	}
	throwingData["nonNullPromiseNest"] = func() interface{} {
		return throwingData
	}

	nullingData["nest"] = func() interface{} {
		return nullingData
	}
	nullingData["nonNullNest"] = func() interface{} {
		return nullingData
	}
	nullingData["promiseNest"] = func() interface{} {
		return nullingData
	}
	nullingData["nonNullPromiseNest"] = func() interface{} {
		return nullingData
	}

	dataType.AddFieldConfig("nest", &graphql.Field{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullNest", &graphql.Field{
		Type: graphql.NewNonNull(dataType),
	})
	dataType.AddFieldConfig("promiseNest", &graphql.Field{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullPromiseNest", &graphql.Field{
		Type: graphql.NewNonNull(dataType),
	})
}

// nulls a nullable field that panics
func TestNonNull_NullsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sync": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: syncError,
				Locations: []location.SourceLocation{
					{
						Line: 3, Column: 9,
					},
				},
				Path: []interface{}{
					"sync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsANullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promise": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{
						Line: 3, Column: 9,
					},
				},
				Path: []interface{}{
					"promise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []interface{}{
					"nest",
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []interface{}{
					"nest",
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []interface{}{
					"promiseNest",
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []interface{}{
					"promiseNest",
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsAComplexTreeOfNullableFieldsThatThrow(t *testing.T) {
	doc := `
      query Q {
        nest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
        promiseNest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
			"promiseNest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []interface{}{
					"nest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 7, Column: 13},
				},
				Path: []interface{}{
					"nest", "nest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 11, Column: 13},
				},
				Path: []interface{}{
					"nest", "promiseNest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 16, Column: 11},
				},
				Path: []interface{}{
					"promiseNest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 19, Column: 13},
				},
				Path: []interface{}{
					"promiseNest", "nest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 23, Column: 13},
				},
				Path: []interface{}{
					"promiseNest", "promiseNest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 5, Column: 11},
				},
				Path: []interface{}{
					"nest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 8, Column: 13},
				},
				Path: []interface{}{
					"nest", "nest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 12, Column: 13},
				},
				Path: []interface{}{
					"nest", "promiseNest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 17, Column: 11},
				},
				Path: []interface{}{
					"promiseNest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 20, Column: 13},
				},
				Path: []interface{}{
					"promiseNest", "nest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 24, Column: 13},
				},
				Path: []interface{}{
					"promiseNest", "promiseNest", "promise",
				},
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsTheFirstNullableObjectAfterAFieldThrowsInALongChainOfFieldsThatAreNonNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        anotherNest: nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
        anotherPromiseNest: promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{Line: 8, Column: 19},
				},
				Path: []interface{}{
					"nest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{Line: 19, Column: 19},
				},
				Path: []interface{}{
					"promiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{Line: 30, Column: 19},
				},
				Path: []interface{}{
					"anotherNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{Line: 41, Column: 19},
				},
				Path: []interface{}{
					"anotherPromiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}

}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNull(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sync": nil,
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promise": nil,
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatReturnsNullSynchronously(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []interface{}{
					"nest",
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []interface{}{
					"nest",
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatReturnsNullSynchronously(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []interface{}{
					"promiseNest",
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []interface{}{
					"promiseNest",
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAComplexTreeOfNullableFieldsThatReturnNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
        promiseNest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
			"promiseNest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
}
func TestNonNull_NullsTheFirstNullableObjectAfterAFieldReturnsNullInALongChainOfFieldsThatAreNonNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        anotherNest: nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
        anotherPromiseNest: promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 8, Column: 19},
				},
				Path: []interface{}{
					"nest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
			},
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 19, Column: 19},
				},
				Path: []interface{}{
					"promiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
			},
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 30, Column: 19},
				},
				Path: []interface{}{
					"anotherNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
			},
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 41, Column: 19},
				},
				Path: []interface{}{
					"anotherPromiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldThrows(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []interface{}{
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldErrors(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []interface{}{
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldReturnsNull(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []interface{}{
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldResolvesNull(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []interface{}{
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
