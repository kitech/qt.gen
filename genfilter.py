# encoding: utf8

import clang.cindex as clidx


class GenFilter(object):
    def __init__(self):
        super(GenFilter, self).__init__()
        return

    def careDecl(self, cursor):
        return True

    def careClass(self, cursor):
        cname = cursor.spelling
        if cname.startswith('QMetaTypeId'): return False
        if cname.startswith('QTypeInfo'): return False
        return True

    def careMethod(self, cursor):
        if cursor.access_specifier != clidx.AccessSpecifier.PUBLIC:
            return False
        return True

    def careArg(self, cursor):
        return True


class GenFilterInline(GenFilter):
    def __init__(self):
        super(GenFilterInline, self).__init__()
        return


class GenFilterGo(GenFilter):
    def __init__(self):
        super(GenFilterGo, self).__init__()
        return


class GenFilterRust(GenFilter):
    def __init__(self):
        super(GenFilterRust, self).__init__()
        return
