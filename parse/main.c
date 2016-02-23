// Toy implementation of infix to postfix expression parser
//
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdbool.h>

typedef struct opstack{
  int count;
  char** stack;
} opstack_t;

opstack_t* stack_new();
int stack_push(opstack_t *s, char* op);
char* stack_pop(opstack_t *s);
char* stack_peek(opstack_t *s);

void print(int outc, char** outv);

int shunt(opstack_t* s, int argc, char** argv);

bool is_num(char*);

bool is_func(char*);

bool is_op(char*);
int prec_op(char*);
int assoc_op(char*);

bool is_lparen(char*);
bool is_rparen(char*);

int
main(int argc, char** argv){
  struct opstack s  = { argc, argv};

  printf("In:\n");
  print(argc - 1, &argv[1]);

  opstack_t *ops = stack_new();
  stack_push(ops, "hello");
  stack_push(ops, "world");

  printf("Before:\n");
  print(ops->count, ops->stack);
  printf("peek: %s\n", stack_peek(ops));
  char* popped = stack_pop(ops);
  printf("poppedk: %s\n", stack_peek(ops));
  printf("After:\n");
  print(ops->count, ops->stack);
  return 0;
}

opstack_t*
stack_new(){
  opstack_t* s = (opstack_t*) malloc(sizeof(opstack_t));
  if (s != NULL){
    s   = memset(s,0,sizeof(opstack_t));
  };
  return s;
}

int
stack_push(opstack_t *s, char* op){
    s->stack = (char**)realloc((void*)(s->stack), sizeof(char*) * s->count);
    if (s->stack == NULL){
      return -1;
    }
    s->count += 1;
    s->stack[s->count -1] = op;
    return 0;
}

char* stack_pop(opstack_t *s){
  char *t;
  if (s->count == 0){
    return NULL;
  }

  t = s->stack[s->count -1];
  s->count -= 1;
  s->stack = (char**)realloc((void*)(s->stack), sizeof(char*) * s->count);

  return t;
}

char* stack_peek(opstack_t *s){
  return s->stack[s->count -1];
}

int
shunt(opstack_t* out, int argc, char** argv){
  for(int i = 0; i < argc; i++){
    char* t = argv[i];
    if (is_num(t)){
    } else if (is_num(t)){
    } else {
      return -1;
    }
  }
}

void
print(int outc, char** outv){
  for(int i = 0; i < outc; i++){
    printf("%s ", outv[i]);
  };
  printf("\n");
}


