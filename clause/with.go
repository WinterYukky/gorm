package clause

// With with clause
type With struct {
	CTEs []CTE
}

// CTE common table expressions
type CTE struct {
	Recursive   bool
	Alias       string
	Columns     []string
	Expressions []Expression
}

// Name with clause name
func (with With) Name() string {
	return "WITH"
}

// Build build with clause
func (with With) Build(builder Builder) {
	for _, cte := range with.CTEs {
		if cte.Recursive {
			builder.WriteString("RECURSIVE ")
			break
		}
	}
	for index, cte := range with.CTEs {
		if index > 0 {
			builder.WriteByte(',')
		}
		cte.Build(builder)
	}
}

// Build build CTE
func (cte CTE) Build(builder Builder) {
	builder.WriteQuoted(cte.Alias)
	if len(cte.Columns) > 0 {
		builder.WriteByte('(')
		for index, column := range cte.Columns {
			if index > 0 {
				builder.WriteByte(',')
			}
			builder.WriteQuoted(column)
		}
		builder.WriteByte(')')
	}

	builder.WriteString(" AS ")

	builder.WriteByte('(')
	for _, expression := range cte.Expressions {
		expression.Build(builder)
	}
	builder.WriteByte(')')
}

// MergeClause merge with clauses
func (with With) MergeClause(clause *Clause) {
	if w, ok := clause.Expression.(With); ok {
		ctes := make([]CTE, len(w.CTEs)+len(with.CTEs))
		copy(ctes, w.CTEs)
		copy(ctes[len(w.CTEs):], with.CTEs)
		with.CTEs = ctes
	}

	clause.Expression = with
}
