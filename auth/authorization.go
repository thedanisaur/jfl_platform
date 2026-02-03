package auth

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"

	"github.com/google/cel-go/cel"
	"github.com/google/uuid"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func compileCelToSQL(expr string, scope map[string]string, request_user types.UserClaims) (string, []interface{}, error) {
	// Parse CEL expression into AST
	// TODO [drd] pull this out of here as the env will need to be built in each app
	env, err := cel.NewEnv(
		cel.Variable("log", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("aircrew", cel.ListType(cel.MapType(cel.StringType, cel.DynType))),
		cel.Variable("request_user", cel.MapType(cel.StringType, cel.DynType)),
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create new env: %w", err)
	}
	ast, iss := env.Parse(expr)
	if iss.Err() != nil {
		return "", nil, fmt.Errorf("failed to parse the expression: %w", iss.Err())
	}
	checked, iss := env.Check(ast)
	if iss.Err() != nil {
		return "", nil, fmt.Errorf("env cannot accept checked ast: %w", iss.Err())
	}
	checked_expression, err := cel.AstToCheckedExpr(checked)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert ast to checked expression: %w", err)
	}

	// Convert AST to SQL
	return compileExpression(checked_expression.GetExpr(), scope, request_user)
}

func compileExpression(expression *exprpb.Expr, scope map[string]string, request_user types.UserClaims) (string, []interface{}, error) {
	switch expression_kind := expression.ExprKind.(type) {

	// boolean constants
	case *exprpb.Expr_ConstExpr:
		switch v := expression_kind.ConstExpr.ConstantKind.(type) {

		case *exprpb.Constant_BoolValue:
			if v.BoolValue {
				return "1=1", nil, nil
			} else {
				return "1=0", nil, nil
			}

		case *exprpb.Constant_Int64Value:
			return "?", []interface{}{v.Int64Value}, nil

		case *exprpb.Constant_DoubleValue:
			return "?", []interface{}{v.DoubleValue}, nil

		case *exprpb.Constant_StringValue:
			return "?", []interface{}{v.StringValue}, nil

		case *exprpb.Constant_NullValue:
			return "NULL", nil, nil
		}

		return "", nil, errors.New("unsupported constant type")

	// identifiers: true / false
	case *exprpb.Expr_IdentExpr:
		switch expression_kind.IdentExpr.Name {
		case "true":
			return "1=1", nil, nil
		case "false":
			return "1=0", nil, nil
		default:
			return "", nil, fmt.Errorf("unsupported identifier: %s", expression_kind.IdentExpr.Name)
		}

	// record.field OR request_user.field
	case *exprpb.Expr_SelectExpr:
		// operand must be an identifier
		ident, ok := expression_kind.SelectExpr.Operand.ExprKind.(*exprpb.Expr_IdentExpr)
		if !ok {
			return "", nil, errors.New("unsupported select operand")
		}

		switch ident.IdentExpr.Name {
		case "log", "aircrew":
			table_name, ok := scope[ident.IdentExpr.Name]
			if !ok {
				return "", nil, fmt.Errorf("unknown table alias: %s", ident.IdentExpr.Name)
			}
			return fmt.Sprintf("%s.%s", table_name, expression_kind.SelectExpr.Field), nil, nil

		case "request_user":
			accessor, ok := types.UserClaimsAccessors[expression_kind.SelectExpr.Field]
			// val, ok := request_user[expression_kind.SelectExpr.Field]
			if !ok {
				return "", nil, fmt.Errorf("request_user.%s not provided", expression_kind.SelectExpr.Field)
			}

			val := accessor(&request_user)

			// If it'slice a slice, return as-is for `in`
			if slice, ok := val.([]interface{}); ok {
				return "?", slice, nil
			}

			return "?", []interface{}{val}, nil

		default:
			return "", nil, fmt.Errorf("unsupported select base: %s", ident.IdentExpr.Name)
		}

	// binary operators
	case *exprpb.Expr_CallExpr:
		if len(expression_kind.CallExpr.Args) != 2 {
			return "", nil, errors.New("only binary operators are supported")
		}

		left_sql, left_args, err := compileExpression(expression_kind.CallExpr.Args[0], scope, request_user)
		if err != nil {
			return "", nil, err
		}

		right_sql, right_args, err := compileExpression(expression_kind.CallExpr.Args[1], scope, request_user)
		if err != nil {
			return "", nil, err
		}

		sql_operation := map[string]string{
			"_==_":     "=",
			"_!=_":     "!=",
			"_<_":      "<",
			"_<=_":     "<=",
			"_>_":      ">",
			"_>=_":     ">=",
			"_&&_":     "AND",
			"_||_":     "OR",
			"_in_":     "IN",
			"@in":      "IN",
			"_not_in_": "NOT IN",
			"@not_in":  "NOT IN",
		}[expression_kind.CallExpr.Function]

		if sql_operation == "" {
			return "", nil, fmt.Errorf("unsupported operator: %s", expression_kind.CallExpr.Function)
		}

		// Support `in`/`not in` with the requesting identifier(s) on the left side
		// i.e.: request.unit_ids in record.unit_ids
		if sql_operation == "IN" || sql_operation == "NOT IN" {
			// left_args must be a slice (even if there's only one value)
			if len(left_args) == 0 {
				return "", nil, errors.New("left side of `in`/`not in` must be a list")
			}

			// build (?, ?, ?)
			placeholders := make([]string, len(left_args))
			for i := range placeholders {
				placeholders[i] = "?"
			}

			sql := fmt.Sprintf("(%s %s (%s))", strings.Join(placeholders, ", "), sql_operation, right_sql)
			return sql, left_args, nil
		}

		return fmt.Sprintf("(%s %s %s)", left_sql, sql_operation, right_sql), append(left_args, right_args...), nil

	// Exists operation
	case *exprpb.Expr_ComprehensionExpr:
		comp := expression_kind.ComprehensionExpr

		// Detect `exists`, exists() always uses:
		// accu_init = false
		// loop_step = accu || predicate
		if _, ok := comp.AccuInit.ExprKind.(*exprpb.Expr_ConstExpr); !ok {
			return "", nil, errors.New("unsupported comprehension: non-boolean accumulator")
		}

		iter_ident, ok := comp.IterRange.ExprKind.(*exprpb.Expr_IdentExpr)
		if !ok {
			return "", nil, errors.New("exists() must iterate over a table alias identifier")
		}

		table_name, ok := scope[iter_ident.IdentExpr.Name]
		if !ok {
			return "", nil, fmt.Errorf("unknown table alias in exists(): %s", iter_ident.IdentExpr.Name)
		}

		// Extract predicate from: accu || predicate
		call, ok := comp.LoopStep.ExprKind.(*exprpb.Expr_CallExpr)
		if !ok || call.CallExpr.Function != "_||_" {
			return "", nil, errors.New("unsupported comprehension form (expected OR in loop step)")
		}

		predicate_expr := call.CallExpr.Args[1]

		// Compile predicate
		predicate_sql, predicate_args, err := compileExpression(predicate_expr, scope, request_user)
		if err != nil {
			return "", nil, err
		}

		sql := fmt.Sprintf(
			`EXISTS (
			SELECT 1
			FROM %s
			WHERE %s
		)`,
			table_name,
			predicate_sql,
		)

		return sql, predicate_args, nil
	}

	return "", nil, errors.New("unsupported expression type")
}

