#include "filebrowserwidget.h"
#include <QDebug>
#include <QDirIterator>
#include <QFileDialog>
#include <QHBoxLayout>
#include <QPushButton>
#include <QVBoxLayout>

FileBrowserWidget::FileBrowserWidget(QWidget *parent)
    : QWidget(parent)
    , m_fsModel(new QFileSystemModel(this))
    , m_treeView(new QTreeView(this))
{
    // Setup layout
    auto *mainLayout = new QVBoxLayout(this);
    
    // Create toolbar with load directory button
    auto *toolbarLayout = new QHBoxLayout();
    auto *loadDirButton = new QPushButton("Load Directory", this);
    toolbarLayout->addWidget(loadDirButton);
    toolbarLayout->addStretch();
    mainLayout->addLayout(toolbarLayout);
    
    // Setup file system model
    m_fsModel->setRootPath(QDir::homePath());
    m_fsModel->setNameFilters({"*.mp3", "*.flac"});
    m_fsModel->setNameFilterDisables(false);
    
    // Setup tree view
    m_treeView->setModel(m_fsModel);
    m_treeView->setRootIndex(m_fsModel->index(QDir::homePath()));
    m_treeView->setHeaderHidden(true);
    m_treeView->setAnimated(true);
    m_treeView->setMinimumWidth(300);
    mainLayout->addWidget(m_treeView);
    
    // Connect signals
    connect(loadDirButton, &QPushButton::clicked, this, [this]() {
        QString dir = QFileDialog::getExistingDirectory(
            this, "Select Music Directory", QDir::homePath(),
            QFileDialog::ShowDirsOnly | QFileDialog::DontResolveSymlinks);
        if (!dir.isEmpty()) {
            loadDirectory(dir);
        }
    });
    
    connect(m_treeView, &QTreeView::doubleClicked,
            this, &FileBrowserWidget::onFileDoubleClicked);
}

FileBrowserWidget::~FileBrowserWidget() = default;

void FileBrowserWidget::loadDirectory(const QString &path)
{
    if (path.isEmpty()) {
        return;
    }
    
    m_treeView->setRootIndex(m_fsModel->index(path));
    m_playlist.clear();
    
    // Scan for music files
    QDirIterator it(path, {"*.mp3", "*.flac"}, QDir::Files,
                    QDirIterator::Subdirectories);
    while (it.hasNext()) {
        m_playlist << it.next();
    }
    
    emit directoryLoaded(path);
    emit playlistGenerated(m_playlist);
    
    qDebug() << "Loaded" << m_playlist.size() << "songs from" << path;
}

QString FileBrowserWidget::currentDirectory() const
{
    return m_fsModel->filePath(m_treeView->rootIndex());
}

void FileBrowserWidget::setFilter(const QStringList &filters)
{
    m_fsModel->setNameFilters(filters);
}

void FileBrowserWidget::onFileDoubleClicked(const QModelIndex &index)
{
    QString path = m_fsModel->filePath(index);
    if (!path.endsWith(".mp3") && !path.endsWith(".flac")) {
        return;
    }
    
    emit fileDoubleClicked(path);
}