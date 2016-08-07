# encoding: utf8

import clang.cindex as clidx
from genutil import GenUtil


class GenFilter(object):
    def __init__(self):
        super(GenFilter, self).__init__()
        self.gutil = GenUtil()
        return

    def skipDecl(self, cursor):
        return False

    def skipClass(self, cursor):
        cname = cursor.spelling
        if cname.startswith('QMetaTypeId'): return True
        if cname.startswith('QTypeInfo'): return True
        # 这个也许是因为qt有bug，也许是因为arch上的qt包有问题。QT_OPENGL_ES_2相关。
        if cname.startswith('QOpenGLFunctions_') and 'CoreBackend' in cname: return True
        if cname.startswith('QOpenGLFunctions_') and 'DeprecatedBackend' in cname: return True
        if cname.startswith('QOpenGLFunctions'): return True
        if cname.startswith('QOpenGLExtraFunctions'): return True
        if cname == 'QAbstractOpenGLFunctionsPrivate': return True
        if cname == 'QOpenGLFunctionsPrivate': return True
        if cname == 'QOpenGLExtraFunctionsPrivate': return True
        if cname.startswith('QOpenGLVersion'): return True
        if cname.startswith('QOpenGL'): return True

        if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True

        return False

    def skipMethod(self, cursor):
        if cursor.access_specifier != clidx.AccessSpecifier.PUBLIC:
            return True

        metamths = ['qt_metacall', 'qt_metacast', 'qt_check_for_']
        for mm in metamths:
            if cursor.spelling.startswith(mm): return True

        if cursor.spelling in ['tr', 'trUtf8', 'data_ptr']: return True

        if cursor.spelling.startswith('operator'): return True

        # retype = cursor.result_type
        # if retype.get_declaration() is not None:
        #     tdef = retype.get_declaration()
        #     if tdef.access_specifier in [clidx.AccessSpecifier.PRIVATE,
        #                                  clidx.AccessSpecifier.PROTECTED]:
        #         return True

        if cursor.spelling in ['rend', 'append', 'insert', 'rbegin', 'prepend', 'crend', 'crbegin']: return True

        return False

    def skipArg(self, cursor):
        return False


class GenFilterInc(GenFilter):
    def __init__(self):
        super(GenFilterInc, self).__init__()
        return

    def skipClass(self, cursor):
        if GenFilter.skipClass(self, cursor): return True

        cname = cursor.spelling
        if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True

        # 对于抽象类，还是要处理的

        return False

    def skipMethod(self, cursor):
        if GenFilter.skipMethod(self, cursor): return True

        return False


class GenFilterGo(GenFilter):
    def __init__(self):
        super(GenFilterGo, self).__init__()
        return

    def skipClass(self, cursor):
        if GenFilter.skipClass(self, cursor): return True

        cname = cursor.spelling
        if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True

        # 对于抽象类，还是要处理的

        return False

    def skipMethod(self, cursor):
        if GenFilter.skipMethod(self, cursor): return True
        if self.gutil.is_pure_virtual_method(cursor): return True

        return False


class GenFilterRust(GenFilter):
    def __init__(self):
        super(GenFilterRust, self).__init__()
        return

    def skipClass(self, cursor):
        if GenFilter.skipClass(self, cursor): return True

        return False

    def skipMethod(self, cursor):
        if GenFilter.skipMethod(self, cursor): return True

        return False
