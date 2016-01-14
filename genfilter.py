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
        return False

    def skipMethod(self, cursor):
        if cursor.access_specifier != clidx.AccessSpecifier.PUBLIC:
            return True
        return False

    def skipArg(self, cursor):
        return False


class GenFilterInline(GenFilter):
    def __init__(self):
        super(GenFilterInline, self).__init__()
        return

    def skipClass(self, cursor):
        if GenFilter.skipClass(self, cursor): return True

        cname = cursor.spelling
        if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True

        # 对于抽象类，还是要处理的

        return False


class GenFilterGo(GenFilter):
    def __init__(self):
        super(GenFilterGo, self).__init__()
        return


class GenFilterRust(GenFilter):
    def __init__(self):
        super(GenFilterRust, self).__init__()
        return
