#ifndef MEOW_HAAR_FEATURES_H_
#define MEOW_HAAR_FEATURES_H_

#include "integral_image.h"

enum meow_feature_type {
  meow_feature_x2,
  meow_feature_y2,
  meow_feature_x3,
  meow_feature_y3,
  meow_feature_x2y2
};

int meow_haar_x2(
  meow_integral_image_t,
  unsigned int,
  unsigned int,
  unsigned int,
  unsigned int
);

int meow_haar_y2(
  meow_integral_image_t,
  unsigned int,
  unsigned int,
  unsigned int,
  unsigned int
);

int meow_haar_x3(
  meow_integral_image_t,
  unsigned int,
  unsigned int,
  unsigned int,
  unsigned int
);

int meow_haar_y3(
  meow_integral_image_t,
  unsigned int,
  unsigned int,
  unsigned int,
  unsigned int
);

int meow_haar_x2y2(
  meow_integral_image_t,
  unsigned int,
  unsigned int,
  unsigned int,
  unsigned int
);

#endif
