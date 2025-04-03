meow_inc = include_directories('include')
subdir('src')

meow_dep = declare_dependency(
	link_with: meow_lib,
	include_directories: meow_inc,
)

subdir('tests')