func EvaluateRead(txid uuid.UUID, resource string, operation string, scope map[string]string, request_user types.UserClaims, policies []types.PermissionDTO) (string, []interface{}, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(EvaluateRead))

	for _, policy := range policies {
		// Implicit deny overrides any allow
		if policy.Effect != "allow" {
			return "", nil, errors.New("not authorized")
		}
	}

	// If we made it here we are authorized to read, let's build the mysql filters
	var filters []string
	var args []interface{}

	for _, policy := range policies {
		// compile CEL -> SQL filter (we'll implement this next)
		sql, p, err := compileCelToSQL(policy.ConditionExpression, scope, request_user)
		if err != nil {
			log.Printf("%s\n", err.Error())
			return "", nil, err
		}
		filters = append(filters, sql)
		args = append(args, p...)
	}

	// Combine filters using OR
	filter_string := strings.Join(filters, " OR ")
	// This shouldn't happen
	if filter_string == "" {
		// Implicit deny overrides any allow
		filter_string = "1=0"
	}
	return filter_string, args, nil
}

func EvaluateWrite(txid uuid.UUID, resource string, operation string, record map[string]interface{}, request_user map[string]interface{}, policies []types.PermissionDTO) (bool, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(EvaluateRead))

	env, err := cel.NewEnv(
		cel.Variable("record", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("request_user", cel.MapType(cel.StringType, cel.DynType)),
	)
	if err != nil {
		return false, err
	}

	vars := map[string]interface{}{
		"record":       record,
		"request_user": request_user,
	}

	allowed := false
	for _, policy := range policies {
		ast, issues := env.Compile(policy.ConditionExpression)
		if issues != nil && issues.Err() != nil {
			return false, issues.Err()
		}

		program, err := env.Program(ast)
		if err != nil {
			return false, err
		}

		result, _, err := program.Eval(vars)
		if err != nil {
			return false, err
		}

		ok, is_boolean := result.Value().(bool)
		if !is_boolean {
			// TODO [drd] log that this is an invalid policy
			return false, fmt.Errorf("not authorized")
		}
		// Deny any policies that the effect is not `allow`
		if policy.Effect == "allow" && ok {
			allowed = true
		} else {
			return false, errors.New("not authorized")
		}
	}
	return allowed, nil
}
