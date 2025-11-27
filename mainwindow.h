#pragma once
#include <QFileSystemModel>
#include <QMainWindow>
#include <QTreeView>
#include <QMediaPlayer>
#include <QLabel>
#include <QPushButton>
class MainWindow : public QMainWindow {
  Q_OBJECT
public:
  explicit MainWindow(QWidget *parent = nullptr);
  ~MainWindow() override;
private slots:
  // Click button
  void onLoadDirectory();
  // play music
  void onFileDoubleClicked(const QModelIndex& index);

  void onPlayPause();
  void onPositionChanged(qint64 pos);
  void onDurationChanged(qint64 dur);
  void onSeek(int value);
private:
    // File Browser
  QFileSystemModel *m_fsModel;
  QTreeView *m_treeView;
  QAction *m_loadDirAction;
  // Player Backend
  QMediaPlayer *m_player;
  QAudioOutput *m_audioOutput;
  // Player UI
  QLabel *m_coverLabel;
  QLabel *m_timeLabel;
  QSlider *m_progress;
  QPushButton *m_playPauseBtn;
  QPushButton *m_randomNextBtn;

  qint64 m_duration=0;
  bool m_userSeeking = false;
};
