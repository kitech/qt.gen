# encoding: utf8

from genbase import GenerateBase, TestBuilder


class GenerateForGo(GenerateBase):

    def generateHeader(self, module):
        tcode = ''
        tcode += "%module qstring\n"
        tcode += "%{\n"
        tcode += "// #include \"QtCore\"\n"
        tcode += "#include <%s>\n" % module
        tcode += "%}\n\n"
        return tcode

    def generateClasses(self, module, class_decls):
        code = ''

        for elems in class_decls:
            class_name, cs, methods = elems
            code += self.generateClass(class_name, cs, methods)

        return code

    def generateClass(self, class_name, cs, methods):

        code = ''
        code += 'class %s{\n' % class_name
        code += 'public:\n'
        code += "//    %s();\n" % class_name
        code += "//    ~%s();\n" % class_name

        for mth in methods:
            cursor = methods[mth]
            if self.check_skip_method(cursor): continue

            method_line1 = self.build_swig_method(class_name, mth, cursor)
            method_line2 = self.restore_method_by_token(cursor)
            method_line = method_line1
            # TODO hotfix
            # print(method_line1, 567)
            # print(method_line2, 234)
            if mth in ['fileInfo']:
                method_line = method_line2
                # exit(0)
            code += method_line

        code += "};\n\n"
        return code

    def build_swig_method(self, class_name, method_name, method_cursor):
        return self.generateMethod(class_name, method_name, method_cursor)

    def generateMethod(self, class_name, method_name, method_cursor):
        cursor = method_cursor

        code = ''
        return_type = cursor.result_type
        return_real_type = self.real_type_name(return_type)
        if '::' in return_real_type: return code
        if self.check_skip_params(cursor): return code

        fixmthname = self.fix_conflict_method_name(method_name)

        if cursor.access_specifier == clang.cindex.AccessSpecifier.PUBLIC:
            code += "  public: "
        elif cursor.access_specifier == clang.cindex.AccessSpecifier.PROTECTED:
            code += "  protected: "
        elif cursor.access_specifier == clang.cindex.AccessSpecifier.PRIVATE:
            code += "  private: "

        if cursor.is_static_method(): code += " static"

        if cursor.kind == clang.cindex.CursorKind.CONSTRUCTOR or \
           cursor.kind == clang.cindex.CursorKind.DESTRUCTOR:
            code += " %s(" % (fixmthname)
        else:
            return_type_name = self.resolve_swig_type_name(class_name, return_type)
            return_type_name2 = self.hotfix_typename_ifenum_asint(class_name, method_cursor, return_type)
            return_type_name = return_type_name2 if return_type_name2 is not None else return_type_name
            code += " %s %s(" % (return_type_name, fixmthname)

        argv = self.generateParams(class_name, method_name, method_cursor)
        if len(argv) > 0: code += ', '.join(argv)

        code += ");\n"

        # if skip_method is True: return ''
        return code

    # @return []
    def generateParams(self, class_name, method_name, method_cursor):
        idx = 0
        argv = []

        for arg in method_cursor.get_arguments():
            idx += 1
            # print('%s, %s, ty:%s, kindty:%s' % (method_name, arg.displayname, arg.type.spelling, arg.kind))
            # print('arg type kind: %s, %s' % (arg.type.kind, arg.type.get_declaration()))
            # param_line2 = self.restore_param_by_token(arg)
            # print(param_line2)

            type_name = self.resolve_swig_type_name(class_name, arg.type)
            type_name2 = self.hotfix_typename_ifenum_asint(class_name, arg, arg.type)
            type_name = type_name2 if type_name2 is not None else type_name

            arg_name = 'arg%s' % idx if arg.displayname == '' else arg.displayname
            # try fix void (*)(void *) 函数指针
            # 实际上swig不需要给定名字，只需要类型即可。
            if arg.type.kind == clang.cindex.TypeKind.POINTER and "(*)" in type_name:
                argelem = "%s" % (type_name.replace('(*)', '(*%s)' % arg_name))
            else:
                argelem = "%s %s" % (type_name, arg_name)
            argv.append(argelem)

        return argv

    # def hotfix_typename_ifenum_asint(self, class_name, arg):
    def hotfix_typename_ifenum_asint(self, class_name, token_cursor, atype):
        type_name = self.resolve_swig_type_name(class_name, atype)
        # if type_name not in ('int', 'int *', 'const int &'): return None
        type_name_segs = type_name.split(' ') 
        if 'int' not in type_name_segs: return None

        tokens = []
        for token in token_cursor.get_tokens():
            tokens.append(token.spelling)
            tkcursor = token.cursor

        # 为什么tokens是空呢，是不能识别的？
        if len(tokens) == 0: return None
        # TODO 全部使用replace方式，而不是这种每个符号的处理
        while tokens[0] in ['const', 'inline']:
            tokens = tokens[1:]

        firstch = tokens[0][0:1]
        if firstch.upper() == firstch and firstch != 'Q':
            print('Warning fix enum-as-int:', type_name, '=> %s::' % class_name, tokens[0])
            return '%s::%s' % (class_name, tokens[0])

        if len(tokens) < 3: return None
        if firstch.upper() == firstch and firstch == 'Q' and tokens[1] == '::':
            print('Warning fix enum-as-int2:', type_name, '=> %s::' % class_name, tokens[2])
            return '%s::%s' % (tokens[0], tokens[2])

        # like QtMsgType
        if firstch.upper() == firstch and firstch == 'Q' and tokens[0][0:2] == 'Qt':
            print('Warning fix enum-as-int3:', type_name, '=> ', tokens[0])
            return '%s' % (tokens[0])

        if firstch.upper() == firstch and firstch == 'Q' and tokens[0][1:1].lower() == tokens[0][1:1]:
            print('Warning fix enum-as-int4:', type_name, '=> ', type_name.replace('int', tokens[0]))
            return '%s' % (type_name.replace('int', tokens[0]))

        # like qint64...
        if firstch.lower() == firstch and tokens[0][0:1] == 'q' and '*' in type_name:
            print('Warning fix qint*-as-int5:', type_name, '=> ', tokens[0])
            return '%s %s' % (tokens[0], tokens[1])

        return None

    def restore_method_by_token(self, method_cursor):
        m = method_cursor
        method_line = ''
        for token in m.get_tokens():
            # print(token, token.kind, token.spelling)
            method_line += token.spelling + ' '

        return method_line

    def restore_param_by_token(self, arg_cursor):
        param_line = ''
        for token in arg_cursor.get_tokens():
            # print(token, token.kind, token.spelling)
            pass
        return param_line

    def fix_conflict_method_name(self, method_name):
        mthname = method_name
        fixmthname = mthname
        if mthname in ['type']:  # , 'select']:
            fixmthname = mthname + '_'
        return fixmthname

    def real_type_name(self, atype):
        type_name = atype.spelling

        if atype.kind == clang.cindex.TypeKind.TYPEDEF:
            # print('underlying type: %s' % atype.get_declaration().underlying_typedef_type.spelling)
            # print('underlying type: %s' % arg.type.underlying_typedef_type.spelling)
            type_name = atype.get_declaration().underlying_typedef_type.spelling
            if type_name.startswith('QFlags<'):
                type_name = type_name[7:len(type_name) - 1]

        return type_name

    # @return str
    def resolve_swig_type_name(self, class_name, atype):
        type_name = atype.spelling
        if type_name in ['QFunctionPointer', 'CategoryFilter',
                         'EasingFunction']:
            type_bclass = atype.get_declaration().semantic_parent
            # if type_name.startswith('Q'):
            # 全局定义的，不需要前缀
            if type_bclass.kind == clang.cindex.CursorKind.TRANSLATION_UNIT: pass
            else: type_name = '%s::%s' % (type_bclass.spelling, type_name)
        else:
            type_name = self.real_type_name(atype)

            # QTextStreamManipulator(void (QTextStream::*)(int) m, int a);
            # int registerNormalizedType(const ::QByteArray & normalizedTypeName, void * destructor, void *(*)(void *, const void *) constructor, int size, QMetaType::TypeFlags flags, const QMetaObject * metaObject);
            # qreal (*)(qreal) customType();
            # if type_name == 'void (*)(void *)':
            #    type_name = "void *"

        return type_name

    # @return True | False
    def check_skip_params(self, cursor):
        for arg in cursor.get_arguments():
            type_name = arg.type.spelling
            if 'QPrivate' in type_name: return True
            if 'QLatin1String' == type_name: return True
            # void directoryChanged(const QString & path, QFileSystemWatcher::QPrivateSignal arg0);
            if arg.displayname == '' and type_name == 'int': return True

        return False

    # @return True | False
    def check_skip_method(self, cursor):
        method_name = cursor.spelling
        if method_name.startswith('operator'):
            # print("Omited operator method: " + mth)
            return True

        # print('pub:' + str(cursor.access_specifier))
        if cursor.access_specifier == clang.cindex.AccessSpecifier.PUBLIC:
            pass
        if cursor.access_specifier == clang.cindex.AccessSpecifier.PROTECTED:
            if cursor.kind == clang.cindex.CursorKind.CONSTRUCTOR or \
               cursor.kind == clang.cindex.CursorKind.DESTRUCTOR:
                pass
            else: return True
        if cursor.access_specifier == clang.cindex.AccessSpecifier.PRIVATE:
            if cursor.kind == clang.cindex.CursorKind.CONSTRUCTOR or \
               cursor.kind == clang.cindex.CursorKind.DESTRUCTOR:
                pass
            else: return True

        istatic = cursor.is_static_method()
        # if istatic is True: return True

        # fix method
        fixmths = ['tr', 'trUtf8', 'qt_metacall', 'qt_metacast', 'data_ptr',
                   'sprintf', 'vsprintf', 'vasprintf', 'asprintf',
                   'entryInfoListcc',]
        if method_name in fixmths: return True
        fixmths_prefix = ['qt_check_for_']
        for p in fixmths_prefix:
            if method_name.startswith(p): return True

        # 实现不知道怎么fix了，已经fix，原来是给clang.cindex.parse中的-I不全，导致找不到类型。
        # fixmths3 = ['setQueryItems']
        # if method_name in fixmths3: return True

        return False

    def write_swig_code(self, module, code):
        tcode = self.generateHeader(module)

        code = tcode + code

        fpath = "core.swigcxx"
        f = os.open(fpath, os.O_CREAT | os.O_TRUNC | os.O_RDWR)
        os.write(f, code)
        os.close(f)
        return


