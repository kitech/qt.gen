# encoding: utf8

from genutil import *
from typeconv import TypeConv, TypeConvForRust


class GenerateBase(object):
    def __init__(self):
        super(GenerateBase, self).__init__()
        self.gctx = None
        self.gutil = GenUtil()
        self.tyconv = TypeConvForRust()
        return

    def setGenContext(self, ctx):
        self.gctx = ctx
        return

    def genpass(self):
        raise 'not impled'
        return

    def generateHeader(self, module):
        return

    def generateClasses(self, module, class_decls):
        return

    def generateClass(self, class_name, cs, methods):
        return

    def generateMethod(self, class_name, method_name, method_cursor):
        return

    # @return []
    def generateParams(self, class_name, method_name, method_cursor):
        return

    def method_is_inline(self, method_cursor):
        parent = method_cursor.semantic_parent
        for token in method_cursor.get_tokens():
            if token.spelling == 'inline':
                # print(111, method_cursor.spelling, parent.spelling)
                return True

        if not method_cursor.is_definition():
            defn = method_cursor.get_definition()
            if defn is not None: return self.method_is_inline(defn)

        return self.method_is_inline_ex(method_cursor)
        return False

    # 只好遍历整个 类的token了
    # 太慢无用
    def method_is_inline_ex(self, method_cursor):
        parent = method_cursor.semantic_parent
        inline_methods = self.gutil.get_inline_methods(parent)
        if method_cursor.mangled_name in inline_methods:
            return True
        return False

    def method_is_pure_virtual(self, method_cursor):
        return self.gutil.is_pure_virtual_method(method_cursor)

    def is_qt_class(self, type_name):
        # should be qt class name
        for seg in type_name.split(' '):
            if seg[0:1] == 'Q' and seg[1:2].upper() == seg[1:2] and '::' not in seg:  # should be qt class name
                return True
        return False

    def get_qt_class(self, type_name):
        # should be qt class name
        for seg in type_name.split(' '):
            if seg[0:1] == 'Q' and seg[1:2].upper() == seg[1:2] and '::' not in seg:  # should be qt class name
                if '<' in seg:
                    return seg.split('<')[0]
                return seg
        return None

    pass


class GenClassContext(object):
    def __init__(self, cursor):
        self.gutil = GenUtil()
        self.tyconv = TypeConvForRust()

        self.ctysz = max(32, cursor.type.get_size())  # 可能这个get_size()的值不准确啊。
        self.class_cursor = cursor
        self.class_name = cursor.type.spelling
        self.cursor = cursor

        self.full_class_name = cursor.type.spelling
        # 类内类处理
        if self.cursor.semantic_parent.kind == clidx.CursorKind.STRUCT_DECL or \
           self.cursor.semantic_parent.kind == clidx.CursorKind.CLASS_DECL:
            self.full_class_name = '%s::%s' % (self.cursor.semantic_parent.spelling, self.class_name)

        # inherit
        self.base_class = None
        self.base_class_name = ''
        self.has_base = False

        # signals
        self.signals = self.gutil.get_signals(cursor)

        # aux
        self.tymap = None
        # simple init

        self.CP = None
        return


class GenMethodContext(object):
    def __init__(self, cursor, class_cursor):
        self.gutil = GenUtil()
        self.tyconv = TypeConvForRust()

        self.ctysz = max(32, class_cursor.type.get_size())  # 可能这个get_size()的值不准确啊。
        self.class_cursor = class_cursor
        self.class_name = class_cursor.type.spelling
        self.cursor = cursor

        self.ctor = cursor.kind == clidx.CursorKind.CONSTRUCTOR
        self.dtor = cursor.kind == clidx.CursorKind.DESTRUCTOR

        self.method_name = cursor.spelling
        self.method_name_rewrite = self.method_name
        self.mangled_name = cursor.mangled_name
        if self.ctor: self.mangled_name = self.mangled_name.replace('C1', 'C2')
        if self.dtor: self.mangled_name = self.mangled_name.replace('D0Ev', 'D2Ev')

        self.full_class_name = self.class_cursor.type.spelling
        # 类内类处理
        if self.class_cursor.semantic_parent.kind == clidx.CursorKind.STRUCT_DECL or \
           self.class_cursor.semantic_parent.kind == clidx.CursorKind.CLASS_DECL:
            self.full_class_name = '%s::%s' % (self.class_cursor.semantic_parent.spelling, self.class_name)

        self.static = cursor.is_static_method()
        self.has_return = True
        self.ret_type = cursor.result_type
        self.ret_type_name_cpp = self.ret_type.spelling
        self.ret_type_name_rs = ''
        self.ret_type_name_ext = ''
        self.ret_type_ref = '&' in self.ret_type_name_cpp

        self.static_str = 'static' if self.static else ''
        self.static_suffix = '_s' if self.static else ''
        self.static_self_struct = '' if self.static else '& self, '
        self.static_self_trait = '' if self.static else ', rsthis: & %s' % (self.class_name)
        self.static_self_call = '' if self.static else 'self'

        self.isinline = False

        self.params_cpp = ''
        self.params_rs = ''
        self.params_call = ''
        self.params_ext_arr = []
        self.params_ext = ''

        self.unique_methods = {}
        self.struct_proto = '%s::%s%s' % (self.class_name, self.method_name, self.static_suffix)
        self.trait_proto = ''  # '%s::%s(%s)' % (class_name, method_name, trait_params)

        self.fn_proto_cpp = ''

        # inherit
        self.base_class = None
        self.base_class_name = ''
        self.has_base = False

        # signals
        self.signals = self.gutil.get_signals(class_cursor)

        # aux
        self.tymap = None
        # simple init

        self.CP = None
        return


# build the generated .swigcxx file
class TestBuilder:
    def tryBuild(self):
        return



class WarkerOne:
    def __init__(self): return

