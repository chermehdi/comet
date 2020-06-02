// Parser rules
grammar Comet;

comet:
    function_declaration+ EOF
    ;

function_declaration:
    'func' IDENTIFIER'('parameterless')' '{' statements '}' # functionDeclaration
    ;

parameterless:
    parameters*
    ;

parameters:
    IDENTIFIER(','IDENTIFIER)* | expression (','expression)*
    ;

statements:
    statement*
    ;

statement:
    assignment_statement                      # assignment
    | if_statement                            # ifStatement
    | return_statement                        # returnStatement
    | for_loop                                # forLoop
    ;

assignment_statement:
    'var' IDENTIFIER '=' expression
    ;

if_statement:
    'if' '('expression')' '{' statements '}' ('else' '{' statements '}')?
    ;

return_statement:
    'return' expression
    ;

for_loop:
    'for' '('expression')' '{' statements '}'
    ;

// struct and object declaration to be defined later

expression:
    function_call
    | IDENTIFIER
    | STRING 
    | NUMBER 
    | BOOL 
    | string_concat
    | math_expression
    | array_declaration
    ;

array_declaration:
    '['(expression (','expression)*)?']'
     ;

function_call:
    IDENTIFIER '('parameterless')'
    ;

string_concat:
    STRING ('+' STRING)*
    ;

math_expression:
    term (('+' | '-') term)*
    ;

term:
    factor (('*' | '/') factor)*
    ;

factor:
    IDENTIFIER
    | NUMBER
    | '('expression')'
    ;

// Lexer rules

WS : [ \t\r\n\u000C]+ -> skip;
COMMENT : '/*' .*? '*/' -> skip;
LINE_COMMENT : '//' ~[\r\n]* -> skip;
IDENTIFIER: ('_'|LETTER)(LETTER|DIGIT|'_')*;
NUMBER: ('-' | '+')? (([1-9][0-9]*) | ([0-9]));
OPERATOR: '+' | '-' | '*' | '/' | '%' | '!';
fragment DIGIT: [0-9];
BOOL: 'true' | 'false';
LETTER: [a-zA-Z];
STRING: '"' ~('\r' | '\n' | '"')* '"' ;
