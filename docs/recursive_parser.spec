a recursive descent parser is a kind of top-down parser 
built from a set of mutually recursive procedures (or a non-recursive equivalent)
 where each such procedure implements one of the nonterminals of the grammar.

i.e for each token we have to write recursive function.


logic :

stmt(){


    switch (nexttoken){

        case expr:
             match(expr);
             match(end);
             return;

        case if:
            match(if)
            match('{')
            match(expr i.e if Condition)
            match('}')
            
            //recursive call to stmt()
            stmt( i.e code under if block)
            return;

    }
}


match( token string){
    if(nexttoken ==token ){
        nexttoken = nextterminal
    }
}







func stmt(){

    switch(nexttoken){

        a) var statements
        b) return statment
       c)  expression
            i) assignment_expression
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
               init_declarator_list
               init_declarator
               expression
               block_item
               compound_statement
               block_item_list
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
               infer_clauses
               indexing_literal
               indexing_slice_literal
               IDENTIFIER 
               external_declaration
               LBRACE
               package_declaration  
               global_declaration   
               function_declaration
               import_declaration
               struct_declaration   
    }
}


// new actions convert yacc code to golang code to handle token.

we have add  approx new  120 actions older ast actions can be used to parse token.


for example :

this yacc code will form new ast actions

function_parameters:
                LPAREN RPAREN
                { $$ = nil }
        |       LPAREN parameter_type_list RPAREN
                { $$ = $2 }
                ;

like ast.function_parameters()  


statement - actions.statement()  // new actions 

labeled_statement - actions.labeled_statement()  // // new actions

global_declaration - actions.DeclareGlobal()

function_header - actions.FunctionHeader()


declare_struct -  actions.DeclareStruct() 

package_declaration  - actions.DeclarePackage() 

function_parameters - actions.function_parameters() 

parameter_type_list  - actions.parameter_type_list()

fields - actions.fields()

import_declaration - actions.import_declaration()

function_header - actions.functions_header()

function_parameters -actions.function_parameters()



function_declaration - actions.FunctionDeclaration()

parameter_type_list  - actions.parameter_type_list()

parameter_list - actions.parameter_list()


parameter_declaration - actions.parameter_declaration()

declarator - actions.direct_declarator()


id_list -actions.indentifier()

types_list - actions.type_list()


struct_literal_fields - actions.StructLiteralAssignment()

array_literal_expression_list - actions.array_literal_expression_list()

indexing_literal - actions.indexing_literal()


indexing_slice_literal - actions.indexing_slice_literal()


array_literal_expression - actions.indexing_slice_literal()


slice_literal_expression_list - actions.slice_literal_expression_list()

infer_action_arg - actions.infer_action_arg()

infer_clauses- actions.infer_clauses()

primary_expression - actions.primary_expression()

after_period - actions.after_period()

postfix_expression - actions.postfix_expression()


argument_expression_list - actions.argument_expression_list()


unary_expression - actions.unary_expression()


multiplicative_expression - actions.multiplicative_expression()


additive_expression - actions.additive_expression()

shift_expression - action.shift_expression()



relational_expression - actions.relational_expression()


assignment_operator - actions.assignment_operator()


expression - actions.expression()

declaration - actions.declaration()




tokens

BOOLEAN_LITERAL

BYTE_LITERAL

SHORT_LITERAL

INT_LITERAL

LONG_LITERAL

UNSIGNED_BYTE_LITERAL

UNSIGNED_SHORT_LITERAL

UNSIGNED_LONG_LITERAL

FLOAT_LITERAL

DOUBLE_LITERAL

FUNC OP LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK IDENTIFIER

VAR COMMA PERIOD COMMENT STRING_LITERAL PACKAGE IF ELSE FOR TYPSTRUCT STRUCT

SEMICOLON NEWLINE

ASSIGN CASSIGN IMPORT RETURN GOTO GT_OP LT_OP GTEQ_OP LTEQ_OP EQUAL COLON NEW

EQUALWORD GTHANWORD LTHANWORD

GTHANEQ LTHANEQ UNEQUAL AND OR

ADD_OP SUB_OP MUL_OP DIV_OP MOD_OP REF_OP NEG_OP AFFVAR

PLUSPLUS MINUSMINUS REMAINDER LEFTSHIFT RIGHTSHIFT EXP
NOT
                        
                        
BITXOR_OP BITOR_OP BITCLEAR_OP

PLUSEQ MINUSEQ MULTEQ DIVEQ REMAINDEREQ EXPEQ
                        
LEFTSHIFTEQ RIGHTSHIFTEQ BITANDEQ BITXOREQ BITOREQ

DEC_OP INC_OP PTR_OP LEFT_OP RIGHT_OP
GE_OP LE_OP EQ_OP NE_OP AND_OP OR_OP
ADD_ASSIGN AND_ASSIGN LEFT_ASSIGN MOD_ASSIGN

