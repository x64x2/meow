#include <stdbool.h>
#include "haar_features.h"
#include "integral_image.h"

#ifndef meow_WEAK_CLASSIFIER_H_
#define meow_WEAK_CLASSIFIER_H_

typedef struct {
  short alpha;

  int threshold;

  char parity;

  enum meow_feature feature_type;

  unsigned int x0;
  unsigned int y0;
  unsigned int width;
  unsigned int height;
} meow_weak_classifier_t;

bool meow_weak_classify(
  meow_weak_classifier_t,
  meow_integral_image_t,
  unsigned int,
  unsigned int,
  unsigned int,
  unsigned char
);

#endif
