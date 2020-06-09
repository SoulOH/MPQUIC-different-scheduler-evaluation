#!/usr/bin/python
# -*- coding: UTF-8 -*-
import os
import sys

def main():
    scheduler = sys.argv[1]
    if scheduler not in ["RR", "lowRTT", "utilRepair", "oppReduntant"]:
        return
    directory = "/home/souloh/testdata-x-50ms-50ms-5-loss-5-loss/"
    os.system("mkdir " + directory)
    smallFileData = directory + scheduler + "-small-file-test-data.log"
    largeFileData = directory + scheduler + "-large-file-test-data.log"

    with open(smallFileData, "a") as f:
        for i in range(0, 35):
            print("small file times " + str(i+1))
            smallFileName = directory + scheduler + "-smallFile-" + str(i) +".log"
            command1 = "go run example/client/main.go -m -o " \
            + smallFileName + " -s " + scheduler \
            + " https://10.0.2.2:6121/demo/small-file"
            os.system(command1)
            smallFileCompletionTime = get_download_time(smallFileName)
            f.write(smallFileCompletionTime)
    
    with open(largeFileData, "a") as f:
        for i in range(0, 35):
            print("large file times " + str(i+1))
            largeFileName = directory + scheduler + "-largeFile-" + str(i) +".log"
            command2 = "go run example/client/main.go -m -o " \
            + largeFileName + " -s " + scheduler \
            + " https://10.0.2.2:6121/demo/large-file"
            os.system(command2)
            largeFileCompletionTime = get_download_time(largeFileName)
            f.write(largeFileCompletionTime)

def get_download_time(fileName):
    with open(fileName, 'rb') as f:  # 打开文件
    # 在文本文件中，没有使用b模式选项打开的文件，只允许从文件头开始,只能seek(offset,0)
        offset = -50  # 设置偏移量
        while True:
            """
            file.seek(off, whence=0)：从文件中移动off个操作标记（文件指针），正往结束方向移动，负往开始方向移动。
            如果设定了whence参数，就以whence设定的起始位为准，0代表从头开始，1代表当前位置，2代表文件最末尾位置。 
            """
            f.seek(offset, 2)  # seek(offset, 2)表示文件指针：从文件末尾(2)开始向前50个字符(-50)
            lines = f.readlines()  # 读取文件指针范围内所有行
            if len(lines) >= 2:  # 判断是否最后至少有两行，这样保证了最后一行是完整的
                last_line = lines[-1]  # 取最后一行
                break
            # 如果off为50时得到的readlines只有一行内容，那么不能保证最后一行是完整的
            # 所以off翻倍重新运行，直到readlines不止一行
            offset *= 2
        return last_line.decode()

if __name__ == "__main__":
    main()