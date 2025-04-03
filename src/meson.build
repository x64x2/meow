meow_src = [
	'haar_features.c',
	#'integral_image.c',
	'weak_classifier.c',
	'strong_classifier.c',
]

meow_src += nasm_gen.process([
	'integral_image.asm',
	#'haar_features.asm',
	'model.asm',
])

meow_lib = static_library(
	'meow',
	meow_src,
	include_directories: meow_inc,
)
