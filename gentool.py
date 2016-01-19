# encoding: utf8

import sys
import os
import clang
import clang.cindex
import clang.cindex as clidx

from genbase import GenerateBase, TestBuilder
from gengo import GenerateForGo  # , TestBuilderForGo
from geninline import GenerateForInlineCXX
from genrust import GenerateForRust
from gencontext import *
from genutil import *


clang.cindex.Config.set_library_file('/usr/lib/libclang.so')

qtmodules = ['QtCore', 'QtGui', 'QtWidgets']
qtmodules.append('QtNetwork')
# qtmodules.append('QtDBus')

compile_args = ['-x', 'c++', '-std=c++11', '-D__CODE_GENERATOR__', '-D_GLIBCXX_USE_C++11ABI=0']
compile_args += "-I/usr/include/qt -std=c++11 -DQT_CORE_LIB -DQT_NO_DEBUG -D_GNU_SOURCE -pipe -fno-exceptions -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -std=c++0x -Wall -W -D_REENTRANT -fPIC".split(' ')
for module in qtmodules: compile_args += ['-I/usr/include/qt/%s' % (module)]


class GenTool:
    def __init__(self):
        self.cursors = {}  # module => clang.cindex.Cursor
        self.generator = GenerateForGo()
        self.generator = GenerateForInlineCXX()
        self.generator = GenerateForRust()
        # self.builder = TestBuilderForGo()
        self.genres = {}  # key => True | False
        self.conflib = clang.cindex.conf.lib
        self.gctx = GenContext()
        self.gutil = GenUtil()
        self.tuc = None
        # self.cmake_code = ''
        # self.cmake_code_module = ''
        return

    # def cmake_header(self):
    #     code = ''
    #     code += "project(qtinline)\n"
    #     code += "cmake_minimum_required(VERSION 3.0)\n"
    #     code += "set(CMAKE_VERBOSE_MAKEFILE on)\n"
    #     code += "find_package(Qt5Core)\n"
    #     code += "find_package(Qt5Gui)\n"
    #     code += "find_package(Qt5Widgets)\n"
    #     code += "find_package(Qt5Network)\n"
    #     code += "set(CMAKE_CXX_FLAGS \"-O2 -std=c++11 -fno-exceptions\") # -std=c++14\")\n"
    #     code += "\n"
    #     return code

    # 单独测试每个qt类
    def walkgo(self):
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
        self.generator.setGenContext(self.gctx)
        self.generator.genpass()
        self.gctx.dumpContext()

        for k in self.gctx.classes:
            if not k.startswith('Q'): print(k)

        return
        headers = self.gen_module_macro_include_path()
        for header in headers:
            module = header.split('/')[-1:].pop()
            # self.walkgo_module(module, header)

            # self.cmake_code += "add_library(%sInline SHARED\n" % (module)
            # self.cmake_code += self.cmake_code_module
            # self.cmake_code += ")\n"
            # self.cmake_code += "qt5_use_modules(%sInline Core Gui Widgets Network)\n\n" % (module)
            # self.cmake_code_module = ''
            break

        # self.cmake_code = self.cmake_header() + self.cmake_code
        # self.generator.write_cmake_code(module, '', self.cmake_code)
        return

    # def walkgo_module(self, module, header):
    #     print(module, header)
    #     fp = open(header, "r")
    #     while True:
    #         line = fp.readline()
    #         if line is None or len(line) == 0: break
    #         if not line.startswith('#include "'): continue
    #         class_file_name = line.split('"')[1]
    #         class_path_name = self.calc_class_path(module, class_file_name)
    #         class_name_lower = self.calc_class_name_lower(module, class_file_name)
    #         # print(line.strip(), class_file_name, class_path_name, class_name)
    #         if class_name_lower in ['qtcoreversion']: continue
    #         self.walkgo_class(module, class_path_name, class_name_lower)
    #     fp.close()
    #     return

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

    # 构造这个类的AST，找到类的定义
    def walkgo_class(self, module, header, class_name_lower):
        cursor = self.build_ast(module)

        class_names = []
        for cs in cursor.get_children():
            class_name = cs.spelling
            if cs.kind != clang.cindex.CursorKind.CLASS_DECL: continue
            # if cs.kind != clang.cindex.CursorKind.CLASS_TEMPLATE: continue
            if class_name == 'QIODevice':
                # print(header, class_name_lower, name)
                if class_name_lower == class_name.lower():
                    skip = self.check_skip_class(header, class_name_lower, cs)
                    # print(header, class_name_lower, class_name)
                    # print(555, skip)
                    # exit(0)

            if self.check_skip_class(header, class_name_lower, cs): continue
            if not self.check_inheader_class(header, class_name_lower, cs): continue
            # print(header, class_name_lower, class_name, 'matched')
            base_classes = self.get_base_class(cs)
            base_class = base_classes[0] if len(base_classes) > 0 else None
            if base_class is not None:
                print(class_name, '->', base_class.spelling)
            methods = self.build_methods(cs)
            class_names.append([class_name, cs, methods, base_class])

        if len(class_names) == 0:
            print('wtf: no class found:', header)
            return

        if len(class_names) > 1:
            print(header, ' has more than one class:' + header)
            pass

        class_names = self.dedup_classes(class_names)
        # self.cmake_code_module += self.generator.generateCMake(module, class_names)
        class_code = self.generator.generateClasses(module, class_names)
        # self.generator.write_swig_code(module, class_code)

        # try build it
        # ret = self.builder.tryBuild()
        # if ret == 0: print('Builder class ok: ' + header)
        # self.genres[header] = True if ret == 0 else False
        # if ret != 0:
        #    print('Builder class faild: ' + header)

        #    exit(0)
        return

    def get_base_class(self, class_cursor):
        semp = class_cursor.semantic_parent
        lexp = class_cursor.lexical_parent
        # print(semp.spelling, lexp.spelling)
        # print(class_cursor.type.spelling)

        def skip_base_class(cursor):
            name = cursor.spelling
            # if name.startswith('QAbstract'): return True
            if name.endswith('Interface'): return True
            if name.endswith('Private'): return True
            return False

        bases = []
        for x in class_cursor.walk_preorder():
            # print(x.kind, x.spelling)
            if x.kind == clidx.CursorKind.CXX_BASE_SPECIFIER:
                decl = x.get_definition().type.get_declaration()
                if skip_base_class(decl): break  # 提前终止，提高查找速度
                # print(x.kind, decl.kind, decl.spelling)
                bases.append(decl)
        # print(bases, len(bases))
        return bases

    def dedup_classes(self, class_names):
        dedup_class_names = []
        unique_classes = {}
        for elems in class_names:
            class_name, cs, methods, base_class = elems
            if class_name in unique_classes: continue
            unique_classes[class_name] = True
            dedup_class_names.append(elems)
        return dedup_class_names

    def build_methods(self, class_cursor):
        method_names = {}

        for m in class_cursor.get_children():
            # print(m.kind, m.spelling)
            # TODO va_list type
            if self.check_skip_method(m): continue
            mangled_name = m.mangled_name
            if m.kind == clang.cindex.CursorKind.CXX_METHOD:
                method_names[mangled_name] = m
            if m.kind == clang.cindex.CursorKind.CONSTRUCTOR:
                method_names[mangled_name] = m
            if m.kind == clang.cindex.CursorKind.DESTRUCTOR:
                method_names[mangled_name] = m

        return method_names

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

    def check_skip_method(self, method_cursor):
        if method_cursor.spelling in ['sprintf', 'vsprintf', 'objectNameChanged']:
            return True
        return False

    def check_skip_class(self, header, class_name_lower, class_cursor):
        cs = class_cursor
        name = cs.spelling

        if name.endswith('Private'):
            # print('Omited private internal class:' + name)
            return True

        if name.endswith('Ref'):
            # print('Omited private internal class:' + name)
            return True

        # if name.startswith('QAbstract'):
            # print('Omited abstract base class:' + name)
        #    return True

        # 如果有虚拟方法，无法实例化，则不生成文类的封装类
        # for subc in class_cursor.get_children():
        #    if self.conflib.clang_CXXMethod_isPureVirtual(subc): return True

        if name in ['QTextStreamManipulator', 'QMetaType', 'QMetaProperty',
                    'QFutureInterfaceBase', 'QFuture', 'QFutureInterface', 'QFutureWatcher',
                    'QCommandLineParser',
                    'QStringBuilder', 'QTypeInfo', 'QNoDebug',
                    'QInternal', 'QVariantAnimation', 'Connection']:
            return True

        # for test
        if name.startswith('QJson'): return True

        if not cs.is_definition():
            # print('Omited non definition class: ' + name)
            return True

        # debuging
        # if not name.startswith('QDir'): return True
        # if name != 'QString': return True
        # if name not in ["QString", "QObject", "QThread", "QUrl"]: return True

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

