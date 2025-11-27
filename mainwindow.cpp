#include "mainwindow.h"
#include <QAudioOutput>
#include <QFileDialog>
#include <QLabel>
#include <QSplitter>
#include <QToolBar>
#include <QVBoxLayout>
MainWindow::MainWindow(QWidget *parent) : QMainWindow(parent) {
  // tool bar with "load directory" button
  auto *toolbar = addToolBar("Main Toolbar");
  m_loadDirAction = new QAction("Load Directory", this);
  toolbar->addAction(m_loadDirAction);
  // bind method
  connect(m_loadDirAction, &QAction::triggered, this,
          &MainWindow::onLoadDirectory);
  // left: file tree view
  auto *splitter = new QSplitter(this);
  setCentralWidget(splitter);
  m_fsModel = new QFileSystemModel(this);
  m_fsModel->setRootPath(QDir::homePath());

  m_fsModel->setNameFilters({"*.mp3", "*.flac"});
  m_fsModel->setNameFilterDisables(false);
  m_treeView = new QTreeView(this);
  m_treeView->setModel(m_fsModel);
  m_treeView->setRootIndex(m_fsModel->index(QDir::homePath()));
  m_treeView->setHeaderHidden(true);
  m_treeView->setAnimated(true);
  connect(m_treeView, &QTreeView::doubleClicked, this,
          &MainWindow::onFileDoubleClicked);
  splitter->addWidget(m_treeView);
  // Right: player panel
  auto *playerPanel = new QWidget(this);
  auto *playerLayout = new QVBoxLayout(playerPanel);
  // Album artwork
  m_coverLabel = new QLabel("No Cover", this);
  m_coverLabel->setAlignment(Qt::AlignCenter);
  m_coverLabel->setMinimumHeight(200);
  playerLayout->addWidget(m_coverLabel);

  // Progress slider
  m_progress = new QSlider(Qt::Horizontal, this);
  m_progress->setRange(0, 1000);
  playerLayout->addWidget(m_progress);

  // Time label
  m_timeLabel = new QLabel("00:00 / 00:00", this);
  m_timeLabel->setAlignment(Qt::AlignCenter);
  playerLayout->addWidget(m_timeLabel);

  // Buttons
  auto *btnLayout = new QHBoxLayout();
  m_playPauseBtn = new QPushButton("▶", this);
  m_randomNextBtn = new QPushButton("🔀 Next", this);

  btnLayout->addWidget(m_playPauseBtn);
  btnLayout->addWidget(m_randomNextBtn);

  playerLayout->addLayout(btnLayout);

  splitter->addWidget(playerPanel);

  //  Audio Engine
  m_player = new QMediaPlayer(this);
  m_audioOutput = new QAudioOutput(this);
  m_player->setAudioOutput(m_audioOutput);
  m_audioOutput->setVolume(0.7);

  connect(m_playPauseBtn, &QPushButton::clicked, this,
          &MainWindow::onPlayPause);

  connect(m_progress, &QSlider::sliderMoved, this, &MainWindow::onSeek);
  connect(m_progress, &QSlider::sliderPressed,
          [this]() { m_userSeeking = true; });
  connect(m_progress, &QSlider::sliderReleased, [this]() {
    if (!m_duration)
      return;
    m_player->setPosition((m_progress->value() * m_duration) / 1000);
    m_userSeeking = false;
  });
  connect(m_player, &QMediaPlayer::positionChanged, this,
          &MainWindow::onPositionChanged);

  connect(m_player, &QMediaPlayer::durationChanged, this,
          &MainWindow::onDurationChanged);
}
MainWindow::~MainWindow() = default;

void MainWindow::onLoadDirectory() {
  QString dir = QFileDialog::getExistingDirectory(
      this, "Select Music Directory", QDir::homePath(),
      QFileDialog::ShowDirsOnly | QFileDialog::DontResolveSymlinks);
  if (dir.isEmpty()) {
    return;
  }
  m_treeView->setRootIndex(m_fsModel->index(dir));
}

//  Double-click to play
void MainWindow::onFileDoubleClicked(const QModelIndex &index) {
  QString path = m_fsModel->filePath(index);
  if (!path.endsWith(".mp3") && !path.endsWith(".flac"))
    return;

  m_player->setSource(QUrl::fromLocalFile(path));
  m_player->play();
  m_playPauseBtn->setText("❚❚");

  //  Clear cover for now (we’ll add TagLib later)
  m_coverLabel->setText("Playing:\n" + QFileInfo(path).fileName());
}

//  Play / Pause
void MainWindow::onPlayPause() {
  if (m_player->playbackState() == QMediaPlayer::PlayingState) {
    m_player->pause();
    m_playPauseBtn->setText("▶");
  } else {
    m_player->play();
    m_playPauseBtn->setText("❚❚");
  }
}

//  Update progress
void MainWindow::onPositionChanged(qint64 pos) {
  if (!m_duration)
    return;
  if (!m_userSeeking) {
    m_progress->setValue(static_cast<int>((pos * 1000) / m_duration));
  }
  QTime curr(0, 0);
  curr = curr.addMSecs(pos);

  QTime total(0, 0);
  total = total.addMSecs(m_duration);

  m_timeLabel->setText(curr.toString("mm:ss") + " / " +
                       total.toString("mm:ss"));
}

void MainWindow::onDurationChanged(qint64 dur) { m_duration = dur; }

//  Seek
void MainWindow::onSeek(int value) {
  if (!m_duration)
    return;
  m_player->setPosition((value * m_duration) / 1000);
}
