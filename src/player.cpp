#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include "core/playercontroller.h"

int main(int argc, char **argv) {
  QGuiApplication app(argc, argv);

  PlayerController controller;

  QQmlApplicationEngine engine;
  engine.rootContext()->setContextProperty("player", &controller);
  engine.load(QUrl("qrc:/qml/Main.qml"));

  return app.exec();
}
