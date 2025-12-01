#pragma once
#include <QFileSystemModel>
#include <QLabel>
#include <QMainWindow>
#include <QMediaPlayer>
#include <QPushButton>
#include <QRandomGenerator64>
#include <QTreeView>
enum class PlayMode { Orderly, Random };

class MainWindow : public QMainWindow {
  Q_OBJECT
public:
  explicit MainWindow(QWidget *parent = nullptr);
  ~MainWindow() override;
  void resizeEvent(QResizeEvent *event) override {
    QMainWindow::resizeEvent(event);
    updateCoverDisplay();
  }
private slots:
  // Click button
  void onLoadDirectory();
  // play music
  void onFileDoubleClicked(const QModelIndex &index);

  void onPlayPause();
  void onPositionChanged(qint64 pos);
  void onDurationChanged(qint64 dur);
  void onSeek(int value);
  void onNextSong();

private:
  QPixmap loadCoverArtwork(const QString &filePath);
  void updateCoverDisplay();
  void playSongAtIndex(int);

private:
  // File Browser
  QFileSystemModel *m_fsModel;
  QTreeView *m_treeView;
  QAction *m_loadDirAction;
  // playlist
  QVector<QString> m_playlist;
  int m_currentIndex = -1;
  PlayMode m_playMode = PlayMode::Orderly;
  // Player Backend
  QMediaPlayer *m_player;
  QAudioOutput *m_audioOutput;
  // Player UI
  QLabel *m_coverLabel;
  QPixmap m_coverPixmapOriginal;
  QLabel *m_timeLabel;
  QSlider *m_progress;
  QPushButton *m_playPauseBtn;
  QPushButton *m_nextBtn;
  QPushButton *m_modeBtn;

  qint64 m_duration = 0;
  bool m_userSeeking = false;
};
