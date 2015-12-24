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

    pass


class GenClassContext(object):
    def __init__(self, cursor):
        self.tyconv = TypeConvForRust()

        self.ctysz = max(32, cursor.type.get_size())  # 可能这个get_size()的值不准确啊。
        self.class_cursor = cursor
        self.class_name = cursor.spelling
        self.cursor = cursor

        # inherit
        self.base_class = None
        self.base_class_name = ''
        self.has_base = False

        # aux
        self.tymap = None
        # simple init

        self.CP = None
        return


class GenMethodContext(object):
    def __init__(self, cursor, class_cursor):
        self.tyconv = TypeConvForRust()

        self.ctysz = max(32, class_cursor.type.get_size())  # 可能这个get_size()的值不准确啊。
        self.class_cursor = class_cursor
        self.class_name = class_cursor.spelling
        self.cursor = cursor
        self.method_name = cursor.spelling
        self.method_name_rewrite = self.method_name
        self.mangled_name = cursor.mangled_name

        self.ctor = cursor.kind == clidx.CursorKind.CONSTRUCTOR
        self.dtor = cursor.kind == clidx.CursorKind.DESTRUCTOR

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

        self.params_cpp = ''
        self.params_rs = ''
        self.params_call = ''
        self.params_ext = ''

        self.unique_methods = {}
        self.struct_proto = '%s::%s%s' % (self.class_name, self.method_name, self.static_suffix)
        self.trait_proto = ''  # '%s::%s(%s)' % (class_name, method_name, trait_params)

        self.fn_proto_cpp = ''

        # inherit
        self.base_class = None
        self.base_class_name = ''
        self.has_base = False

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

