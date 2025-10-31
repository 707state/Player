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
    @Published var single: String = ""
    @Published var album: String = ""
    @Published var artist: String = ""
    @Published var isLiked: Bool = false
    
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
            let data = try Data(contentsOf: file, options: .mappedIfSafe)
            
            // Step 2: 使用内存数据初始化 AVAudioPlayer
            player = try AVAudioPlayer(data: data)
            player?.delegate = self
            player?.prepareToPlay()
            player?.play()
            duration = player?.duration ?? 0
            currentFile = file
            isPlaying = true
            // Extract metadata (artwork + basic tags)
            Task.detached { @MainActor in
                let meta = try? await self.extractBasicMetadata(from: file)
                if let meta = meta {
                    self.single = meta.title
                    self.artist = meta.artist
                    self.album = meta.album
                    self.artwork = meta.artwork
                }
                await self.checkLikeStatus()
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
    // MARK: - Metadata
    struct BasicMeta { let title: String; let artist: String; let album: String; let artwork: NSImage? }
    func extractBasicMetadata(from url: URL) async throws -> BasicMeta {
        let asset = AVURLAsset(url: url)
        var title = url.deletingPathExtension().lastPathComponent
        var artist = ""
        var album = ""
        var artwork: NSImage?=nil
        for format in try await asset.load(.availableMetadataFormats) {
            let metadata = try await asset.loadMetadata(for: format)
            for item in metadata {
                if item.commonKey == .commonKeyTitle, let t = try? await item.load(.stringValue), !t.isEmpty { title = t }
                if item.commonKey == .commonKeyArtist, let a = try? await item.load(.stringValue), !a.isEmpty { artist = a }
                if item.commonKey == .commonKeyAlbumName, let al = try? await item.load(.stringValue), !al.isEmpty { album = al }
                if item.commonKey == .commonKeyArtwork{
                    if let data = try? await item.load(.dataValue) {
                        if let image = NSImage(data: data) {
                            artwork = image
                        }
                    }
                    
                    if let dict = try? await item.load(.value) as? [String: Any],
                       let data = dict["data"] as? Data {
                        if let image = NSImage(data: data) {
                            artwork = image
                        }
                    }
                }
            }
        }
        return BasicMeta(title: title, artist: artist, album: album,artwork: artwork)
    }
}

// MARK: - Backend integration
extension MusicModel {
    private func logNetworkError(context: String, url: URL?, response: URLResponse?, data: Data?, error: Error?) {
        print("[Network] Context=\(context)")
        if let url { print("[Network] URL=\(url.absoluteString)") }
        if let http = response as? HTTPURLResponse {
            print("[Network] Status=\(http.statusCode)")
        }
        if let error { print("[Network] Error=\(error.localizedDescription)") }
        if let data, let text = String(data: data, encoding: .utf8) {
            print("[Network] Body=\(text)")
        }
    }
    private func backendBaseURL() -> URL? {
        let str = UserDefaults.standard.string(forKey: "backend_url") ?? ""
        guard !str.trimmingCharacters(in: .whitespaces).isEmpty, let url = URL(string: str) else {
            return nil
        }
        return url
    }

    private func buildArtistsArray(from artistString: String) -> [String] {
        // Split by common delimiters
        let delimiters: CharacterSet = [",", "/", "&", "、"]
        let parts = artistString.components(separatedBy: delimiters)
            .map { $0.trimmingCharacters(in: .whitespacesAndNewlines) }
            .filter { !$0.isEmpty }
        return parts.isEmpty ? [artistString].filter { !$0.isEmpty } : parts
    }

    @MainActor
    func checkLikeStatus() async {
        guard let base = backendBaseURL(), !album.isEmpty else {
            self.isLiked = false
            return
        }
        // Query by album + one artist
        let oneArtist = buildArtistsArray(from: artist).joined(separator:",")
        var comps = URLComponents(url: base.appendingPathComponent("/single"), resolvingAgainstBaseURL: false)!
        comps.queryItems = [
            URLQueryItem(name: "album", value: album),
            URLQueryItem(name: "artists", value: oneArtist),
            URLQueryItem(name:"title",value: single),
        ]
        guard let url = comps.url else { self.isLiked = false; return }
        do {
            let (data, response) = try await URLSession.shared.data(from: url)
            if let http = response as? HTTPURLResponse, http.statusCode == 200 {
                if let obj = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
                   let exists = obj["exists"] as? Bool {
                    self.isLiked = exists
                } else {
                    self.isLiked = false
                }
            } else {
                logNetworkError(context: "checkLikeStatus", url: url, response: response, data: data, error: nil)
                self.isLiked = false
            }
        } catch {
            logNetworkError(context: "checkLikeStatus", url: url, response: nil, data: nil, error: error)
            self.isLiked = false
        }
    }

    func toggleLike() {
        Task {
            guard let base = backendBaseURL() else { return }
            let artistsArray = buildArtistsArray(from: artist)
            let payload: [String: Any] = [
                "title": single,
                "artists": artistsArray,
                "album": album,
            ]
            do {
                var request = URLRequest(url: base.appendingPathComponent("/single"))
                request.httpMethod = isLiked ? "DELETE" : "POST"
                request.addValue("application/json", forHTTPHeaderField: "Content-Type")
                request.httpBody = try JSONSerialization.data(withJSONObject: payload)
                let (data, response) = try await URLSession.shared.data(for: request)
                if let http = response as? HTTPURLResponse, (200...299).contains(http.statusCode) {
                    await MainActor.run { self.isLiked.toggle() }
                } else {
                    logNetworkError(context: "toggleLike", url: request.url, response: response, data: data, error: nil)
                }
            } catch {
                logNetworkError(context: "toggleLike", url: nil, response: nil, data: nil, error: error)
            }
        }
    }
}
