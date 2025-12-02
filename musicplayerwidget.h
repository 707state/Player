#pragma once

#include "common.h"
#include <QAudioOutput>
#include <QLabel>
#include <QMediaPlayer>
#include <QPushButton>
#include <QSlider>
#include <QWidget>

enum class PlayMode { Orderly, Random };

class MusicPlayerWidget : public QWidget {
    Q_OBJECT

public:
    explicit MusicPlayerWidget(QWidget *parent = nullptr);
    ~MusicPlayerWidget() override;

    void playFile(const QString &filePath);
    void play();
    void pause();
    void stop();
    void next();
    void previous();
    
    void setVolume(float volume);
    float volume() const;
    
    void setPlayMode(PlayMode mode);
    PlayMode playMode() const;
    
    void setPlaylist(const QVector<QString> &playlist);
    QVector<QString> playlist() const;
    
    void setCurrentIndex(int index);
    int currentIndex() const;
    
    bool isPlaying() const;
    qint64 duration() const;
    qint64 position() const;

signals:
    void playStateChanged(bool playing);
    void currentSongChanged(const QString &filePath);
    void positionChanged(qint64 position);
    void durationChanged(qint64 duration);
    void playModeChanged(PlayMode mode);
    void metaDataChanged(const BasicMeta &meta);

public slots:
    void seek(qint64 position);
    void seekToPercentage(float percentage);

private slots:
    void onPlayPause();
    void onNext();
    void onPrevious();
    void onModeChanged();
    void onPositionChanged(qint64 pos);
    void onDurationChanged(qint64 dur);
    void onSeek(int value);
    void onMediaStatusChanged(QMediaPlayer::MediaStatus status);
    void updateCoverDisplay();

private:
    void setupUI();
    void setupConnections();
    void playSongAtIndex(int index);
    BasicMeta extractBasicMetadata(const QString &filePath);
    void updateCoverWithMeta(const BasicMeta &meta);
    void updateModeButtonText();
    
    // Player backend
    QMediaPlayer *m_player;
    QAudioOutput *m_audioOutput;
    
    // Playlist
    QVector<QString> m_playlist;
    int m_currentIndex = -1;
    PlayMode m_playMode = PlayMode::Orderly;
    
    // UI elements
    QLabel *m_coverLabel;
    QLabel *m_songInfo;
    QPixmap m_coverPixmapOriginal;
    QLabel *m_timeLabel;
    QSlider *m_progress;
    QPushButton *m_playPauseBtn;
    QPushButton *m_nextBtn;
    QPushButton *m_prevBtn;
    QPushButton *m_modeBtn;
    QSlider *m_volumeSlider;
    
    // State
    qint64 m_duration = 0;
    bool m_userSeeking = false;
    
protected:
    void resizeEvent(QResizeEvent *event) override;
};