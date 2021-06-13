/*

Literal :

BOOLEAN_LITERAL

BYTE_LITERAL

SHORT_LITERAL

INT_LITERAL

LONG_LITERAL

UNSIGNED_BYTE_LITERAL

UNSIGNED_SHORT_LITERAL

UNSIGNED_INT_LITERAL

UNSIGNED_LONG_LITERAL

FLOAT_LITERAL

DOUBLE_LITERAL


tok :


FUNC

OP




Expression :


assignment_expression


constant_expression


conditional_expression


logical_or_expression


logical_and_expression


exclusive_or_expression


inclusive_or_expression


and_expression


relational_expression

shift_expression


additive_expression


multiplicative_expression



unary_expression

argument_expression_list

postfix_expression


primary_expression


struct_literal_expression

array_literal_expression_list


array_literal_expression


slice_literal_expression_list

slice_literal_expression

return_expression

selector

struct_literal_fields

elseif

elseif_list

declaration


init_declarator_list


initializer


expression

block_item


block_item_list


compound_statement

else_statement


labeled_statement

expression_statement


selection_statement


iteration_statement


jump_statement

statement


function_header

infer_action_arg


infer_action

infer_actions


infer_clauses


indexing_literal

indexing_slice_literal

IDENTIFIER

LBRACE

translation_unit





actions :



// DeclareGlobal creates a global variable in the current package.

global_declaration : actions.DeclareGlobal($2, $3, nil, false)


struct_declaration : actions.DeclareStruct($2, $4)


struct_literal_fields :

actions.StructLiteralAssignment([]*ast.CXExpression{actions.StructLiteralFields($1)}, $3)

package_declaration : actions.DeclarePackage($2)


function_header  : actions.FunctionHeader($2, nil, false)


function_parameters :

function_declaration  : actions.FunctionDeclaration($1, $2, nil, $3)



direct_declarator :


array_literal_expression :
	actions.ArrayLiteralExpression($1, $2, $4)


primary_expression:

 actions.PrimaryIdentifier($1)


 .
 .
 .


postfix_expression:
	actions.PostfixExpressionArray($1, $3)



multiplicative_expression:
	actions.OperatorExpression($1, $3, constants.OP_MUL)

*/

/*

Execution of AST :


RunCompiled  : starting point

RunCxAst

RunCompiled

ToCall :


callback()



*/

package ast
