project(
	'meow',
	'c',
	version: '0.1.0',
    meson_version: '>= 0.59.0',
    default_options: [
    	'warning_level=2',
    	'werror=false',
    	'c_std=c99',
    ],
)

add_project_arguments(
	[ '-I' + meson.project_build_root(), ],
	language: 'c',
)

nasm = find_program('nasm')
nasm_gen = generator(
	nasm,
 	output  : '@BASENAME@.o',
	arguments : [
		'-f', meson.get_external_property('nasm_format'),
		'@INPUT@',
		'-o', '@OUTPUT@'
	]
)

spng_proj = subproject('spng', default_options:['enable_opt=false'])
spng_dep = spng_proj.get_variable('spng_dep')

subdir('meow')
subdir('src')
