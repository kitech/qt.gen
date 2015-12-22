# encoding: utf8

import logging

import clang
import clang.cindex as clidx

FORMAT = '%(asctime)-15s %(filename)s:%(lineno)d %(funcName)s %(message)s'
LOGLEVEL= logging.DEBUG
# LOGLEVEL = logging.ERROR
logging.basicConfig(format=FORMAT, level=LOGLEVEL)
glog = logging.getLogger()


class GenUtil:
    def get_code_file(self, cursor):
        loc = cursor.location
        code_file = loc.file.name.split('/')[-1].split('.')[0]
        return code_file

    # like QtCore
    def get_decl_module(self, cursor):
        loc = cursor.location
        decl_module = loc.file.name.split('/')[-2]
        return decl_module

    # like core
    def get_decl_mod(self, cursor):
        loc = cursor.location
        decl_module = loc.file.name.split('/')[-2][2:].lower()

        # 自测
        if decl_module not in ['core', 'gui', 'widgets']:
            raise 'unknown module: %s, %s' % (decl_module, cursor.spelling)

        return decl_module

    # TODO 好像有点bug。
    # QListData会推导出基类是NotIndirectLayout。而实际上QListData没有基类。
    def get_base_class(self, cursor):
        bases = []
        for x in cursor.walk_preorder():
            # print(x.kind, x.spelling)
            if x.kind == clidx.CursorKind.CXX_BASE_SPECIFIER:
                decl = x.get_definition().type.get_declaration()
                # print(x.kind, decl.kind, decl.spelling)
                # fix, 需要decl.semantic_parent.kind == TRANSLATION_UNIT
                # 而这个遍历是有可能进入到类内部的，所以不准确
                if decl.semantic_parent.kind == clidx.CursorKind.TRANSLATION_UNIT:
                    bases.append(decl)
                else: break  # 提前跳出结束执行
        return bases

    def get_methods(self, class_cursor):
        method_names = {}

        for m in class_cursor.get_children():
            # print(m.kind, m.spelling)
            # TODO va_list type
            # if self.check_skip_method(m): continue
            mangled_name = m.mangled_name
            if m.kind == clang.cindex.CursorKind.CXX_METHOD:
                method_names[mangled_name] = m
            if m.kind == clang.cindex.CursorKind.CONSTRUCTOR:
                method_names[mangled_name] = m
            if m.kind == clang.cindex.CursorKind.DESTRUCTOR:
                method_names[mangled_name] = m

        return method_names

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
        self.newline = '\n'
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
        return self.newline.join(self.insert_points[name])

    def removePoint(self, name):
        codes = self.insert_points.pop(name)
        return self.newline.join(codes)

    # 按照names给出的顺序合并并导出代码。
    def exportCode(self, names):
        self.export_times += 1
        code = ''
        for name in names:
            code += self.newline.join(self.insert_points[name]) + self.newline
        return code

    def totalLength(self):
        tlen = 0
        for name in self.insert_points.keys():
            for line in self.insert_points[name]:
                tlen += len(line)
        return tlen

    def totalLine(self):
        tline = 0
        for name in self.insert_points.keys():
            tline += len(self.insert_points[name])
        return tline

    def reset(self):
        if self.export_times == 0:
            print('Warning, code maybe not export')
        self.insert_points = {}
        self.export_times = 0
        return
