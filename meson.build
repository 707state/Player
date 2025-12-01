project(
  'Player',
  'cpp',
  version : '0.1',
  meson_version : '>= 1.3.0',
  default_options : ['warning_level=3', 'cpp_std=c++20'],
)

qt6=import('qt6')
qt_modules=[
    'Core',
    'Widgets',
    'Gui',
    'Multimedia',
    'Qml'
]
qt6_deps=dependency('qt6',modules: qt_modules)
qrc=qt6.compile_resources(sources: 'src/app.qrc')
taglib_deps=dependency('taglib')
sources=[
    'src/player.cpp',
    'core/playercontroller.cpp',
    qrc
]
headers=[
    'core/playercontroller.h'
]
incdir=include_directories('.')
moc_headers=qt6.preprocess(
    moc_headers: headers
)

exe = executable(
  'Player',
  sources+moc_headers,
  install : true,
  include_directories:incdir,
  dependencies : [qt6_deps,taglib_deps],
)
