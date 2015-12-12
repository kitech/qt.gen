# encoding: utf8

import logging

import clang
import clang.cindex

FORMAT = '%(asctime)-15s %(filename)s:%(lineno)d %(funcName)s %(message)s'
LOGLEVEL= logging.DEBUG
# LOGLEVEL = logging.ERROR
logging.basicConfig(format=FORMAT, level=LOGLEVEL)
glog = logging.getLogger()


class GenUtil:
    pass


# 可以多点写入的代码编辑类
# 支持多点写入
# 支持前身写入
# 支持唯一写入
class CodePaper:
    def __init__(self):
        self.code = ''
        self.insert_points = {}  # name => [codes]
        self.export_times = 0
        return

    def addPoint(self, name):
        if name not in self.insert_points:
            self.insert_points[name] = []
        return

    def hasPoint(self, name):
        if name in self.insert_points: return True
        return False

    def allPoints(self):
        return self.insert_points.keys()

    def append(self, name, code):
        if not self.hasPoint(name): self.addPoint(name)
        self.insert_points[name].append(code)
        return

    def appendUnique(self, name, code):
        if not self.hasPoint(name): self.addPoint(name)
        if code not in self.insert_points[name]:
            self.insert_points[name].append(code)
        return

    def prepend(self, name, code):
        if not self.hasPoint(name): self.addPoint(name)
        self.insert_points[name].insert(0, code)
        return

    def AP(self, name, code): return self.append(name, code)

    def APU(self, name, code): return self.appendUnique(name, code)

    def PP(self, name, code): return self.prepend(name, code)

    def getPoint(self, name):
        return ''.join(self.insert_points[name])

    def removePoint(self, name):
        codes = self.insert_points.pop(name)
        return ''.join(codes)

    # 按照names给出的顺序合并并导出代码。
    def exportCode(self, names):
        self.export_times += 1
        code = ''
        for name in names:
            code += ''.join(self.insert_points[name])
        return code

    def totalLength(self):
        tlen = 0
        for name in self.insert_points.keys():
            for line in self.insert_points[name]:
                tlen += len(line)
        return tlen

    def reset(self):
        if self.export_times == 0:
            print('Warning, code maybe not export')
        self.insert_points = {}
        self.export_times = 0
        return
