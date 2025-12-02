#pragma once

#include "common.h"
#include <QFileSystemModel>
#include <QTreeView>
#include <QWidget>

class FileBrowserWidget : public QWidget {
    Q_OBJECT

public:
    explicit FileBrowserWidget(QWidget *parent = nullptr);
    ~FileBrowserWidget() override;

    void loadDirectory(const QString &path);
    QString currentDirectory() const;
    void setFilter(const QStringList &filters);

signals:
    void directoryLoaded(const QString &directory);
    void fileDoubleClicked(const QString &filePath);
    void playlistGenerated(const QVector<QString> &playlist);

private slots:
    void onFileDoubleClicked(const QModelIndex &index);

private:
    QFileSystemModel *m_fsModel;
    QTreeView *m_treeView;
    QVector<QString> m_playlist;
};