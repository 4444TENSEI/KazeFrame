# 输出的ZIP包路径和文件名
$ZIPFILE = "KazeFrame_Source.zip"
 
# 要压缩的文件和文件夹路径
$PathsToCompress = @(
    ".air.toml"
    ".gitignore"
    "go.mod"
    "go.sum"
    "internal"
    "pkg"
    "script"
    "static"
    "main.go"
)
 
# 压缩
Compress-Archive -Path $PathsToCompress -DestinationPath $ZIPFILE -Force
# 输出执行成功提示信息
Write-Output "------The source code has been backed up successfully, final output path: [$ZIPFILE]------"