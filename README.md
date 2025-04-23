# 🔋 Battery Notifier (Go + Zenity)

A lightweight battery monitoring tool written in Go that alerts you via graphical popups using **Zenity** when your battery is either **low** or **fully charged**.
Perfect for Linux window managers like i3, XFCE, or full DEs like GNOME and KDE.

---

## 🚀 Features
- Alerts at **10%**, **20%**, **90%**, and **95–100%**
- Uses `upower` for battery info and `zenity` for clean GUI alerts
- Systemd-ready — runs automatically in the background

---

## ⚙️ How to Run

```bash
# 1. Build the binary
go build -o battery_alert battery_alert.go

# 2. Move it to your binary directory
mkdir -p ~/bin
mv battery_alert ~/bin/
chmod +x ~/bin/battery_alert

# 3. Set up as a systemd user service (recommended)
# See `systemd/` folder for service & timer files
