#include "musicplayerwidget.h"
#include <QAudioOutput>
#include <QDirIterator>
#include <QFileDialog>
#include <QHBoxLayout>
#include <QLabel>
#include <QMediaPlayer>
#include <QPushButton>
#include <QRandomGenerator>
#include <QSlider>
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

MusicPlayerWidget::MusicPlayerWidget(QWidget *parent)
    : QWidget(parent)
    , m_player(new QMediaPlayer(this))
    , m_audioOutput(new QAudioOutput(this))
    , m_coverLabel(new QLabel(this))
    , m_songInfo(new QLabel(this))
    , m_timeLabel(new QLabel("00:00 / 00:00", this))
    , m_progress(new QSlider(Qt::Horizontal, this))
    , m_playPauseBtn(new QPushButton("▶", this))
    , m_nextBtn(new QPushButton("▶|", this))
    , m_prevBtn(new QPushButton("|◀", this))
    , m_modeBtn(new QPushButton("➡ Order", this))
    , m_volumeSlider(new QSlider(Qt::Horizontal, this))
{
    setupUI();
    setupConnections();

    // Setup audio
    m_player->setAudioOutput(m_audioOutput);
    m_audioOutput->setVolume(0.7);
    m_volumeSlider->setValue(70);
}

MusicPlayerWidget::~MusicPlayerWidget() = default;

void MusicPlayerWidget::setupUI()
{
    auto *mainLayout = new QVBoxLayout(this);

    // Album artwork
    m_coverLabel->setAlignment(Qt::AlignCenter);
    m_coverLabel->setMinimumSize(240, 240);
    m_coverLabel->setStyleSheet("background:#222;color:#aaa");
    m_coverLabel->setAlignment(Qt::AlignCenter);
    m_coverLabel->setMinimumHeight(200);
    mainLayout->addWidget(m_coverLabel);

    // Song info
    m_songInfo->setAlignment(Qt::AlignCenter);
    m_songInfo->setStyleSheet("color:#aaa; font-weight:bold; font-size:16px;");
    m_songInfo->setText("No song selected");
    mainLayout->addWidget(m_songInfo);

    // Progress slider
    m_progress->setRange(0, 1000);
    mainLayout->addWidget(m_progress);

    // Time label
    m_timeLabel->setAlignment(Qt::AlignCenter);
    mainLayout->addWidget(m_timeLabel);

    // Control buttons
    auto *controlLayout = new QHBoxLayout();
    controlLayout->addWidget(m_prevBtn);
    controlLayout->addWidget(m_playPauseBtn);
    controlLayout->addWidget(m_nextBtn);
    controlLayout->addWidget(m_modeBtn);

    // Volume control
    auto *volumeLayout = new QHBoxLayout();
    auto *volumeLabel = new QLabel("Volume:", this);
    m_volumeSlider->setRange(0, 100);
    m_volumeSlider->setValue(70);
    volumeLayout->addWidget(volumeLabel);
    volumeLayout->addWidget(m_volumeSlider);

    auto *bottomLayout = new QVBoxLayout();
    bottomLayout->addLayout(controlLayout);
    bottomLayout->addLayout(volumeLayout);

    mainLayout->addLayout(bottomLayout);
}

void MusicPlayerWidget::setupConnections()
{
    connect(m_playPauseBtn, &QPushButton::clicked,
            this, &MusicPlayerWidget::onPlayPause);
    connect(m_nextBtn, &QPushButton::clicked,
            this, &MusicPlayerWidget::onNext);
    connect(m_prevBtn, &QPushButton::clicked,
            this, &MusicPlayerWidget::onPrevious);
    connect(m_modeBtn, &QPushButton::clicked,
            this, &MusicPlayerWidget::onModeChanged);

    connect(m_progress, &QSlider::sliderMoved,
            this, &MusicPlayerWidget::onSeek);
    connect(m_progress, &QSlider::sliderPressed,
            [this]() { m_userSeeking = true; });
    connect(m_progress, &QSlider::sliderReleased, [this]() {
        if (!m_duration) return;
        m_player->setPosition((m_progress->value() * m_duration) / 1000);
        m_userSeeking = false;
    });

    connect(m_volumeSlider, &QSlider::valueChanged, [this](int value) {
        m_audioOutput->setVolume(value / 100.0f);
    });

    connect(m_player, &QMediaPlayer::positionChanged,
            this, &MusicPlayerWidget::onPositionChanged);
    connect(m_player, &QMediaPlayer::durationChanged,
            this, &MusicPlayerWidget::onDurationChanged);
    connect(m_player, &QMediaPlayer::mediaStatusChanged,
            this, &MusicPlayerWidget::onMediaStatusChanged);
}

