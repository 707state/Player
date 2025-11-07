//
//  PlayerApp.swift
//  Player
//
//  Created by jask on 2025/9/18.
//

import SwiftUI

@main
struct PlayerApp: App {
    @StateObject private var player = MusicModel()
    @NSApplicationDelegateAdaptor(AppDelegate.self) var appDelegate
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(player)
                .onAppear {                    
                    appDelegate.player = player
                }
        }
        .commands{
            CommandMenu("Playback"){
                Button(player.isPlaying ? "Pause" : "Play"){
                    if player.isPlaying{
                        player.pause()
                    }else{
                        player.resume()
                    }
                }
                .keyboardShortcut("p",modifiers: [.command])
            }
            CommandMenu("Next"){
                Button("Next"){
                    player.playNext()
                }
                .keyboardShortcut("n",modifiers: [.command])
            }
        }
    }
}