import os
import subprocess
import shlex
class TestBuilderForGo(TestBuilder):
    def tryBuild(self):
        ret = self.cleanup_files()
        ret = self.run_swig()
        if ret == 0:
            ret = self.run_go_build()
        return ret

    def run_go_build(self):
        env = self.make_env()
        cmd = "/usr/bin/go install -v -x qt/core"
        cmd = shlex.split(cmd)
        proc = subprocess.Popen(cmd, env=env, cwd=os.getenv('PWD'), shell=False,
                                stderr=subprocess.STDOUT)
        outdata, errdata = proc.communicate()
        ret = proc.poll()
        # print(ret, proc.returncode, outdata, errdata)
        return ret

    def run_swig(self):
        cmd = self.make_cmd()
        env = self.make_env()
        # print(' '.join(cmd))
        # cmd = ["ls", "-l"]
        proc = subprocess.Popen(cmd, env=env, cwd=os.getenv('PWD'), shell=False,
                                stderr=subprocess.STDOUT)
        outdata, errdata = proc.communicate()
        ret = proc.poll()
        # print(ret, proc.returncode, outdata, errdata)
        # if ret == 0: print('ok')
        return ret

    def cleanup_files(self):
        workdir = os.getenv('PWD')
        cmd = 'rm -vf %s/src/qt/core/core*' % workdir
        cmd = shlex.split(cmd)
        proc = subprocess.Popen(cmd, env=None, cwd=os.getenv('PWD'), shell=False,
                                stderr=subprocess.STDOUT)
        outdata, errdata = proc.communicate()
        ret = proc.poll()
        # print(ret, proc.returncode, outdata, errdata)
        return ret

    def make_cmd(self):
        workdir = os.getenv('PWD')
        # -v
        print("Running swig...")
        cmd = "/usr/bin/swig -go -cgo -intgosize 64 -module core -o %s/src/qt/core/core_wrap.cxx -outdir %s/src/qt/core/ -I/usr/include/qt/QtCore -I/usr/include/qt -c++ ./core.swigcxx" % (workdir, workdir)
        # return cmd.split(' ')
        return shlex.split(cmd)

    def make_env(self):
        env = {}
        env['PATH'] = os.getenv('PATH')
        env['CC'] = '/home/dev/clang3.7/bin/clang'
        env['CXX'] = '/home/dev/clang3.7/bin/clang++'
        env['GOPATH'] = os.getenv('PWD')
        env['CGO_ENABLED'] = '1'
        env['CGO_LDFLAGS'] = '-lQt5Core'
        env['CGO_CXXFLAGS'] = "-I/usr/include/qt/QtCore -I/usr/include/qt -std=c++11 -DQT_CORE_LIB -DQT_NO_DEBUG -D_GNU_SOURCE -pipe -fno-exceptions -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -std=c++0x -Wall -W -D_REENTRANT -fPIC"
        env['CGO_CXXFLAGS'] += " -std=c++11"

        return env


def WarkerGo(WarkerOne):
    def __init__(self):
        super(self, WarkerGo).__init__()
        return