void MusicPlayerWidget::playFile(const QString &filePath)
{
    if (filePath.isEmpty()) {
        return;
    }

    // Find the file in playlist
    int index = m_playlist.indexOf(filePath);
    if (index >= 0) {
        playSongAtIndex(index);
    } else {
        // Play as single file
        m_player->setSource(QUrl::fromLocalFile(filePath));
        m_player->play();
        m_playPauseBtn->setText("❚❚");

        // Extract metadata
        BasicMeta meta = extractBasicMetadata(filePath);
        updateCoverWithMeta(meta);

        emit playStateChanged(true);
        emit currentSongChanged(filePath);
        emit metaDataChanged(meta);
    }
}

void MusicPlayerWidget::play()
{
    m_player->play();
    m_playPauseBtn->setText("❚❚");
    emit playStateChanged(true);
}

void MusicPlayerWidget::pause()
{
    m_player->pause();
    m_playPauseBtn->setText("▶");
    emit playStateChanged(false);
}

void MusicPlayerWidget::stop()
{
    m_player->stop();
    m_playPauseBtn->setText("▶");
    emit playStateChanged(false);
}

void MusicPlayerWidget::next()
{
    onNext();
}

void MusicPlayerWidget::previous()
{
    onPrevious();
}

void MusicPlayerWidget::setVolume(float volume)
{
    m_audioOutput->setVolume(volume);
    m_volumeSlider->setValue(volume * 100);
}

float MusicPlayerWidget::volume() const
{
    return m_audioOutput->volume();
}

void MusicPlayerWidget::setPlayMode(PlayMode mode)
{
    if (m_playMode != mode) {
        m_playMode = mode;
        updateModeButtonText();
        emit playModeChanged(mode);
    }
}

PlayMode MusicPlayerWidget::playMode() const
{
    return m_playMode;
}

void MusicPlayerWidget::setPlaylist(const QVector<QString> &playlist)
{
    m_playlist = playlist;
    m_currentIndex = -1;
}

QVector<QString> MusicPlayerWidget::playlist() const
{
    return m_playlist;
}

void MusicPlayerWidget::setCurrentIndex(int index)
{
    if (index >= 0 && index < m_playlist.size()) {
        playSongAtIndex(index);
    }
}

int MusicPlayerWidget::currentIndex() const
{
    return m_currentIndex;
}

bool MusicPlayerWidget::isPlaying() const
{
    return m_player->playbackState() == QMediaPlayer::PlayingState;
}

qint64 MusicPlayerWidget::duration() const
{
    return m_player->duration();
}

qint64 MusicPlayerWidget::position() const
{
    return m_player->position();
}

void MusicPlayerWidget::seek(qint64 position)
{
    m_player->setPosition(position);
}

void MusicPlayerWidget::seekToPercentage(float percentage)
{
    if (percentage < 0.0f) percentage = 0.0f;
    if (percentage > 1.0f) percentage = 1.0f;

    qint64 pos = static_cast<qint64>(m_duration * percentage);
    m_player->setPosition(pos);
}

void MusicPlayerWidget::onPlayPause()
{
    if (m_player->playbackState() == QMediaPlayer::PlayingState) {
        pause();
    } else {
        play();
    }
}

void MusicPlayerWidget::onNext()
{
    if (m_playlist.isEmpty()) {
        return;
    }

    if (m_playMode == PlayMode::Random) {
        m_currentIndex = QRandomGenerator::global()->bounded(m_playlist.size());
    } else {
        m_currentIndex++;
        if (m_currentIndex >= m_playlist.size()) {
            m_currentIndex = 0;
        }
    }

    playSongAtIndex(m_currentIndex);
}

void MusicPlayerWidget::onPrevious()
{
    if (m_playlist.isEmpty()) {
        return;
    }

    if (m_playMode == PlayMode::Random) {
        m_currentIndex = QRandomGenerator::global()->bounded(m_playlist.size());
    } else {
        m_currentIndex--;
        if (m_currentIndex < 0) {
            m_currentIndex = m_playlist.size() - 1;
        }
    }

    playSongAtIndex(m_currentIndex);
}

void MusicPlayerWidget::onModeChanged()
{
    if (m_playMode == PlayMode::Orderly) {
        setPlayMode(PlayMode::Random);
    } else {
        setPlayMode(PlayMode::Orderly);
    }
}

void MusicPlayerWidget::onPositionChanged(qint64 pos)
{
    if (!m_duration) {
        return;
    }

    if (!m_userSeeking) {
        m_progress->setValue(static_cast<int>((pos * 1000) / m_duration));
    }

    QTime curr(0, 0);
    curr = curr.addMSecs(pos);

    QTime total(0, 0);
    total = total.addMSecs(m_duration);

    m_timeLabel->setText(curr.toString("mm:ss") + " / " + total.toString("mm:ss"));
    emit positionChanged(pos);
}

