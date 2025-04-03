#include <stdbool.h>
#include "integral_image.h"
#include "weak_classifier.h"

#ifndef MEOW_STRONG_CLASSIFIER_H_
#define MEOW_STRONG_CLASSIFIER_H_

#define NO_WEAK_CLASSIFIERS 128

typedef struct {
  meow_weak_classifier_t classifiers[NO_WEAK_CLASSIFIERS];
} meow_strong_classifier_t;

bool meow_strong_classify(
  meow_strong_classifier_t,
  meow_integral_image_t,
  unsigned int,
  unsigned char
);

#endif
