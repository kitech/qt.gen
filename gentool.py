# encoding: utf8

import sys
import os
import clang
import clang.cindex

from genbase import GenerateBase, TestBuilder
from gengo import GenerateForGo, TestBuilderForGo
from geninline import GenerateForInlineCXX
from genrust import GenerateForRust


clang.cindex.Config.set_library_file('/usr/lib/libclang.so')

qtmodules = ['QtCore', 'QtGui', 'QtWidgets', 'QtNetwork', 'QtDBus']
qtmodules = ['QtGui', 'QtWidgets', 'QtNetwork', 'QtDBus']
# qtmodules = ['QtWidgets', 'QtNetwork', 'QtDBus']
compile_args = ['-x', 'c++', '-std=c++11', '-D__CODE_GENERATOR__']
compile_args += "-I/usr/include/qt -std=c++11 -DQT_CORE_LIB -DQT_NO_DEBUG -D_GNU_SOURCE -pipe -fno-exceptions -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -std=c++0x -Wall -W -D_REENTRANT -fPIC".split(' ')
for module in qtmodules: compile_args += ['-I/usr/include/qt/%s' % (module)]


class GenTool:
    def __init__(self):
        self.cursors = {}  # module => clang.cindex.Cursor
        self.generator = GenerateForGo()
        self.generator = GenerateForInlineCXX()
        self.generator = GenerateForRust()
        self.builder = TestBuilderForGo()
        self.genres = {}  # key => True | False
        self.conflib = clang.cindex.conf.lib
        self.cmake_code = ''
        self.cmake_code_module = ''
        return

    def cmake_header(self):
        code = ''
        code += "project(qtinline)\n"
        code += "cmake_minimum_required(VERSION 3.0)\n"
        code += "set(CMAKE_VERBOSE_MAKEFILE on)\n"
        code += "find_package(Qt5Core)\n"
        code += "find_package(Qt5Gui)\n"
        code += "find_package(Qt5Widgets)\n"
        code += "find_package(Qt5Network)\n"
        code += "set(CMAKE_CXX_FLAGS \"-O2 -std=c++11 -fno-exceptions\") # -std=c++14\")\n"
        code += "\n"
        return code

    # 单独测试每个qt类
    def walkgo(self):
        headers = self.gen_module_macro_include_path()
        for header in headers:
            module = header.split('/')[-1:].pop()
            self.walkgo_module(module, header)

            self.cmake_code += "add_library(%sInline SHARED\n" % (module)
            self.cmake_code += self.cmake_code_module
            self.cmake_code += ")\n"
            self.cmake_code += "qt5_use_modules(%sInline Core Gui Widgets Network)\n\n" % (module)
            self.cmake_code_module = ''
            break

        self.cmake_code = self.cmake_header() + self.cmake_code
        # self.generator.write_cmake_code(module, '', self.cmake_code)
        return

    def walkgo_module(self, module, header):
        print(module, header)
        fp = open(header, "r")
        while True:
            line = fp.readline()
            if line is None or len(line) == 0: break
            if not line.startswith('#include "'): continue
            class_file_name = line.split('"')[1]
            class_path_name = self.calc_class_path(module, class_file_name)
            class_name_lower = self.calc_class_name_lower(module, class_file_name)
            # print(line.strip(), class_file_name, class_path_name, class_name)
            if class_name_lower in ['qtcoreversion']: continue
            self.walkgo_class(module, class_path_name, class_name_lower)
        fp.close()
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
            methods = self.build_methods(cs)
            class_names.append([class_name, cs, methods])

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

    def dedup_classes(self, class_names):
        dedup_class_names = []
        unique_classes = {}
        for elems in class_names:
            class_name, cs, methods = elems
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

        if name.startswith('QAbstract'):
            # print('Omited abstract base class:' + name)
            return True

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

    def build_ast(self, module):
        if module in self.cursors:
            return self.cursors[module]

        global compile_args
        module_header = self.build_module_header(module)
        index = clang.cindex.Index.create()
        tu = index.parse(module_header, compile_args)
        cursor = tu.cursor
        self.cursors[module] = cursor
        return cursor

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

