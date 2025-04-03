#include "meow/haar_features.h"

// TODO This can be optimised much further (sum_area's parameters can be pre-
// calculated during training).

// TODO Actually think about this code (i.e. maximum values and such).

int meow_haar_x2(
	meow_integral_image_t img,
	unsigned int x0,
	unsigned int y0,
	unsigned int width,
	unsigned int height
) {
	unsigned int half_width = width >> 1;
	unsigned int left_value = meow_sum_area(img, x0, y0, x0 + half_width, y0 + height);
	unsigned int right_value = meow_sum_area(img, x0 + half_width, y0, x0 + width, y0 + height);
  return left_value - right_value;
}

int meow_haar_y2(
	meow_integral_image_t img,
	unsigned int x0,
	unsigned int y0,
	unsigned int width,
	unsigned int height
) {
	unsigned int half_height = height >> 1;
	unsigned int top_value = meow_sum_area(img, x0, y0, x0 + height, y0 + half_height);
	unsigned int bottom_value = meow_sum_area(img, x0, y0 + half_height, x0 + width, y0 + height);
	return top_value - bottom_value;
}

int meow_haar_x3(
  meow_integral_image_t img,
  unsigned int x0,
  unsigned int y0,
  unsigned int width,
  unsigned int height
) {
  unsigned int third_width = width / 3;
  unsigned int left_value = meow_sum_area(img, x0, y0, x0 + third_width, y0 + height);
  unsigned int middle_value = meow_sum_area(img, x0 + third_width, y0, x0 + third_width * 2, y0 + height);
  unsigned int right_value = meow_sum_area(img, x0 + third_width * 2, y0, x0 + width, y0 + height);
  return left_value + right_value - middle_value;
}

int meow_haar_y3(
  meow_integral_image_t img,
  unsigned int x0,
  unsigned int y0,
  unsigned int width,
  unsigned int height
) {
  unsigned int third_height = height / 3;
  unsigned int top_value = meow_sum_area(img, x0, y0, x0 + width, y0 + third_height);
  unsigned int middle_value = meow_sum_area(img, x0, y0 + third_height, x0 + width, y0 + third_height * 2);
  unsigned int bottom_value = meow_sum_area(img, x0, y0 + third_height * 2, x0 + width, y0 + height);
  return top_value + bottom_value - middle_value;
}

int meow_haar_x2y2(
  meow_integral_image_t img,
  unsigned int x0,
  unsigned int y0,
  unsigned int width,
  unsigned int height
) {
  unsigned int half_width = width >> 1;
  unsigned int half_height = height >> 1;
  unsigned int top_left_value = meow_sum_area(img, x0, y0, x0 + half_width, y0 + half_height);
  unsigned int top_right_value = meow_sum_area(img, x0 + half_width + 1, y0, x0 + width, y0 + half_height);
  unsigned int bottom_left_value = meow_sum_area(img, x0, y0 + half_height + 1, x0 + half_width, y0 + height);
  unsigned int bottom_right_value = meow_sum_area(img, x0 + half_width + 1, y0 + half_height + 1, x0 + width, y0 + height);
  return top_left_value + bottom_right_value - top_right_value - bottom_left_value;
}
