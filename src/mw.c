#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <spng.h>

#include "meow/model.h"

#define HELP_MESSAGE \
  "Usage: %s [FILE]\n" \
  "Detect cat faces.\n" \
  "\n" \
  "With no FILE, or when FILE is -, read standard input.\n" \
  "\n" \
  "  -h, --help     display this help and exit\n"

int main(int argc, char *argv[]) {
  // check validity of the cmd arg
  if (argc != 2) {
    fputs("ERROR: invalid usage\n", stderr);
    return EXIT_FAILURE;
  }

  // print an help msg
  if (!strcmp("-h", argv[1]) || !strcmp("--help", argv[1])) {
    printf(HELP_MESSAGE, argv[0]);
    return EXIT_SUCCESS;
  } 
  
  // read
  FILE *image_file;
  if (!strcmp("-", argv[1])) {
    image_file = stdin;
  } else {
    image_file = fopen(argv[1], "rb");
  }
  
  if (image_file == NULL) {
    perror("ERROR: cannot read FILE");
    return EXIT_FAILURE;
  }
  
  // decode image.
  int status;
  spng_ctx *spng_handle = spng_ctx_new(0);
  
  spng_set_png_file(spng_handle, image_file);
  
  struct spng_ihdr image_ihdr;
  if ((status = spng_get_ihdr(spng_handle, &image_ihdr))) {
    fprintf(stderr, "ERROR: cannot decode FILE: %s\n", spng_strerror(status));
    spng_ctx_free(spng_handle);
    fclose(image_file);
    return EXIT_FAILURE;
  }
  
  size_t image_buffer_size;
  spng_decoded_image_size(spng_handle, SPNG_FMT_RGB8, &image_buffer_size);
  unsigned char *image_buffer = malloc(image_buffer_size);
  
  status = spng_decode_image(spng_handle, image_buffer, image_buffer_size, SPNG_FMT_RGB8, 0);
  spng_ctx_free(spng_handle);
  fclose(image_file);
  if (status) {
    fprintf(stderr, "ERROR: cannot decode FILE: %s\n", spng_strerror(status));
    return EXIT_FAILURE;
  }
   
   
  return 0;
}
