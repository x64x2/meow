#include "meow/strong_classifier.h"

#include "meow/weak_classifier.h"

bool meow_strong_classify(
  meow_strong_classifier_t classifier,
  meow_integral_image_t image,
  unsigned int scale,
  unsigned char scale_shift // TODO: Make scale_shift constanttttt
) {
  int sum = 0;

  for (size_t i = 0; i < NO_WEAK_CLASSIFIERS; i++) {
    // TODO Use for loop
    if (meow_weak_classify(classifier.classifiers[i], image, scale, scale_shift)) {
      sum += classifier.classifiers[i].alpha;
    }
  }

  return sum > 0;
}
