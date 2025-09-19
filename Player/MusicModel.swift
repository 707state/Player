//
//  MusicModel.swift
//  Player
//
//  Created by jask on 2025/9/18.
//
import AVFoundation
import Foundation

struct FileNode: Identifiable {
    let id = UUID()
    let url: URL
    var isDirectory: Bool
    var children: [FileNode]?
}
class MusicModel: ObservableObject{
    @Published var rootNode: FileNode?
    @Published var currentFile: URL?
    @Published var isPlaying = false
    @Published var currentTime: TimeInterval = 0
    @Published var duration: TimeInterval = 0
    private var player: AVAudioPlayer?
    private var timer: Timer?
    
    // Recursive
    func loadDirectoryTree(url: URL) {
        rootNode = buildNode(for: url)
    }
    
    private func buildNode(for url: URL) -> FileNode {
        var children: [FileNode] = []
        var isDir: ObjCBool = false
        FileManager.default.fileExists(atPath: url.path, isDirectory: &isDir)
        
        if isDir.boolValue {
            if let items = try? FileManager.default.contentsOfDirectory(
                at: url,
                includingPropertiesForKeys: nil,
                options: [.skipsHiddenFiles]) {
                
                for item in items {
                    // 只收录文件夹 或 mp3/flac 文件
                    var isChildDir: ObjCBool = false
                    FileManager.default.fileExists(atPath: item.path, isDirectory: &isChildDir)
                    if isChildDir.boolValue {
                        children.append(buildNode(for: item))
                    } else {
                        let ext = item.pathExtension.lowercased()
                        if ext == "mp3" || ext == "flac" {
                            children.append(buildNode(for: item))
                        }
                    }
                }
            }
        }
        return FileNode(url: url, isDirectory: isDir.boolValue,
                        children: children.isEmpty ? nil : children)
    }
    func play(file: URL) {
        stop()
        do {
            player = try AVAudioPlayer(contentsOf: file)
            player?.prepareToPlay()
            player?.play()
            duration = player?.duration ?? 0
            currentFile = file
            isPlaying = true
            startTimer()
        } catch {
            print("Failed to play:", error)
        }
    }
    
    func pause() {
        player?.pause()
        isPlaying = false
    }
    
    func resume() {
        player?.play()
        isPlaying = true
    }
    
    func stop() {
        player?.stop()
        isPlaying = false
        currentTime = 0
        stopTimer()
    }
    
    func seek(to time: TimeInterval) {
        player?.currentTime = time
        currentTime = time
    }
    
    private func startTimer() {
        timer = Timer.scheduledTimer(withTimeInterval: 0.5, repeats: true) { _ in
            self.currentTime = self.player?.currentTime ?? 0
        }
    }
    private func stopTimer() {
        timer?.invalidate()
        timer = nil
    }
}
