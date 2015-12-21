# encoding: utf8

import clang.cindex as clidx

from genutil import CodePaper, GenUtil


class GenContext(object):
    def __init__(self):
        super(GenContext, self).__init__()

        self.modules = {}  # module name => tu cursor
        self.classes = {}  # class name => cursor
        self.clsmths = {}  # method name => cursor
        self.funcs = {}
        self.codes = {}  # file => CodePaper
        self.gutil = GenUtil()
        return

    def addClass(self, class_cursor):
        clskey = class_cursor.displayname
        # clskey = class_cursor.spelling

        if clskey not in self.classes:
            self.classes[clskey] = class_cursor
        else:
            # print('wtf', class_cursor.spelling, class_cursor.displayname)
            pass

        self.resolve_code_file(class_cursor)
        return

    def addFunction(self, func_cursor):
        funckey = func_cursor.spelling
        funckey = func_cursor.displayname
        funckey = func_cursor.mangled_name

        if funckey not in self.funcs:
            self.funcs[funckey] = func_cursor
        else:
            prev_cursor = self.funcs[funckey]
            # print('wtf', func_cursor.spelling, func_cursor.displayname,
            # func_cursor.is_definition(), prev_cursor.is_definition())
            pass

        self.resolve_code_file(func_cursor)
        return

    def getCodePager(self, cursor):
        decl_file = self.get_decl_file(cursor)
        return self.codes[decl_file]

    def is_template(self, cursor):
        if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE \
           or cursor.kind == clidx.CursorKind.FUNCTION_TEMPLATE:
            return True
        return False

    def resolve_code_file(self, cursor):
        decl_file = self.get_decl_file(cursor)
        if decl_file not in self.codes:
            self.codes[decl_file] = CodePaper()
            # print(decl_file)
            mod = self.get_decl_mod_by_path(decl_file)
        return

    def get_code_file(self, cursor):
        loc = cursor.location
        code_file = loc.file.name.split('/')[-1].split('.')[0]
        return code_file

    def get_code_file_by_path(self, path):
        # /QtWidgets/{qwhatthis}.h
        return path.split('/')[-1].split('.')[0]

    def get_decl_file(self, cursor):
        loc = cursor.location
        prefix = '/usr/include/qt'
        if loc.file.name.startswith(prefix):
            decl_file = loc.file.name[len(prefix):]
        else:
            print(cursor.kind, cursor.spelling, loc)
            raise '123'
        return decl_file

    def get_decl_module(self, cursor):
        return self.gutil.get_decl_module(cursor)

    def get_decl_mod(self, cursor):
        return self.gutil.get_decl_mod(cursor)

    def get_decl_mod_by_path(self, path):
        # /Qt{Widgets}/qwhatthis.h
        try:
            mod = path.split('/')[-2][2:].lower()
        except:
            print(path)
            raise '123'
        return mod

    def dumpContext(self, verbose = False):
        print('==========BEGIN')
        if verbose:
            print(self.classes.keys())
            print(self.funcs.keys())

        class_template_cnt = 0
        for key in self.classes:
            if self.is_template(self.classes[key]): class_template_cnt += 1
        func_template_cnt = 0
        for key in self.funcs:
            if self.is_template(self.funcs[key]): func_template_cnt += 1

        print('class count:', len(self.classes), class_template_cnt,
              'func count:', len(self.funcs), func_template_cnt)

        colen = 0
        coline = 0
        for key in self.codes: colen += self.codes[key].totalLength()
        for key in self.codes: coline += self.codes[key].totalLine()
        print('code line:', coline, 'code len:', colen)
        print('==========ENDDDDDDDD')
        return


class GenContextGlobal(GenContext):
    def __init__(self):
        super(GenContextGlobal, self).__init__()
        return


