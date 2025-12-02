#pragma once

#include "filebrowserwidget.h"
#include "musicplayerwidget.h"
#include <QMainWindow>
#include <QSplitter>

class MainWindow : public QMainWindow {
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = nullptr);
    ~MainWindow() override;

private slots:
    void onFileDoubleClicked(const QString &filePath);
    void onPlaylistGenerated(const QVector<QString> &playlist);

private:
    void setupUI();
    void setupConnections();

private:
    QSplitter *m_splitter;
    FileBrowserWidget *m_fileBrowser;
    MusicPlayerWidget *m_musicPlayer;
};