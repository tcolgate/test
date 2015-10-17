#include <string.h>
#include "../bloom.h"

int main(int argc, char** argv){
  char *thing = "thing";
  char *othing = "thing2";

  bloom_t b = bloom_new(1000,0.001);

  if (b->m != 1598)
    return 1;
  if (b->k != 10)
    return 1;

  if (bloom_test(b,thing,strlen(thing)+1))
    return 1;

  bloom_add(b,thing,strlen(thing)+1);

  if (b->count != 1)
    return 1;

  if (!bloom_test(b,thing,strlen(thing)+1))
    return 1;
  if (bloom_test(b,othing,strlen(othing)+1))
    return 1;

  return 0;
}

