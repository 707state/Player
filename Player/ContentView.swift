//
//  ContentView.swift
//  Player
//
//  Created by jask on 2025/9/18.
//

import SwiftUI
import AVFoundation

struct ContentView: View {
    @EnvironmentObject var player: MusicModel
    
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
                .frame(minWidth: 150)
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

func formatTime(_ time: TimeInterval) -> String {
    let minutes = Int(time) / 60
    let seconds = Int(time) % 60
    return String(format: "%02d:%02d", minutes, seconds)
}

struct PlayerPanel: View {
    @ObservedObject var player: MusicModel
    
    var body: some View {
        ZStack {
            // 背景层：模糊的 artwork
            if let artwork = player.artwork {
                Image(nsImage: artwork)
                    .resizable()
                    .scaledToFill()
                    .blur(radius: 30)
                    .overlay(
                        LinearGradient(
                            gradient: Gradient(colors: [.black.opacity(0.4), .clear]),
                            startPoint: .bottom,
                            endPoint: .top
                        )
                    )
                    .ignoresSafeArea()
            } else {
                // 默认模糊背景
                Color.gray.opacity(0.2)
                    .ignoresSafeArea()
            }
            
            // 前景层：播放器内容
            VStack(spacing: 20) {
                if let artwork = player.artwork {
                    Image(nsImage: artwork)
                        .resizable()
                        .scaledToFit()
                        .frame(width: 200, height: 200)
                        .cornerRadius(10)
                        .shadow(radius: 4)
                } else {
                    Image(systemName: "music.note")
                        .resizable()
                        .scaledToFit()
                        .frame(width: 120, height: 120)
                        .opacity(0.3)
                }
                
                if let currentItem = player.currentFile {
                    Text("Now playing：\(currentItem.lastPathComponent)")
                        .font(.headline)
                } else {
                    Text("Please select an mp3/flac file")
                }
                
                // 播放进度滑杆
                Slider(
                    value: Binding(
                        get: { player.currentTime },
                        set: { player.seek(to: $0) }
                    ),
                    in: 0...(player.duration)
                )
                .disabled(player.currentFile == nil)
                
                HStack {
                    Text(formatTime(player.currentTime))
                    Spacer()
                    Text(formatTime(player.duration))
                }
                .font(.caption)
                
                HStack(spacing: 40) {
                    Button {
                        if player.isPlaying { player.pause() } else { player.resume() }
                    } label: {
                        Image(systemName: player.isPlaying ? "pause.fill" : "play.fill")
                    }
                    Button { player.stop() } label: {
                        Image(systemName: "stop.fill")
                    }
                    Button { player.playNext() } label: {
                        Image(systemName: "forward.fill")
                    }
                }
                
                // 音量控制
                Slider(value: $player.volume, in: 0...1)
                    .frame(width: 120, height: 20)
                    .padding(.horizontal, 4)
            }
            .padding()
        }
    }
}
