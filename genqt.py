# encoding: utf8

import sys
import os
import clang
import clang.cindex

# TODO 类中的enum
# 方法参数的默认值。
# signal/slots，这个有希望吗
# 全局qt函数
# 打包到一起编译，估计编译出来的库会非常大，按一个类1M，1000个类就是1G哇。
# 而pyqt5生成的binding层，一个qt模块最大3M，加起来30M的样子。

# 改进：
# 采用编译器的pass方式，从AST中不断推进

from gentool import GenTool


# 解析参数，调用不同的GenTool工具方法
def main():
    tool = GenTool()
    tool.walkgo()
    okcnt = 0
    for header in tool.genres:
        ok = tool.genres[header]
        if ok: okcnt += 1
    print('ok/total: %s/%s ' % (okcnt, len(tool.genres)))
    return


if __name__ == '__main__':
    main()
