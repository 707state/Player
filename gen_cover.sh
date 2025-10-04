#!/bin/bash
root="${1:-/Volumes/NO NAME/4040 - Good Thing Bad Thing Who Knows/}"
# 遍历所有子目录
find "$root" -type d | while read -r dir; do
    # 找第一个 mp3 或 flac
    firstfile=$(find "$dir" -maxdepth 1 -type f \( -iname "*.mp3" -o -iname "*.flac" \) | head -n 1)
    if [ -n "$firstfile" ]; then
        coverfile="$dir/cover.jpg"
        if [ ! -f "$coverfile" ]; then
            echo "Extracting cover for: $dir"
            exiftool -b -Picture "$firstfile" > "$coverfile"
        fi
    fi
done
