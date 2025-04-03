#include <stdlib.h>
#include <string.h>

#include "meow/integral_image.h"

#define IN_VALUES { \
  255, 255, 255, \
  0,   0,   0, \
  255, 255, 255, \
  0,   0,   0, \
  255, 255, 255, \
  0,   0,   0, \
  255, 255, 255, \
  0,   0,   0, \
  255, 255, 255, \
  0,   0,   0, \
  255, 255, 255, \
  0,   0,   0, \
  255, 255, 255, \
  0,   0,   0, \
  255, 255, 255, \
}
#define IN_WIDTH 5
#define IN_HEIGHT 3

#define EXPECTED_VALUES { \
  0,    0,    0,    0,    0,    0, \
  0,  255,  255,  510,  510,  765, \
  0,  255,  510,  765,  1020, 1275, \
  0,  510,  765, 1275, 1530, 2040 \
}

int main() {
  unsigned char in_values[] = IN_VALUES;

  meow_integral_image_t integral_img;
  integral_img.width = 69;
  meow_integrate_image(&integral_img, in_values, IN_WIDTH, IN_HEIGHT);
  
  const unsigned int expected_values[] = EXPECTED_VALUES;

  // TODO Use memcmp.
  for (int i = 0; i < 24; i++) {
    if (integral_img.values[i] != expected_values[i]) {
      return EXIT_FAILURE;
    }
  }

  free(integral_img.values);

  return EXIT_SUCCESS;
}

