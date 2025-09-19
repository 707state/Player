//
//  MusicModel.swift
//  Player
//
//  Created by jask on 2025/9/18.
//
import AVFoundation
import Foundation

class MusicModel: ObservableObject{
    @Published var mp3files: [URL]=[]
    @Published var currentFile: URL?
    @Published var isPlaying = false
    @Published var currentTime: TimeInterval = 0
    @Published var duration: TimeInterval = 0
    private var player: AVAudioPlayer?
    private var timer: Timer?
    
    func loadMP3Files(url: URL){
        let fm=FileManager.default
        if let items=try? fm.contentsOfDirectory(  at: url,
                                                   includingPropertiesForKeys: nil,
                                                   options: [.skipsHiddenFiles]){
            mp3files = items.filter { $0.pathExtension.lowercased() == "mp3" || $0.pathExtension.lowercased()=="flac" }
        }
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
               print("播放失败:", error)
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
