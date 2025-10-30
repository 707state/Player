//
//  AppDelegate.swift
//  Player
//
//  Created by jask on 2025/9/20.
//
import Foundation
import AppKit
class AppDelegate: NSObject, NSApplicationDelegate {
    var statusItem: NSStatusItem!
    var player: MusicModel?

    func applicationDidFinishLaunching(_ notification: Notification) {
        DispatchQueue.main.async {
            self.statusItem = NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
            if let button = self.statusItem.button {
                button.image = NSImage(systemSymbolName: "music.note", accessibilityDescription: "Player")
                button.action = #selector(self.toggleMenu(_:))
                button.target = self
            }

            let menu = NSMenu()
            menu.addItem(NSMenuItem(title: "Select Directory", action: #selector(self.selectDirectory), keyEquivalent: "s"))
            menu.addItem(NSMenuItem(title: "Play/Pause", action: #selector(self.togglePlayPause), keyEquivalent: "t"))
            menu.addItem(NSMenuItem(title: "Next", action: #selector(self.playNext), keyEquivalent: "n"))
            menu.addItem(NSMenuItem.separator())
            menu.addItem(NSMenuItem(title: "Quit", action: #selector(self.quit), keyEquivalent: "q"))

            self.statusItem.menu = menu
        }
    }

    @objc func toggleMenu(_ sender: Any?) {
        statusItem.menu?.popUp(positioning: nil, at: NSEvent.mouseLocation, in: nil)
    }

    @objc func togglePlayPause() {
        player?.isPlaying ?? false ? player?.pause() : player?.resume()
    }

    @objc func playNext() {
        player?.playNext()
    }
    @objc func selectDirectory() {
        let panel = NSOpenPanel()
        panel.canChooseDirectories = true
        panel.canChooseFiles = false
        panel.allowsMultipleSelection = false
        panel.title = "Select Music Directory"
        panel.begin { response in
                guard response == .OK, let url = panel.url else { return }
                // Asynchronous
                DispatchQueue.global(qos: .userInitiated).async {
                    DispatchQueue.main.async {
                        self.player?.loadDirectoryTree(url: url)
                    }
                }
            }
    }

    @objc func quit() {
        NSApplication.shared.terminate(nil)
    }
}
