import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import QtQuick.Dialogs

Window {
  width: 900
  height: 600
  visible: true
  title: "QML Music Player"

  ColumnLayout {
    anchors.fill: parent
    spacing: 10

    Image {
      source: player.cover
      fillMode: Image.PreserveAspectFit
      Layout.fillWidth: true
      Layout.fillHeight: true
      smooth: true
    }

    Slider {
      from: 0
      to: 1000
      onMoved: player.seek(value)
    }

    Text {
      text: player.timeText
      horizontalAlignment: Text.AlignHCenter
      Layout.fillWidth: true
    }

    RowLayout {
      Layout.fillWidth: true

      Button {
        text: player.playing ? "❚❚" : "▶"
        onClicked: player.togglePlay()
      }

      Button {
        text: "Next ▶▶"
        onClicked: player.playNext()
      }

      Button {
        text: "Load Dir"
        onClicked: folderDialog.open()
      }
    }
  }

  FolderDialog {
    id: folderDialog
    title: "Select Music Directory"
    onAccepted: player.loadDirectory(selectedFolder)
  }
}
