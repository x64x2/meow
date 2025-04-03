#ifndef MEOW_INTEGRAL_IMAGE_H_
#define MEOW_INTEGRAL_IMAGE_H_

#include <stdint.h>

// See https://en.wikipedia.org/wiki/Summed-area_table.

// This represents an integral image in memory.
typedef struct {
  // Store the width and height of the image in pixels.
  unsigned int width, height;
  
  // Store the summed pixel values.
  unsigned int *values;
} meow_integral_image_t;

// This function generates integral images from regular images.
void meow_integrate_image(meow_integral_image_t *, unsigned char *, uint32_t, uint32_t);

// This function returns the sum of all the values in a given rectuangular area
// in an image. The ranges are exclusive.
unsigned int meow_sum_area(
  meow_integral_image_t,
  unsigned int,
  unsigned int,
  unsigned int,
  unsigned int
);

#endif
