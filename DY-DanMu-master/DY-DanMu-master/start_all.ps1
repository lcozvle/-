# 设置环境变量
$env:MYSQLPWD="1234"
$env:REDISPWD=""
$env:EMAILUSER=""
$env:EMAILPWD=""

Write-Host "🚀 正在启动斗鱼弹幕监控系统..." -ForegroundColor Green

# 1. 启动 PersistServer
Write-Host "正在启动 PersistServer (数据存储服务)..."
Start-Process -FilePath "go" -ArgumentList "run", "persistServerRun.go" -WorkingDirectory "$PSScriptRoot\persistServer\run" -WindowStyle Minimized
Start-Sleep -Seconds 2

# 2. 启动 Spider
Write-Host "正在启动 Spider (弹幕抓取服务)..."
Start-Process -FilePath "go" -ArgumentList "run", "spiderRun.go" -WorkingDirectory "$PSScriptRoot\spider\run" -WindowStyle Minimized
Start-Sleep -Seconds 2

# 3. 启动 Web Server
Write-Host "正在启动 Web Server (API 服务)..."
Start-Process -FilePath "go" -ArgumentList "run", "web.go" -WorkingDirectory "$PSScriptRoot\web\server" -WindowStyle Minimized
Start-Sleep -Seconds 3

Write-Host "✅ 所有服务已在后台启动！" -ForegroundColor Green
Write-Host "正在打开监控页面..."

# 4. 打开监控页面
Start-Process "$PSScriptRoot\monitor.html"

Write-Host "按任意键退出此窗口（服务将继续在后台运行）..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
