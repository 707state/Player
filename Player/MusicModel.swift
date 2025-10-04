//
//  MusicModel.swift
//  Player
//
//  Created by jask on 2025/9/18.
//
import AVFoundation
import Foundation
import AppKit
struct FileNode: Identifiable {
    let id = UUID()
    let url: URL
    var isDirectory: Bool
    var children: [FileNode]?
}
class MusicModel: NSObject,ObservableObject, AVAudioPlayerDelegate{
    // Only one MusicModel in this app!
    static let shared = MusicModel()
    
    @Published var rootNode: FileNode?
    @Published var currentFile: URL?
    @Published var isPlaying = false
    @Published var currentTime: TimeInterval = 0
    @Published var duration: TimeInterval = 0
    @Published var volume: Float = 1.0 {
        didSet {
            player?.volume = volume
        }
    }
    @Published var artwork: NSImage?=nil
    
    private var player: AVAudioPlayer?
    private var timer: Timer?
    
    // Flattend file list for next play
    private var allFiles: [URL] = []
    // Recursive
    func loadDirectoryTree(url: URL) {
        let url_=url;
        DispatchQueue.global(qos: .userInitiated).async {
            let root = self.buildNode(for: url_)
            DispatchQueue.main.async {
                self.rootNode = root
                self.allFiles=self.collectAllFiles(from: root)
            }
        }
    }
    // Traverse whole tree
    private func collectAllFiles(from node: FileNode?) -> [URL] {
        guard let node = node else { return [] }
        if node.isDirectory {
            return node.children?.flatMap { collectAllFiles(from: $0) } ?? []
        } else {
            return [node.url]
        }
    }
    
    func buildNode(for url: URL) -> FileNode {
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
            player?.delegate = self
            player?.prepareToPlay()
            player?.play()
            duration = player?.duration ?? 0
            currentFile = file
            isPlaying = true
            Task.detached {
                self.artwork = try await self.extractArtwork(from: file)
            }
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
    
    func playNext() {
        guard !allFiles.isEmpty else { return }
        var nextFile: URL
        repeat {
            nextFile = allFiles.randomElement()!
        } while nextFile == currentFile && allFiles.count > 1
        play(file: nextFile)
    }
    
    func audioPlayerDidFinishPlaying(_ player: AVAudioPlayer, successfully flag: Bool) {
        if flag {
            playNext() //Automatically play next song
        }
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
        let queue = DispatchQueue(label: "music.timer")
            queue.async {
                self.timer = Timer.scheduledTimer(withTimeInterval: 0.5, repeats: true) { _ in
                    DispatchQueue.main.async {
                        self.currentTime = self.player?.currentTime ?? 0
                    }
                }
                RunLoop.current.add(self.timer!, forMode: .common)
                RunLoop.current.run()
            }
    }
    private func stopTimer() {
        timer?.invalidate()
        timer = nil
    }
}

extension MusicModel{
    func extractArtwork(from url:URL) async throws ->NSImage? {
        let asset=AVURLAsset(url: url)
        
        let metadata=try await asset.loadMetadata(for: .init(rawValue: "org.xiph.vorbis-comment"))
        for item in metadata {

            if item.commonKey == .commonKeyArtwork {

                if let data = try? await item.load(.dataValue) {
                    if let image = NSImage(data: data) {
                        return image
                    }
                }
                
                if let dict = try? await item.load(.value) as? [String: Any],
                   let data = dict["data"] as? Data {
                    if let image = NSImage(data: data) {
                        return image
                    }
                }
            }
            
        }
        return nil
    }
}
