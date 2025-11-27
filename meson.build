project(
  'Player',
  'cpp',
  version : '0.1',
  meson_version : '>= 1.3.0',
  default_options : ['warning_level=3', 'cpp_std=c++17'],
)

qt6=import('qt6')
qt_modules=[
    'Core',
    'Widgets',
    'Gui',
    'Multimedia'
]
qt6_deps=dependency('qt6',modules: qt_modules)

sources=[
    'player.cpp',
    'mainwindow.cpp'
]
headers=[
    'mainwindow.h'
]

moc_sources=qt6.preprocess(
    moc_headers: headers
)

exe = executable(
  'Player',
  sources + moc_sources,
  install : true,
  dependencies : qt6_deps,
)
