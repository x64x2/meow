#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "meow/haar_features.h"

#define IN_VALUES { \
    0,    0,    0,    0,    0, \
    0,    0,    0,  255,  510, \
    0,    0,    0,  510, 1020, \
    0,  255,  510, 1020, 1530, \
    0,  510, 1020, 1530, 2040 \
}
#define IN_WIDTH 5
#define IN_HEIGHT 5
#define IN_SIZE IN_WIDTH * IN_HEIGHT * sizeof(unsigned int)

#define EXPECTED_VALUE -1020

int main() {
  meow_integral_image_t integral_img;
  integral_img.width = IN_WIDTH;
  integral_img.height = IN_HEIGHT;
  const unsigned int in_values[] = IN_VALUES;
  integral_img.values = in_values;

  if (meow_haar_x2y2(integral_img, 0, 0, 4, 4) != EXPECTED_VALUE) {
    return EXIT_FAILURE;
  }
  
  return EXIT_SUCCESS;
}
