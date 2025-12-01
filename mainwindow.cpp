#include "mainwindow.h"
#include <QAudioOutput>
#include <QDirIterator>
#include <QFileDialog>
#include <QLabel>
#include <QSplitter>
#include <QToolBar>
#include <QVBoxLayout>
#include <taglib.h>
#include <taglib/attachedpictureframe.h>
#include <taglib/fileref.h>
#include <taglib/flacfile.h>
#include <taglib/flacpicture.h>
#include <taglib/id3v2.h>
#include <taglib/id3v2frame.h>
#include <taglib/id3v2tag.h>
#include <taglib/mpegfile.h>
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
  m_coverLabel = new QLabel(this);
  m_coverLabel->setAlignment(Qt::AlignCenter);
  m_coverLabel->setMinimumSize(240, 240);
  m_coverLabel->setStyleSheet("background:#222;color:#aaa");
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
  m_playPauseBtn = new QPushButton("П", this);
  m_nextBtn = new QPushButton("▶|", this);
  m_modeBtn = new QPushButton("➡ Order", this);

  btnLayout->addWidget(m_playPauseBtn);
  btnLayout->addWidget(m_nextBtn);
  btnLayout->addWidget(m_modeBtn);

  playerLayout->addLayout(btnLayout);

  splitter->addWidget(playerPanel);

  //  Audio Engine
  m_player = new QMediaPlayer(this);
  m_audioOutput = new QAudioOutput(this);
  m_player->setAudioOutput(m_audioOutput);
  m_audioOutput->setVolume(0.7);

  connect(m_playPauseBtn, &QPushButton::clicked, this,
          &MainWindow::onPlayPause);
  connect(m_nextBtn, &QPushButton::clicked, this, &MainWindow::onNextSong);
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
  connect(m_player, &QMediaPlayer::mediaStatusChanged, this,
          [this](QMediaPlayer::MediaStatus status) {
            if (status == QMediaPlayer::EndOfMedia) {
              onNextSong();
            }
          });
  connect(m_modeBtn, &QPushButton::clicked, this, [this]() {
    if (m_playMode == PlayMode::Orderly) {
      m_playMode = PlayMode::Random;
      m_modeBtn->setText("🔀 Random");
    } else {
      m_playMode = PlayMode::Orderly;
      m_modeBtn->setText("➡ Order");
    }
  });
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
  m_playlist.clear();
  QDirIterator it(dir, {"*.mp3", "*.flac"}, QDir::Files,
                  QDirIterator::Subdirectories);
  while (it.hasNext()) {
    m_playlist << it.next();
    qDebug() << m_playlist.back() << "\n";
  }
  m_currentIndex = -1;
  qDebug() << "Loaded songs: " << m_playlist.size();
}

//  Double-click to play
void MainWindow::onFileDoubleClicked(const QModelIndex &index) {
  QString path = m_fsModel->filePath(index);
  if (!path.endsWith(".mp3") && !path.endsWith(".flac"))
    return;
  qDebug() << path;
  m_currentIndex = m_playlist.indexOf(path);
  if (m_currentIndex < 0)
    return;
  playSongAtIndex(m_currentIndex);
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

void MainWindow::onNextSong() {
  if (m_playlist.isEmpty())
    return;

  if (m_playMode == PlayMode::Random) {
    m_currentIndex = QRandomGenerator::global()->bounded(m_playlist.size());
  } else {
    m_currentIndex++;
    if (m_currentIndex >= m_playlist.size())
      m_currentIndex = 0;
  }
  playSongAtIndex(m_currentIndex);
}

QPixmap MainWindow::loadCoverArtwork(const QString &filePath) {
  // ---- MP3 ----
  if (filePath.endsWith(".mp3", Qt::CaseInsensitive)) {
    TagLib::MPEG::File file(filePath.toStdString().c_str());
    auto *tag = file.ID3v2Tag();
    if (!tag)
      return {};

    auto frames = tag->frameListMap()["APIC"];
    if (frames.isEmpty())
      return {};

    auto *frame =
        static_cast<TagLib::ID3v2::AttachedPictureFrame *>(frames.front());

    QByteArray imgData(frame->picture().data(), frame->picture().size());

    QPixmap pix;
    pix.loadFromData(imgData);
    return pix;
  }
  // ---- FLAC ----
  else if (filePath.endsWith(".flac", Qt::CaseInsensitive)) {
    TagLib::FLAC::File file(filePath.toStdString().c_str());
    auto pics = file.pictureList();
    if (pics.isEmpty())
      return {};

    auto *pic = pics.front();
    QByteArray imgData(pic->data().data(), pic->data().size());

    QPixmap pix;
    pix.loadFromData(imgData);
    return pix;
  }

  return {};
}

void MainWindow::updateCoverDisplay() {
  if (m_coverPixmapOriginal.isNull()) {
    m_coverLabel->clear();
    return;
  }

  QSize targetSize = m_coverLabel->size();
  if (targetSize.isEmpty())
    return;

  QPixmap scaled = m_coverPixmapOriginal.scaled(targetSize, Qt::KeepAspectRatio,
                                                Qt::SmoothTransformation);

  m_coverLabel->setPixmap(scaled);
}

void MainWindow::playSongAtIndex(int index) {
  if (index < 0 || index >= m_playlist.size())
    return;

  m_currentIndex = index;
  QString path = m_playlist[m_currentIndex];

  // --- Audio ---
  m_player->setSource(QUrl::fromLocalFile(path));
  m_player->play();
  m_playPauseBtn->setText("❚❚");

  m_coverPixmapOriginal = loadCoverArtwork(path);

  if (!m_coverPixmapOriginal.isNull()) {
    updateCoverDisplay(); // scale from ORIGINAL (no blur)
  } else {
    m_coverLabel->setText("No Cover");
  }
}
