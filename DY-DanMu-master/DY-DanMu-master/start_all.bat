@echo off
cd /d %~dp0
chcp 65001 >nul
echo ==========================================
echo 🚀 正在启动斗鱼弹幕监控系统...
echo ==========================================

:: 杀死可能残留的进程
taskkill /F /IM persistServerRun.exe >nul 2>&1
taskkill /F /IM spiderRun.exe >nul 2>&1
taskkill /F /IM web.exe >nul 2>&1
taskkill /F /IM go.exe >nul 2>&1

:: 设置环境变量
set MYSQLPWD=1234
set REDISPWD=
set EMAILUSER=
set EMAILPWD=

echo [1/4] 正在启动 PersistServer (数据存储服务)...
echo 请耐心等待 15 秒，让它完成编译和启动...
:: 使用 cmd /c 以便在进程结束时自动关闭窗口
start "PersistServer" cmd /c "cd persistServer\run && go run persistServerRun.go"

:: 等待 15 秒确保端口就绪 (关键！)
timeout /t 15

echo [2/4] 正在启动 Spider (弹幕抓取服务)...
start "Spider" cmd /c "cd spider\run && go run spiderRun.go"

timeout /t 5

echo [3/4] 正在启动 Web Server (API 服务)...
start "WebServer" cmd /c "cd web\server && go run web.go"

timeout /t 5

echo [4/4] ✅ 所有服务已启动！正在打开监控页面...
start monitor.html

echo.
echo ==========================================
echo 🎉 系统运行中！
echo.
echo 若要停止服务，请在网页上点击“关闭系统”按钮，
echo 所有黑窗口将自动消失。
echo ==========================================
echo.
echo 按任意键退出此启动窗口...
pause >nul
