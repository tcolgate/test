#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>
#include <math.h>
#include "lookup3.h"

#include <stdio.h>

typedef uint8_t* bitset;
bitset bitset_new(uint32_t size){
  return (bitset) calloc(1 + size/8,sizeof(char));
};

void bitset_set(bitset bs, uint32_t n){
  uint32_t B = n / 8;
  uint32_t b = n % 8;
  bs[B] = bs[B] | ((uint8_t)1 << b);
};

bool bitset_isset(bitset bs, uint32_t n){
  uint32_t B = n / 8;
  uint32_t b = n % 8;
  return ((bs[B] & (1 << b)) == ((uint8_t)1<<b));
};

void bitset_clear(bitset b){
};

typedef struct bloom {
  uint64_t count;
  bitset bits;
  uint32_t m;
  uint8_t k;
}* bloom_t;

bloom_t bloom_new(uint64_t n, double fprate){
  double fillRatio = 0.5;

  bloom_t b = (bloom_t) malloc(sizeof(bloom_t));

  b->m = (uint32_t) ceil((double)n / (log10(fillRatio) *
                  log10(1-fillRatio)) / fabs(log(fprate)));
  b->k = (uint8_t) ceil(log2(1/fprate));

  b->bits = bitset_new(b->m);

  return b;
}

bool bloom_test(bloom_t b, void* key, uint32_t len){
  uint32_t h=0,l=0;
  hashlittle2(key,len,&h,&l);
  
  // Check the K bits.
  for (uint32_t i = 0; i < b->k; i++ ){
    if (!bitset_isset(b->bits,(l+h*i) % b->m))
      return false;
  }

  return true;
};

void bloom_add(bloom_t b, void* key, uint32_t len){
  uint32_t h=0,l=0;
  hashlittle2(key,len,&h,&l);

  // Set the K bits.
  for (uint32_t i = 0; i < b->k; i++ ){
    bitset_set(b->bits,(l+h*i) % b->m);
  }
  b->count++;
};

