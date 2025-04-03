#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "meow/haar_features.h"

#define IN_VALUES { \
    0,    0,    0,    0,    0,    0, \
    0,    0,    0,  255,  510,  765, \
    0,    0,    0,  510, 1020,  1530, \
    0,    0,    0,  765,  1530, 2295 \
}
#define IN_WIDTH 6
#define IN_HEIGHT 4
#define IN_SIZE IN_WIDTH * IN_HEIGHT * sizeof(unsigned int)

#define EXPECTED_VALUE -510

int main() {
  meow_integral_image_t integral_img;
  integral_img.width = IN_WIDTH;
  integral_img.height = IN_HEIGHT;
  const unsigned int in_values[] = IN_VALUES;
  integral_img.values = in_values;

  if (meow_haar_y2(integral_img, 1, 0, 3, 3) != EXPECTED_VALUE) {
    return EXIT_FAILURE;
  }
  
  return EXIT_SUCCESS;
}
