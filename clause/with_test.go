package clause_test

import (
	"fmt"
	"testing"

	"gorm.io/gorm/clause"
)

func TestWith(t *testing.T) {
	results := []struct {
		Clauses []clause.Interface
		Result  string
		Vars    []interface{}
	}{
		{
			[]clause.Interface{
				clause.With{
					CTEs: []clause.CTE{
						{
							Alias:       "cte",
							Expressions: []clause.Expression{clause.Expr{SQL: "SELECT * FROM `users`"}},
						},
					},
				},
				clause.Select{},
				clause.From{
					Tables: []clause.Table{{Name: "cte"}},
				},
			},
			"WITH `cte` AS (SELECT * FROM `users`) SELECT * FROM `cte`", nil,
		},
		{
			[]clause.Interface{
				clause.With{
					CTEs: []clause.CTE{
						{
							Alias:       "cte",
							Expressions: []clause.Expression{clause.Expr{SQL: "SELECT * FROM `users` WHERE `name` = ?", Vars: []interface{}{"jinzhu"}}},
						},
					},
				},
				clause.Select{},
				clause.From{
					Tables: []clause.Table{{Name: "cte"}},
				},
			},
			"WITH `cte` AS (SELECT * FROM `users` WHERE `name` = ?) SELECT * FROM `cte`", []interface{}{"jinzhu"},
		},
		{
			[]clause.Interface{
				clause.With{
					CTEs: []clause.CTE{
						{
							Alias:       "cte",
							Columns:     []string{"id", "name"},
							Expressions: []clause.Expression{clause.Expr{SQL: "SELECT `id`,`name` FROM `users` WHERE `name` = ?", Vars: []interface{}{"jinzhu"}}},
						},
					},
				},
				clause.Select{},
				clause.From{
					Tables: []clause.Table{{Name: "cte"}},
				},
			},
			"WITH `cte`(`id`,`name`) AS (SELECT `id`,`name` FROM `users` WHERE `name` = ?) SELECT * FROM `cte`", []interface{}{"jinzhu"},
		},
		{
			[]clause.Interface{
				clause.With{
					CTEs: []clause.CTE{
						{
							Recursive:   true,
							Alias:       "cte",
							Expressions: []clause.Expression{clause.Expr{SQL: "SELECT * FROM `users`"}},
						},
					},
				},
				clause.Select{},
				clause.From{
					Tables: []clause.Table{{Name: "cte"}},
				},
			},
			"WITH RECURSIVE `cte` AS (SELECT * FROM `users`) SELECT * FROM `cte`", nil,
		},
		{
			[]clause.Interface{
				clause.With{
					CTEs: []clause.CTE{
						{
							Recursive:   true,
							Alias:       "cte1",
							Expressions: []clause.Expression{clause.Expr{SQL: "SELECT * FROM `users`"}},
						},
					},
				},
				clause.With{
					CTEs: []clause.CTE{
						{
							Recursive:   true,
							Alias:       "cte2",
							Expressions: []clause.Expression{clause.Expr{SQL: "SELECT * FROM `users`"}},
						},
					},
				},
				clause.With{
					CTEs: []clause.CTE{
						{
							Alias:       "cte3",
							Expressions: []clause.Expression{clause.Expr{SQL: "SELECT * FROM `users`"}},
						},
					},
				},
				clause.Select{},
				clause.From{
					Tables: []clause.Table{{Name: "cte1"}, {Name: "cte2"}, {Name: "cte3"}},
				},
			},
			"WITH RECURSIVE `cte1` AS (SELECT * FROM `users`),`cte2` AS (SELECT * FROM `users`),`cte3` AS (SELECT * FROM `users`) SELECT * FROM `cte1`,`cte2`,`cte3`", nil,
		},
	}
	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			checkBuildClauses(t, result.Clauses, result.Result, result.Vars)
		})
	}
}
