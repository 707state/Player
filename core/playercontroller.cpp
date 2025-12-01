#include "core/playercontroller.h"
#include <QDirIterator>
#include <QRandomGenerator>
#include <QTime>

#include <taglib/mpegfile.h>
#include <taglib/id3v2tag.h>
#include <taglib/attachedpictureframe.h>
#include <taglib/flacfile.h>
#include <taglib/flacpicture.h>

PlayerController::PlayerController(QObject *parent) : QObject(parent) {
  m_player = new QMediaPlayer(this);
  m_audio = new QAudioOutput(this);
  m_player->setAudioOutput(m_audio);
  m_audio->setVolume(0.7);

  connect(m_player, &QMediaPlayer::positionChanged, this, [&](qint64 p) {
    m_position = p;
    emit positionChanged();
  });

  connect(m_player, &QMediaPlayer::durationChanged, this, [&](qint64 d) {
    m_duration = d;
    emit durationChanged();
  });

  connect(m_player, &QMediaPlayer::mediaStatusChanged, this,
          [&](QMediaPlayer::MediaStatus s) {
            if (s == QMediaPlayer::EndOfMedia)
              playNext();
          });
}

void PlayerController::loadDirectory(const QString &path) {
    qDebug()<<"Folder path is: "<<path;
  m_playlist.clear();
  QDirIterator it(QUrl(path).toLocalFile(), {"*.mp3", "*.flac"}, QDir::Files,
                  QDirIterator::Subdirectories);
  while (it.hasNext()){
    m_playlist << it.next();
    qDebug()<<m_playlist.back();
  }

  m_currentIndex = -1;
}

void PlayerController::togglePlay() {
  if (m_player->playbackState() == QMediaPlayer::PlayingState)
    m_player->pause();
  else
    m_player->play();
  emit playingChanged();
}

void PlayerController::playNext() {
  if (m_playlist.isEmpty())
    return;

  if (m_playMode == PlayMode::Random)
    m_currentIndex = QRandomGenerator::global()->bounded(m_playlist.size());
  else {
    m_currentIndex++;
    if (m_currentIndex >= m_playlist.size())
      m_currentIndex = 0;
  }

  playAt(m_currentIndex);
}

void PlayerController::seek(int value) {
  if (!m_duration)
    return;
  m_player->setPosition((value * m_duration) / 1000);
}

void PlayerController::playAt(int index) {
  if (index < 0 || index >= m_playlist.size())
    return;

  auto path = m_playlist[index];
  m_player->setSource(QUrl::fromLocalFile(path));
  m_player->play();
  emit playingChanged();

  m_cover = loadCoverArtwork(path);
  emit coverChanged();
}

bool PlayerController::playing() const {
  return m_player->playbackState() == QMediaPlayer::PlayingState;
}

QString PlayerController::timeText() const {
  QTime c(0, 0);
  c = c.addMSecs(m_position);
  QTime t(0, 0);
  t = t.addMSecs(m_duration);
  return c.toString("mm:ss") + " / " + t.toString("mm:ss");
}

QUrl PlayerController::cover() const {
  if (m_cover.isNull())
    return {};
  QString path = QDir::temp().filePath("cover.png");
  m_cover.save(path);
  return QUrl::fromLocalFile(path);
}

QPixmap PlayerController::loadCoverArtwork(const QString &filePath) {
  if (filePath.endsWith(".mp3")) {
    TagLib::MPEG::File file(filePath.toStdString().c_str());
    auto *tag = file.ID3v2Tag();
    if (!tag)
      return {};
    auto frames = tag->frameListMap()["APIC"];
    if (frames.isEmpty())
      return {};
    auto *frame =
        static_cast<TagLib::ID3v2::AttachedPictureFrame *>(frames.front());

    QPixmap pix;
    pix.loadFromData(QByteArray(frame->picture().data(),
                                frame->picture().size()));
    return pix;
  }

  if (filePath.endsWith(".flac")) {
    TagLib::FLAC::File file(filePath.toStdString().c_str());
    if (file.pictureList().isEmpty())
      return {};
    auto *pic = file.pictureList().front();
    QPixmap pix;
    pix.loadFromData(QByteArray(pic->data().data(), pic->data().size()));
    return pix;
  }

  return {};
}
