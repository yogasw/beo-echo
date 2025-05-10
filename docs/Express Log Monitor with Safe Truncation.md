# 📄 Express Log Monitor with Safe Truncation

## 🧩 Tujuan

Aplikasi ini digunakan untuk:

* Memantau banyak file log (`logs/mock-*.log`)
* Menjaga agar isi setiap file log hanya maksimal **100 baris**
* Menyediakan **HTTP streaming API** agar client bisa memantau isi log secara realtime
* Menghindari konflik dengan proses penulisan (Mockoon) dengan pendekatan non-blocking

---

## 🗂 Struktur Direktori

```
project/
├── logs/                  # Folder tempat file log aktif (mock-4011.log, dll)
├── temp/                  # Folder sementara untuk copy log sebelum diproses
├── server.js              # File utama Express server
├── log-monitor.js         # Worker background untuk memantau dan mengelola logs
```

---

## 🔧 Fitur Aplikasi

### ✅ File Log Monitoring

* Setiap file `mock-*.log` dimonitor secara paralel
* Setiap 1 detik, isi file disalin ke file sementara
* Salinan tersebut digunakan untuk update buffer internal
* Jika tidak ada perubahan selama 5 detik, file asli ditulis ulang (trim ke 100 baris)

### ✅ HTTP Stream API

* Client dapat mengakses `/stream/:logId` untuk memantau log tertentu
* Data dikirim secara streaming (Transfer-Encoding: chunked)
* Update dikirim ke client setiap kali log diperbarui

---

## 💻 Implementasi

### `server.js`

```js
const express = require('express');
const fs = require('fs');
const path = require('path');
const { startLogMonitor, getLogBuffer, onLogUpdate } = require('./log-monitor');

const app = express();
const PORT = 3000;

startLogMonitor();

app.get('/stream/:logId', (req, res) => {
  const logId = req.params.logId;
  const fileName = `mock-${logId}.log`;
  const filePath = path.join('logs', fileName);

  if (!fs.existsSync(filePath)) {
    return res.status(404).send('Log file not found');
  }

  res.setHeader('Content-Type', 'text/plain');
  res.setHeader('Transfer-Encoding', 'chunked');

  const currentBuffer = getLogBuffer(fileName);
  res.write(currentBuffer.join('\n') + '\n');

  const updateHandler = (logName, lines) => {
    if (logName === fileName) {
      res.write(lines.join('\n') + '\n');
    }
  };

  onLogUpdate(updateHandler);

  req.on('close', () => {
    onLogUpdate(updateHandler, true);
  });
});

app.listen(PORT, () => {
  console.log(`Server listening at http://localhost:${PORT}`);
});
```

---

### `log-monitor.js`

```js
const fs = require('fs');
const path = require('path');
const readline = require('readline');

const LOG_DIR = 'logs';
const TEMP_DIR = 'temp';
const MAX_LINES = 100;
const CHECK_INTERVAL = 1000;
const IDLE_TIMEOUT = 5000;

const logBuffers = {};
const lastModified = {};
const listeners = new Set();

function onLogUpdate(callback, remove = false) {
  if (remove) {
    listeners.delete(callback);
  } else {
    listeners.add(callback);
  }
}

function emitLogUpdate(fileName, lines) {
  for (const cb of listeners) {
    cb(fileName, lines);
  }
}

function getLogBuffer(fileName) {
  return logBuffers[fileName] || [];
}

function listLogFiles() {
  return fs.readdirSync(LOG_DIR).filter(f => /^mock-\d+\.log$/.test(f));
}

async function readLines(filePath) {
  const lines = [];
  const rl = readline.createInterface({
    input: fs.createReadStream(filePath),
    crlfDelay: Infinity
  });

  for await (const line of rl) {
    lines.push(line);
  }

  return lines;
}

function monitorLog(fileName) {
  const fullPath = path.join(LOG_DIR, fileName);
  const tempPath = path.join(TEMP_DIR, `${fileName}.tmp`);

  if (!fs.existsSync(fullPath)) return;

  if (!logBuffers[fileName]) logBuffers[fileName] = [];

  setInterval(async () => {
    try {
      const stats = fs.statSync(fullPath);
      const now = Date.now();

      if (!lastModified[fileName]) {
        lastModified[fileName] = stats.mtimeMs;
      }

      fs.copyFileSync(fullPath, tempPath);

      const newLines = await readLines(tempPath);
      logBuffers[fileName] = (logBuffers[fileName] || []).concat(newLines).slice(-MAX_LINES);

      emitLogUpdate(fileName, newLines);

      if (stats.mtimeMs > lastModified[fileName]) {
        lastModified[fileName] = stats.mtimeMs;
      }

      const idle = now - lastModified[fileName] > IDLE_TIMEOUT;

      if (idle) {
        fs.writeFileSync(fullPath, logBuffers[fileName].join('\n') + '\n');
        console.log(`[Trimmed] ${fileName}`);
      }
    } catch (err) {
      console.error(`[Monitor Error] ${fileName}:`, err.message);
    }
  }, CHECK_INTERVAL);
}

function startLogMonitor() {
  if (!fs.existsSync(TEMP_DIR)) {
    fs.mkdirSync(TEMP_DIR);
  }

  const files = listLogFiles();
  files.forEach(monitorLog);
}

module.exports = {
  startLogMonitor,
  getLogBuffer,
  onLogUpdate,
};
```

---

## 🚀 Cara Menjalankan

1. Buat folder `logs/` dan `temp/`
2. Taruh file log seperti `mock-4011.log` dalam folder `logs/`
3. Jalankan:

   ```bash
   node server.js
   ```
4. Akses:

   ```bash
   curl http://localhost:3000/stream/4011
   ```

---

## 📌 Catatan Pengembangan

* Gunakan sistem polling agar tidak tergantung pada fitur `fs.watch` yang bisa tidak andal pada banyak file.
* Bisa ditambah `WebSocket` untuk pengiriman real-time jika ingin lebih efisien dari HTTP chunked.
* Untuk skala besar, tambahkan fitur re-scan direktori setiap 60 detik agar log baru otomatis dimonitor.

---

Silakan salin semua isi di atas dan simpan ke file `.md` atau dokumentasi kamu. Mau saya bantu juga buat skrip yang otomatis scan file baru setiap menit?
