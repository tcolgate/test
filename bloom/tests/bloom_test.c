#include <string.h>
#include "../bloom.h"

int main(int argc, char** argv){
  char *thing = "thing";
  char *othing = "thing2";

  bloom_t b = bloom_new(1000,0.001);
  printf ("m: %d\n",b->m);
  printf ("k: %d\n",b->k);

  bloom_add(b,othing,strlen(othing)+1);

  return bloom_test(b,thing,strlen(thing)+1) ? 0 : 1;
}

