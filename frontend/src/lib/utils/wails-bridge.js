// wails-bridge.js - Bridge untuk komunikasi dengan Wails desktop app
// File ini akan di-embed oleh Wails saat runtime

// Check if running in Wails desktop environment
export const isDesktop = typeof window !== 'undefined' && window.wails !== undefined;

// Desktop-specific API endpoints - akan menggunakan Wails context API
export const desktopApi = {
  // Method untuk mengecek apakah berjalan di desktop
  isDesktopMode: () => isDesktop,

  // Method untuk mendapatkan versi aplikasi
  getAppVersion: () => {
    if (isDesktop && window.wails?.Environment) {
      return window.wails.Environment.appVersion || '1.0.0';
    }
    return '1.0.0';
  },

  // Method untuk mendapatkan informasi platform
  getPlatform: () => {
    if (isDesktop && window.wails?.Environment) {
      return window.wails.Environment.platform || 'unknown';
    }
    return 'web';
  },

  // Method untuk membuka URL eksternal di browser default
  openExternal: (url) => {
    if (isDesktop && window.wails?.BrowserOpenURL) {
      return window.wails.BrowserOpenURL(url);
    }
    // Fallback untuk web
    window.open(url, '_blank');
  },

  // Method untuk menampilkan dialog
  showDialog: (title, message, type = 'info') => {
    if (isDesktop && window.wails?.MessageDialog) {
      return window.wails.MessageDialog({
        type: type,
        title: title,
        message: message
      });
    }
    // Fallback untuk web
    return Promise.resolve(alert(`${title}: ${message}`));
  },

  // Method untuk quit aplikasi
  quit: () => {
    if (isDesktop && window.wails?.Quit) {
      return window.wails.Quit();
    }
    // Fallback untuk web
    window.close();
  }
};

// Export default untuk mudah digunakan
export default {
  isDesktop,
  ...desktopApi
};
