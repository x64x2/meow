#include "meow/integral_image.h"
#include <stdlib.h>
#include <string.h>

unsigned char meow_integrate_image(
  meow_integral_image_t *dst,
  unsigned char *values,
  uint32_t width,
  uint32_t height
) {
  dst->width = width + 1;
  dst->height = height + 1;
    
  dst->values = calloc(dst->width * dst->height, sizeof(*dst->values));
  
  // Integrate the image
  for (unsigned int y = 0; y < height; y++) {
    for (unsigned int x = 0; x < width; x++) {
      dst->values[dst->width * (y + 1) + x + 1] =
        values[3 * (width * y + x)]
        - dst->values[dst->width * y + x]
        + dst->values[dst->width * y + x + 1]
        + dst->values[dst->width * (y + 1) + x];
    }
  }
  
  return EXIT_SUCCESS;
}

unsigned int meow_sum_area(
  meow_integral_image_t img,
  unsigned int x0,
  unsigned int y0,
  unsigned int x1,
  unsigned int y1
) {
  return
    img.values[img.width * y1 + x1]
    - img.values[img.width * y0 + x1]
    - img.values[img.width * y1 + x0]
    + img.values[img.width * y0 + x0];
}
