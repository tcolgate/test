#include <stdlib.h>
#include <stdbool.h>
#include <math.h>
#include "lookup3.h"

#include <stdio.h>

typedef char* bitset;
bitset bitset_new(int size){
  return (bitset) calloc(1 + size/8,sizeof(char));
};

void bitset_set(bitset bs, unsigned int n){
  int B = n / 8;
  int b = n % 8;
  bs[B] = bs[B] | (1 << b);
};

bool bitset_isset(bitset bs, unsigned int n){
  int B = n / 8;
  int b = n % 8;
  return ((bs[B] & (1 << b)) == (1<<b));
};

void bitset_clear(bitset b){
};

typedef struct bloom {
  uint64_t count;
  bitset bits;
  unsigned int m;
  unsigned int k;
}* bloom_t;

bloom_t bloom_new(int n, double fprate){
  float fillRatio = 0.5;

  bloom_t b = (bloom_t) malloc(sizeof(bloom_t));

  b->m = (unsigned int) ceil((double)n / (log10(fillRatio) *
                  log10(1-fillRatio)) / fabs(log(fprate)));
  b->k = (unsigned int) ceil(log2(1/fprate));

  b->bits = bitset_new(b->m);

  return b;
}

bool bloom_test(bloom_t b, void* key, int len){
  uint32_t h=0,l=0;
  hashlittle2(key,len,&h,&l);
  
  // Check the K bits.
  for (uint32_t i = 0; i < b->k; i++ ){
    if (!bitset_isset(b->bits,(l+h*i) % b->m))
      return false;
  }

  return true;
};

void bloom_add(bloom_t b, void* key, int len){
  uint32_t h=0,l=0;
  hashlittle2(key,len,&h,&l);

  // Set the K bits.
  for (uint32_t i = 0; i < b->k; i++ ){
    bitset_set(b->bits,(l+h*i) % b->m);
  }
  b->count++;
};

