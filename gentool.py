# encoding: utf8

import sys
import os
import clang
import clang.cindex
import clang.cindex as clidx

from genbase import GenerateBase, TestBuilder
from gengo import GenerateForGo  # , TestBuilderForGo
from geninc import GenerateForInc
from genrust import GenerateForRust
from gencontext import *
from genutil import *


clang.cindex.Config.set_library_file('/usr/lib/libclang.so')

qtmodules = ['QtCore', 'QtGui', 'QtWidgets']
# qtmodules.append('QtNetwork')
# qtmodules.append('QtQml')
# qtmodules.append('QtQuick')
# qtmodules.append('QtDBus')

compile_args = ['-x', 'c++', '-std=c++11', '-D__CODE_GENERATOR__', '-D_GLIBCXX_USE_C++11ABI=1']
compile_args += "-I/usr/include/qt -std=c++11 -DQT_CORE_LIB -DQT_NO_DEBUG -D_GNU_SOURCE -pipe -fno-exceptions -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -std=c++0x -Wall -W -D_REENTRANT -fPIC".split(' ')
for module in qtmodules: compile_args += ['-DQT_%s_LIB' % (module[2:].lower())]
for module in qtmodules: compile_args += ['-I/usr/include/qt/%s' % (module)]


class GenTool:
    def __init__(self):
        self.cursors = {}  # module => clang.cindex.Cursor
        self.generator = None  # GenerateBase
        # self.builder = TestBuilderForGo()
        self.genres = {}  # key => True | False
        self.conflib = clang.cindex.conf.lib
        self.gctx = GenContext()
        self.gutil = GenUtil()
        self.tuc = None
        return

    # 单独测试每个qt类
    def walkinc(self):
        self.generator = GenerateForInc()
        self.walkCommon()
        return

    def walkgo(self):
        self.generator = GenerateForGo()
        self.walkCommon()
        return

    def walkrust(self):
        self.generator = GenerateForRust()
        self.walkCommon()
        return

    def walkCommon(self):
        self.tuc = cursor = self.build_ast()
        print(self.tuc.kind)

        idx = 0
        for sub in cursor.get_children(): idx += 1
        print('unit count:', idx)

        for c in cursor.get_children():
            idx += 1
            # nmodule = self.gutil.get_decl_module(c)
            # if nmodule != module: continue
            # if 'QMetaObject' in c.spelling: print(c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
            if c.kind == clidx.CursorKind.CLASS_TEMPLATE and c.is_definition():
                # print(c.kind, c.spelling, c.displayname, c.location)
                self.gctx.addClass(c)
            elif c.kind == clidx.CursorKind.CLASS_DECL and c.is_definition():
                # print(c, c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
                self.gctx.addClass(c)
            elif c.kind == clidx.CursorKind.STRUCT_DECL and c.is_definition():
                # print(c, c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
                if c.spelling.startswith('Q'):
                    self.gctx.addClass(c)
                pass
            elif c.kind == clidx.CursorKind.FUNCTION_TEMPLATE:
                self.gctx.addFunction(c)
            elif c.kind == clidx.CursorKind.FUNCTION_DECL:
                # print(c, c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
                if not c.spelling.startswith('_'): self.gctx.addFunction(c)
            elif c.kind == clidx.CursorKind.ENUM_DECL and c.is_definition():
                # print(c, c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
                pass
            elif c.is_definition():
                # print(c, c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
                pass

        self.gctx.dumpContext()
        print('Running...', self.generator.__class__.__name__)
        self.generator.setGenContext(self.gctx)
        self.generator.genpass()
        self.gctx.dumpContext()

        for k in self.gctx.classes:
            if not k.startswith('Q'): print(k)

        return

    def walkgo_module(self, module, header):
        cursor = self.build_ast(module)
        idx = 0
        for sub in cursor.get_children(): idx += 1
        print('unit count:', idx)

        for c in cursor.get_children():
            idx += 1
            # nmodule = self.gutil.get_decl_module(c)
            # if nmodule != module: continue

            if c.kind == clidx.CursorKind.CLASS_TEMPLATE and c.is_definition():
                print(c.kind, c.spelling, c.displayname, c.location)
            if c.kind == clidx.CursorKind.CLASS_DECL and c.is_definition():
                # print(c, c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
                self.gctx.addClass(c)
            if c.kind == clidx.CursorKind.FUNCTION_DECL and c.is_definition():
                # print(c, c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
                self.gctx.addFunction(c)
            if c.kind == clidx.CursorKind.ENUM_DECL and c.is_definition():
                # print(c, c.kind, c.spelling, c.displayname, c.is_definition(), c.location)
                pass

        self.gctx.dumpContext()
        self.generator.setGenContext(self.gctx)
        self.generator.genpass(module)
        self.gctx.dumpContext()
        return

    def dump_method(self, method_cursor):
        m = method_cursor
        print(m.location, m.result_type,
              self.generator.resolve_swig_type_name(m.result_type),
              self.generator.resolve_swig_type_name(m.type),
              m.mangled_name, m.get_tokens())
        for token in m.get_tokens():
            print(token, token.kind, token.spelling)
        return

    # 判断类是否属于这个头文件
    def check_inheader_class(self, header, class_name_lower, class_cursor):
        cs = class_cursor
        name = cs.spelling

        # print(cs.location.file.name, header)
        if cs.location.file.name == header:
            return True

        return False

    def build_ast(self):
        # if module in self.cursors:
        #    return self.cursors[module]
        import os

        hdrsrc = './qthdrsrc.h'
        astfile = './qthdrsrc.ast'

        index = clang.cindex.Index.create()
        if os.path.exists(astfile):
            tu = index.read(astfile)
        else:
            global compile_args
            tu = index.parse(hdrsrc, compile_args)
            tu.save(astfile)

        cursor = tu.cursor
        # print(cursor.kind)
        return cursor

    def load_ast(self):
        return

    def calc_class_path(self, module, header):
        path = '/usr/include/qt/%s/%s' % (module, header)
        return path

    def calc_class_name_lower(self, module, header):
        class_name = header.split('.')[0]
        return class_name

    def build_module_header(self, module):
        return '/usr/include/qt/%s/%s' % (module, module)

    def gen_module_macro_include_path(self):
        global qtmodules
        prefix = '/usr/include/qt'
        paths = []
        for m in qtmodules:
            paths.append(prefix + '/' + m + '/' + m)
        return paths

