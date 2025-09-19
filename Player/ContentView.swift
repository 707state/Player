//
//  ContentView.swift
//  Player
//
//  Created by jask on 2025/9/18.
//

import SwiftUI
import AVFoundation
struct ContentView: View {
    @StateObject private var player=MusicModel()
    @State private var showPanel=false
    var body: some View {
        HStack{
            List(player.mp3files,id: \.self){
                file in Button(file.lastPathComponent){
                    player.play(file: file)
                }
            }
            .frame(minWidth: 200)
            VStack(spacing:20){
                if let currentItem=player.currentFile{
                    Text("正在播放：\(currentItem.lastPathComponent)")
                        .font(.headline)
                }else{
                    Text("请选择一个mp3/flac文件")
                }
                Slider(value: Binding(
                    get: {player.currentTime},
                    set: {player.seek(to: $0)}
                ),
                       in: 0...(player.duration))
                .disabled(player.currentFile==nil)
                HStack(spacing: 40){
                    Button(action: {
                        if player.isPlaying{
                            player.pause()
                        }else{
                            player.resume()
                        }
                    }){
                        if player.isPlaying{
                            Image(systemName: "pause.fill")
                        }else{
                            Image(systemName: "play.fill")
                        }
                    }
                    Button(action: {player.stop()}){
                        Image(systemName: "stop.fill")
                    }
                }
            }.frame(maxWidth: .infinity)
        }
        .toolbar{
            Button("选择音乐目录"){
                selectDirectory()
            }
        }
    }
    // 打开选择目录对话框
    private func selectDirectory() {
        let panel = NSOpenPanel()
        panel.canChooseDirectories = true
        panel.canChooseFiles = false
        panel.allowsMultipleSelection = false
        if panel.runModal() == .OK, let url = panel.url {
            player.loadMP3Files(url: url)
        }
    }
}

#Preview {
    ContentView()
}
