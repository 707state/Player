//
//  ContentView.swift
//  Player
//
//  Created by jask on 2025/9/18.
//

import SwiftUI
import AVFoundation

struct ContentView: View {
    @StateObject private var player = MusicModel()
    
    var body: some View {
        NavigationSplitView {
            // Scrollable file view
            if let root = player.rootNode {
                List {
                    OutlineGroup(root, children: \.children) { node in
                        if node.isDirectory {
                            Label(node.url.lastPathComponent, systemImage: "folder")
                        } else {
                            Button {
                                player.play(file: node.url)
                            } label: {
                                Label(node.url.lastPathComponent, systemImage: "music.note")
                            }
                        }
                    }
                }
                .frame(minWidth: 150)     // 设置最小宽度
            } else {
                Text("Please select music directory")
                    .frame(minWidth: 150)
            }
        } detail: {
            // 右边播放控制区
            PlayerPanel(player: player)
        }
        .toolbar {
            Button("Select music directory") {
                selectDirectory()
            }
        }
    }
    
    private func selectDirectory() {
        let panel = NSOpenPanel()
        panel.canChooseDirectories = true
        panel.canChooseFiles = false
        panel.allowsMultipleSelection = false
        if panel.runModal() == .OK, let url = panel.url {
            player.loadDirectoryTree(url: url)
        }
    }
}

// 将右侧播放器单独封装
struct PlayerPanel: View {
    @ObservedObject var player: MusicModel
    var body: some View {
        VStack(spacing: 20) {
            if let currentItem = player.currentFile {
                Text("Now playing：\(currentItem.lastPathComponent)")
                    .font(.headline)
            } else {
                Text("Please select on mp3/flac file")
            }
            
            Slider(value: Binding(
                get: { player.currentTime },
                set: { player.seek(to: $0) }),
                   in: 0...(player.duration))
            .disabled(player.currentFile == nil)
            
            HStack(spacing: 40) {
                Button {
                    if player.isPlaying { player.pause() } else { player.resume() }
                } label: {
                    Image(systemName: player.isPlaying ? "pause.fill" : "play.fill")
                }
                Button { player.stop() } label: {
                    Image(systemName: "stop.fill")
                }
            }
        }
        .padding()
    }
}
#Preview {
    ContentView()
}
