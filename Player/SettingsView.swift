//
//  SettingsView.swift
//  Player
//
//  Created by jask on 2025/10/29.
//
import SwiftUI

struct SettingsView: View {
    @AppStorage("backend_url") private var backendUrl = ""
    @Environment(\.dismiss) private var dismiss
    var body: some View {
        NavigationStack {
            Form {
                Section(header: Text("Backend URL")) {
                    TextField("http://ip:port", text: $backendUrl)
                        .textFieldStyle(.roundedBorder)
                        .disableAutocorrection(true)
                }
            }
            .navigationTitle("Settings")
            .toolbar {
                ToolbarItem(placement: .automatic) {
                    Button("Apply") {
                        dismiss()  // 退出设置页
                    }
                }
            }
        }
    }

}
