$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o ./tmp/KazeFrame.exe ./

# 输出的ZIP包路径和文件名
$ZIPFILE = "./KazeFrame_win_dist.zip"

# 要压缩的文件和文件夹路径
$PathsToCompress = @(
    ".\tmp\KazeFrame.exe"
    ".\static"
)
# 压缩
Compress-Archive -Path $PathsToCompress -DestinationPath $ZIPFILE -Force
# 输出执行成功提示信息
Write-Output "------ Build success(windows version), final output path: [$ZIPFILE] ------"
