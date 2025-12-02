#include "mainwindow.h"
#include <QSplitter>

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , m_splitter(new QSplitter(this))
    , m_fileBrowser(new FileBrowserWidget(this))
    , m_musicPlayer(new MusicPlayerWidget(this))
{
    setupUI();
    setupConnections();
    
    // Set initial window size
    resize(1000, 600);
}

MainWindow::~MainWindow() = default;

void MainWindow::setupUI()
{
    // Create splitter for left and right panels
    m_splitter->setOrientation(Qt::Horizontal);
    
    // Add widgets to splitter
    m_splitter->addWidget(m_fileBrowser);
    m_splitter->addWidget(m_musicPlayer);
    
    // Set initial splitter sizes (left: 40%, right: 60%)
    m_splitter->setSizes({400, 600});
    
    // Set splitter as central widget
    setCentralWidget(m_splitter);
    
    // Set window title
    setWindowTitle("Music Player");
}

void MainWindow::setupConnections()
{
    // Connect file browser signals to music player
    connect(m_fileBrowser, &FileBrowserWidget::fileDoubleClicked,
            this, &MainWindow::onFileDoubleClicked);
    
    connect(m_fileBrowser, &FileBrowserWidget::playlistGenerated,
            this, &MainWindow::onPlaylistGenerated);
}

void MainWindow::onFileDoubleClicked(const QString &filePath)
{
    // Play the double-clicked file
    m_musicPlayer->playFile(filePath);
}

void MainWindow::onPlaylistGenerated(const QVector<QString> &playlist)
{
    // Update music player with the new playlist
    m_musicPlayer->setPlaylist(playlist);
}