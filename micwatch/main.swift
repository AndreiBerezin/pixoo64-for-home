import Foundation
import CoreAudio
import AudioToolbox

// MARK: - Config

struct Config {
    let baseURL: URL
    let logPath: String

    static func fromArgs() -> Config {
        let args = CommandLine.arguments

        func value(after flag: String) -> String? {
            guard let i = args.firstIndex(of: flag), i + 1 < args.count else { return nil }
            return args[i + 1]
        }

        let logPath = value(after: "--log") ?? "/tmp/micwatch.out.log"

        if let base = value(after: "--base-url"), let url = URL(string: base) {
            return Config(baseURL: normalize(url), logPath: logPath)
        }

        if let host = value(after: "--host") {
            let s = (host.hasPrefix("http://") || host.hasPrefix("https://"))
                ? host
                : "http://\(host)"
            if let url = URL(string: s) {
                return Config(baseURL: normalize(url), logPath: logPath)
            }
        }

        return Config(baseURL: URL(string: "http://127.0.0.1")!, logPath: logPath)
    }

    private static func normalize(_ url: URL) -> URL {
        var s = url.absoluteString
        while s.hasSuffix("/") { s.removeLast() }
        return URL(string: s) ?? url
    }

    func url(path: String) -> URL {
        baseURL.appendingPathComponent(path.hasPrefix("/") ? String(path.dropFirst()) : path)
    }
}

// MARK: - Logger

final class Logger {
    private let iso = ISO8601DateFormatter()
    private let fileURL: URL
    private let lock = NSLock()

    init(path: String) {
        self.fileURL = URL(fileURLWithPath: path)
    }

    func log(_ msg: String) {
        let line = "[\(iso.string(from: Date()))] \(msg)\n"
        lock.lock()
        defer { lock.unlock() }

        do {
            let data = Data(line.utf8)

            if FileManager.default.fileExists(atPath: fileURL.path) {
                let handle = try FileHandle(forWritingTo: fileURL)
                _ = try handle.seekToEnd()
                try handle.write(contentsOf: data)
                try handle.close()
            } else {
                try data.write(to: fileURL, options: .atomic)
            }
        } catch {
        }
    }
}

// MARK: - MicWatch

final class MicWatch {
    private var defaultInputDeviceID: AudioDeviceID = 0
    private var lastState: Bool? = nil

    private let config: Config
    private let logger: Logger

    init(config: Config, logger: Logger) {
        self.config = config
        self.logger = logger
    }

    func start() {
        logger.log("MicWatch started. Base URL: \(config.baseURL.absoluteString)")

        updateDefaultInputDevice()
        installDefaultInputListener()
        installRunningListenerForCurrentDevice()

        reportIfChanged()
        RunLoop.current.run()
    }

    // MARK: - Listeners

    private func installDefaultInputListener() {
        var addr = AudioObjectPropertyAddress(
            mSelector: kAudioHardwarePropertyDefaultInputDevice,
            mScope: kAudioObjectPropertyScopeGlobal,
            mElement: kAudioObjectPropertyElementMain
        )

        AudioObjectAddPropertyListenerBlock(
            AudioObjectID(kAudioObjectSystemObject),
            &addr,
            DispatchQueue.global(qos: .utility)
        ) { [weak self] _, _ in
            guard let self else { return }
            let old = self.defaultInputDeviceID
            self.updateDefaultInputDevice()
            if self.defaultInputDeviceID != old {
                self.installRunningListenerForCurrentDevice()
                self.reportIfChanged()
            }
        }
    }

    private func installRunningListenerForCurrentDevice() {
        guard defaultInputDeviceID != 0 else { return }

        var addr = AudioObjectPropertyAddress(
            mSelector: kAudioDevicePropertyDeviceIsRunningSomewhere,
            mScope: kAudioObjectPropertyScopeGlobal,
            mElement: kAudioObjectPropertyElementMain
        )

        AudioObjectAddPropertyListenerBlock(
            AudioObjectID(defaultInputDeviceID),
            &addr,
            DispatchQueue.global(qos: .utility)
        ) { [weak self] _, _ in
            self?.reportIfChanged()
        }
    }

    // MARK: - CoreAudio

    private func updateDefaultInputDevice() {
        var addr = AudioObjectPropertyAddress(
            mSelector: kAudioHardwarePropertyDefaultInputDevice,
            mScope: kAudioObjectPropertyScopeGlobal,
            mElement: kAudioObjectPropertyElementMain
        )

        var deviceID = AudioDeviceID(0)
        var size = UInt32(MemoryLayout<AudioDeviceID>.size)

        let status = AudioObjectGetPropertyData(
            AudioObjectID(kAudioObjectSystemObject),
            &addr,
            0,
            nil,
            &size,
            &deviceID
        )

        defaultInputDeviceID = (status == noErr) ? deviceID : 0
    }

    private func isMicRunning() -> Bool {
        guard defaultInputDeviceID != 0 else { return false }

        var addr = AudioObjectPropertyAddress(
            mSelector: kAudioDevicePropertyDeviceIsRunningSomewhere,
            mScope: kAudioObjectPropertyScopeGlobal,
            mElement: kAudioObjectPropertyElementMain
        )

        var running: UInt32 = 0
        var size = UInt32(MemoryLayout<UInt32>.size)

        let status = AudioObjectGetPropertyData(
            defaultInputDeviceID,
            &addr,
            0,
            nil,
            &size,
            &running
        )

        return status == noErr && running != 0
    }

    // MARK: - State + HTTP

    private func reportIfChanged() {
        let running = isMicRunning()
        if lastState == running { return }
        lastState = running

        logger.log("MIC \(running ? "ON" : "OFF")")

        let url = config.url(path: running ? "/mic/on" : "/mic/off")
        sendPOST(url: url)
    }

    private func sendPOST(url: URL) {
        var req = URLRequest(url: url)
        req.httpMethod = "POST"
        req.timeoutInterval = 3

        URLSession.shared.dataTask(with: req) { [weak self] _, resp, err in
            if let err = err {
                self?.logger.log("HTTP error: \(err.localizedDescription). URL: \(url.absoluteString)")
                return
            }
            if let http = resp as? HTTPURLResponse,
               !(200...299).contains(http.statusCode) {
                self?.logger.log("HTTP status: \(http.statusCode). URL: \(url.absoluteString)")
            }
        }.resume()
    }
}

// MARK: - Main

let config = Config.fromArgs()
let logger = Logger(path: config.logPath)
let watcher = MicWatch(config: config, logger: logger)
watcher.start()