MUL_ASSIGN DIV_ASSIGN OR_ASSIGN RIGHT_ASSIGN
SUB_ASSIGN XOR_ASSIGN
BOOL F32 F64
I8 I16 I32 I64

STR
UI8 UI16 UI32 UI64
UNION ENUM CONST CASE DEFAULT SWITCH BREAK CONTINUE
TYPE
                        
/* Types */
BASICTYPE

/* Removers */
REM DEF EXPR FIELD CLAUSES OBJECT OBJECTS

/* Debugging */
DSTACK DPROGRAM DSTATE

/* Affordances */

AFF CAFF TAG INFER VALUE

/* Pointers */
ADDR

int_value

after_period
unary_operator
assignment_operator



**Actions called**
 - `actions.DeclarationSpecifier`
 - `actions.StructLiteralAssignment`
 - `actions.Assignment`
 - `actions.SliceLiteralExpression`
 - ...


Actions API

Functionalities:

Assignments
 - Assignment of struct literals
 - Assignment of array literals
 - short assignment(:=, +=, *=, etc)

Declaration
 - DeclareGlobal
 - DeclareGlobalInPackage
 - DeclareStruct
 - DeclarePackage
 - DeclareImport
 - DeclareLocal
 - DeclarationSpecifiers
 - DeclarationSpecifiersBasic
 - DeclarationSpecifiersStruct

Expressions
 -IterationExpressions
 -trueJmpExpressions
 -BreakExpressions
 -ContinueExpressions
 -SelectionExpressions
 -OperatorExpression
 -UnaryExpression
 -AssociateReturnExpressions
 -AddJmpToReturnExpressions

Functions
 -FunctionHeader (takes a function name ('ident') and either creates the function if it's not known before or returns the already existing function)
 -FunctionAddParameters
 -FunctionProcessParameters
 -FunctionDeclaration
 -FunctionCall
 -ProcessOperatorExpression
 -ProcessPointerStructs
 -ProcessExpressionArguments
 -AddPointer
 -CheckRedeclared
 -ProcessLocalDeclaration
 -ProcessGoTos
 -CheckTypes
 -ProcessStringAssignment
 -ProcessReferenceAssignment
 -ProcessSlice
 -ProcessSliceAssignment
 -UpdateSymbolsTable
 -ProcessMethodCall
 -GiveOffset
 -ProcessTempVariable
 -CopyArgFields
 -ProcessSymbolFields
 -SetFinalSize
 -GetGlobalSymbol
 -PreFinalSize

Literals
 -SliceLiteralExpression
 -PrimaryStructLiteral
 -PrimaryStructLiteralExternal
 -ArrayLiteralExpression

PostFix
 -PostfixExpressionArray
 -PostfixExpressionNative
 -PostfixExpressionEmptyFunCall
 -PostfixExpressionFunCall
 -PostFixExpressionIncDec
 -PostfixExpressionField

Scope
 -DefineNewScope

Statements
 -SelectionStatement


AddEmptyFunctionToPackage -

 https://github.com/skycoin/cx/blob/develop/cx/astapi/functions.go#L27


RemoveFunctionFromPackage - https://github.com/skycoin/cx/blob/develop/cx/astapi/functions.go#L59


RemoveFunctionFromPackage - https://github.com/skycoin/cx/blob/develop/cx/astapi/functions.go#L93


RemoveFunctionInput - https://github.com/skycoin/cx/blob/develop/cx/astapi/functions.go#L132

AddNativeOutputToFunction - https://github.com/skycoin/cx/blob/develop/cx/astapi/functions.go#L166

RemoveFunctionOutput - https://github.com/skycoin/cx/blob/develop/cx/astapi/functions.go#L166

GetPackagesNameList - https://github.com/skycoin/cx/blob/develop/cx/astapi/packages.go#L22


AddEmptyPackage -https://github.com/skycoin/cx/blob/develop/cx/astapi/packages.go#L42


AddNativeInputToExpression -
https://github.com/skycoin/cx/blob/develop/cx/astapi/arguments.go#L35

RemoveInputFromExpression -https://github.com/skycoin/cx/blob/develop/cx/astapi/arguments.go#L81


AddNativeExpressionToFunction - https://github.com/skycoin/cx/blob/develop/cx/astapi/arguments.go#L81

AddNativeExpressionToFunction -https://github.com/skycoin/cx/blob/develop/cx/astapi/expressions.go#L100

GetPackagesNameList -https://github.com/skycoin/cx/blob/develop/cx/astapi/expressions.go#L100


AddEmptyPackage - https://github.com/skycoin/cx/blob/develop/cx/astapi/packages.go#L42
