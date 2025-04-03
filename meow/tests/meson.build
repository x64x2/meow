images_dir = join_paths(meson.current_source_dir(), 'images')

test(
	'Integral image generation',
	executable(
		'integral_image',
		'integral_image.c',
		dependencies: meow_dep
	),
)
test(
	'Area summing',
	executable(
		'sum_area',
		'sum_area.c',
		dependencies: meow_dep
	),
)
test(
	'X2 Haar feature',
	executable(
		'haar_x2',
		'haar_x2.c',
		dependencies: meow_dep
	),
)
test(
	'Y2 Haar feature',
	executable(
		'haar_y2',
		'haar_y2.c',
		dependencies: meow_dep
	),
)
test(
	'X3 Haar feature',
	executable(
		'haar_x3',
		'haar_x3.c',
		dependencies: meow_dep
	),
)
test(
	'Y3 Haar feature',
	executable(
		'haar_y3',
		'haar_y3.c',
		dependencies: meow_dep
	),
)
test(
	'X2Y2 Haar feature',
	executable(
		'haar_x2y2',
		'haar_x2y2.c',
		dependencies: meow_dep
	),
	should_fail: true,
)
