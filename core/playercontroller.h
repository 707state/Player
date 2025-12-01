#pragma once
#include <QObject>
#include <QMediaPlayer>
#include <QAudioOutput>
#include <QPixmap>
#include <QVector>

enum class PlayMode { Orderly, Random };

class PlayerController : public QObject {
  Q_OBJECT
  Q_PROPERTY(qint64 position READ position NOTIFY positionChanged)
  Q_PROPERTY(qint64 duration READ duration NOTIFY durationChanged)
  Q_PROPERTY(QString timeText READ timeText NOTIFY positionChanged)
  Q_PROPERTY(QUrl cover READ cover NOTIFY coverChanged)
  Q_PROPERTY(bool playing READ playing NOTIFY playingChanged)

public:
  explicit PlayerController(QObject *parent = nullptr);

  Q_INVOKABLE void loadDirectory(const QString &path);
  Q_INVOKABLE void togglePlay();
  Q_INVOKABLE void playNext();
  Q_INVOKABLE void seek(int value);

  qint64 position() const { return m_position; }
  qint64 duration() const { return m_duration; }
  bool playing() const;
  QString timeText() const;
  QUrl cover() const;

signals:
  void positionChanged();
  void durationChanged();
  void coverChanged();
  void playingChanged();

private:
  void playAt(int index);
  QPixmap loadCoverArtwork(const QString &filePath);

  QVector<QString> m_playlist;
  int m_currentIndex = -1;
  PlayMode m_playMode = PlayMode::Orderly;

  QMediaPlayer *m_player;
  QAudioOutput *m_audio;

  QPixmap m_cover;
  qint64 m_position = 0;
  qint64 m_duration = 0;
};
