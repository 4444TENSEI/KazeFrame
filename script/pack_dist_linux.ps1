$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o ./tmp/KazeFrame ./

# 输出的ZIP包路径和文件名
$ZIPFILE = "./KazeFrame_linux_dist.zip"

# 要压缩的文件和文件夹路径
$PathsToCompress = @(
    ".\tmp\KazeFrame"
    ".\static"
)

# 压缩
Compress-Archive -Path $PathsToCompress -DestinationPath $ZIPFILE -Force
# 输出执行成功提示信息
Write-Output "------ Build success(linux version), final output path: [$ZIPFILE] ------"