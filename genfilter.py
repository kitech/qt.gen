# encoding: utf8

import clang.cindex as clidx


class GenFilter(object):
    def __init__(self):
        super(GenFilter, self).__init__()
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

        return False


class GenFilterGo(GenFilter):
    def __init__(self):
        super(GenFilterGo, self).__init__()
        return


class GenFilterRust(GenFilter):
    def __init__(self):
        super(GenFilterRust, self).__init__()
        return