void MusicPlayerWidget::onDurationChanged(qint64 dur)
{
    m_duration = dur;
    emit durationChanged(dur);
}

void MusicPlayerWidget::onSeek(int value)
{
    if (!m_duration) {
        return;
    }

    qint64 pos = (value * m_duration) / 1000;
    m_player->setPosition(pos);
}

void MusicPlayerWidget::onMediaStatusChanged(QMediaPlayer::MediaStatus status)
{
    if (status == QMediaPlayer::EndOfMedia) {
        onNext();
    }
}

void MusicPlayerWidget::updateCoverDisplay()
{
    if (m_coverPixmapOriginal.isNull()) {
        m_coverLabel->clear();
        m_coverLabel->setText("No Cover");
        return;
    }

    QSize targetSize = m_coverLabel->size();
    if (targetSize.isEmpty()) {
        return;
    }

    QPixmap scaled = m_coverPixmapOriginal.scaled(targetSize, Qt::KeepAspectRatio,
                                                  Qt::SmoothTransformation);
    m_coverLabel->setPixmap(scaled);
}

void MusicPlayerWidget::playSongAtIndex(int index)
{
    if (index < 0 || index >= m_playlist.size()) {
        return;
    }

    m_currentIndex = index;
    QString path = m_playlist[m_currentIndex];

    // Play audio
    m_player->setSource(QUrl::fromLocalFile(path));
    m_player->play();
    m_playPauseBtn->setText("❚❚");

    // Extract and display metadata
    BasicMeta meta = extractBasicMetadata(path);
    updateCoverWithMeta(meta);

    emit playStateChanged(true);
    emit currentSongChanged(path);
    emit metaDataChanged(meta);
}

void MusicPlayerWidget::updateCoverWithMeta(const BasicMeta &meta)
{
    if (!meta.artwork.isNull()) {
        m_coverPixmapOriginal = meta.artwork;
        updateCoverDisplay();
    } else {
        m_coverLabel->setText("No Cover");
        m_coverPixmapOriginal = QPixmap();
    }

    QString infoText;
    if (!meta.title.isEmpty() && !meta.artist.isEmpty()) {
        infoText = QString("%1 - %2").arg(meta.artist, meta.title);
    } else if (!meta.title.isEmpty()) {
        infoText = meta.title;
    } else {
        infoText = "Unknown";
    }

    m_songInfo->setText(infoText);
}

void MusicPlayerWidget::updateModeButtonText()
{
    if (m_playMode == PlayMode::Orderly) {
        m_modeBtn->setText("➡ Order");
    } else {
        m_modeBtn->setText("🔀 Random");
    }
}



BasicMeta MusicPlayerWidget::extractBasicMetadata(const QString &filePath)
{
    BasicMeta meta;

    // Fallback: file name as title
    meta.title = QFileInfo(filePath).baseName();
    meta.artist = "";
    meta.album = "";
    meta.artwork = QPixmap();

    if (filePath.endsWith(".mp3", Qt::CaseInsensitive)) {
        TagLib::MPEG::File file(filePath.toStdString().c_str());
        auto *tag = file.ID3v2Tag();
        if (tag) {
            meta.title = QString::fromStdString(tag->title().to8Bit(true));
            meta.artist = QString::fromStdString(tag->artist().to8Bit(true));
            meta.album = QString::fromStdString(tag->album().to8Bit(true));

            auto frames = tag->frameListMap()["APIC"];
            if (!frames.isEmpty()) {
                auto *frame =
                    static_cast<TagLib::ID3v2::AttachedPictureFrame *>(frames.front());
                QByteArray imgData(frame->picture().data(), frame->picture().size());
                meta.artwork.loadFromData(imgData);
            }
        }
    } else if (filePath.endsWith(".flac", Qt::CaseInsensitive)) {
        TagLib::FLAC::File file(filePath.toStdString().c_str());
        auto *tag = file.tag();
        if (tag) {
            meta.title = QString::fromStdString(tag->title().to8Bit(true));
            meta.artist = QString::fromStdString(tag->artist().to8Bit(true));
            meta.album = QString::fromStdString(tag->album().to8Bit(true));
        }

        auto pics = file.pictureList();
        if (!pics.isEmpty()) {
            auto *pic = pics.front();
            QByteArray imgData(pic->data().data(), pic->data().size());
            meta.artwork.loadFromData(imgData);
        }
    }

    return meta;
}

void MusicPlayerWidget::resizeEvent(QResizeEvent *event)
{
    QWidget::resizeEvent(event);
    updateCoverDisplay();
}
